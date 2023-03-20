// Copyright 2023 The Gidari CLI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alpstable/cli/api/service/web"
	"github.com/alpstable/csvpb"
	"github.com/alpstable/gidari"
	"github.com/alpstable/gidari/auth"
	"github.com/alpstable/mongopb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

// parseWriteRequestTable parses a WriteRequest and extracts the table name. If
// the table name is not specified in the WriteRequest, it is extracted from the
// path of the URL. If an extension is provided, it is appended to the table
// name. The final table name is returned along with a possible error if URL
// parsing fails.
func parseWriteRequestTable(wr *web.WriteRequest, table string, ext string) (string, error) {
	if table == "" {
		u, err := url.Parse(wr.Url)
		if err != nil {
			return "", fmt.Errorf("failed to parse URL: %w", err)
		}

		table = path.Base(u.Path)
	}

	if ext != "" {
		table = fmt.Sprintf("%s.%s", table, ext)
	}

	return table, nil
}

type listWriter struct {
	glw gidari.ListWriter
}

func newListWriter(lw gidari.ListWriter) *listWriter {
	return &listWriter{glw: lw}
}

// newAuthRoundTrip returns a function that can be used as a RoundTripper in an
// HTTP client, with authentication provided by the given web.Auth object.
func newAuthRoundTrip(wa *web.Auth) func(*http.Request) (*http.Response, error) {
	if cbauth := wa.GetCoinbase(); cbauth != nil {
		return auth.NewCoinbaseRoundTrip(cbauth.Key, cbauth.Secret, cbauth.Passphrase)
	}

	return nil
}

type csvConfig struct {
	filename string
}

// newListWriterCSV creates a new CSV file at the given filename, using the
// web.Writer's database path as the root directory. It returns a
// gidari.ListWriter that writes to the CSV file, as well as a cleanup function
// to be called when writing is complete. If an error occurs, it returns nil for
// both the gidari.ListWriter and cleanup function, along with the error.
//
// The returned gidari.ListWriter is backed by a csv.Writer, and is responsible
// for writing CSV records to the output file. The cleanup function must be
// called when writing is complete to flush the csv.Writer buffer to the output
// file and close the file handle. If an error occurs during cleanup, it panics.
func newListWriterCSV(writer *web.Writer, cfg csvConfig) (*listWriter, func(), error) {
	if cfg.filename == "" {
		return nil, nil, fmt.Errorf("filename is required to create a CSV writer")
	}

	database := writer.GetDatabase()
	if database == "" {
		return nil, nil, fmt.Errorf("database is required to create a CSV writer")
	}

	loc := filepath.Join(database, cfg.filename)

	// Make sure the filename ends with .csv to avoid vulnerabilities.
	if !strings.HasSuffix(loc, ".csv") {
		return nil, nil, fmt.Errorf("filename must end with .csv")
	}

	// Santize the filepath before using it.
	loc = filepath.Clean(loc)

	out, err := os.Create(loc) //nolint:gosec
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create output file: %w", err)
	}

	// Create a new CSV Writer for the output file.
	csvWriter := csv.NewWriter(out)

	return newListWriter(csvpb.NewListWriter(csvWriter)), func() {
		csvWriter.Flush()

		if err := out.Close(); err != nil {
			panic(err)
		}
	}, nil
}

type mongoClientPool struct {
	mu sync.Mutex

	clients map[string]*mongo.Client
}

func newMongoClientPool() *mongoClientPool {
	return &mongoClientPool{
		clients: make(map[string]*mongo.Client),
	}
}

func (mcp *mongoClientPool) close(ctx context.Context) error {
	mcp.mu.Lock()
	defer mcp.mu.Unlock()

	for _, client := range mcp.clients {
		if err := client.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
		}
	}

	return nil
}

func (mcp *mongoClientPool) getClient(ctx context.Context, connStr string) (*mongo.Client, error) {
	mcp.mu.Lock()
	defer mcp.mu.Unlock()

	client, ok := mcp.clients[connStr]
	if ok {
		return client, nil
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	mcp.clients[connStr] = client

	return client, nil
}

type mongoConfig struct {
	pool  *mongoClientPool
	coll  string
	index *mongo.IndexModel
}

// parseMongoIndexModel parses a list of strings into a mongo.IndexModel.
func parseMongoIndexModel(index []string) *mongo.IndexModel {
	if len(index) == 0 {
		return nil
	}

	keys := make(bson.D, len(index))
	for i, key := range index {
		keys[i] = bson.E{Key: key, Value: 1}
	}

	return &mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	}
}

// addWebWriterMongo will set the HTTP Request and writer on the service for
// transport to a MongoDB instance.
func addWebWriterMongo(ctx context.Context, w *web.Writer, cfg mongoConfig) (*listWriter, func(), error) {
	connStr := w.GetConnString()
	if connStr == "" {
		return nil, nil, fmt.Errorf("connection string is required to create a mongodb writer")
	}

	client, err := cfg.pool.getClient(ctx, connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create mongodb client: %w", err)
	}

	database := w.GetDatabase()

	coll := client.Database(database).Collection(cfg.coll)
	lwriter := mongopb.NewListWriter(coll)

	// Create an index if one was provided and it is not empty.
	if cfg.index != nil && cfg.index.Keys != nil {
		_, err := coll.Indexes().CreateOne(ctx, *cfg.index)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create mongodb index: %w", err)
		}
	}

	return newListWriter(lwriter), func() {}, nil
}

// httpRequestBuilder will build a gidari HTTP request from a web.WriteRequest.
type httpRequestBuilder struct {
	gidariService   *gidari.Service
	writeRequest    *web.WriteRequest
	mongoClientPool *mongoClientPool

	writers []*web.Writer
	auth    *web.Auth

	errs   chan error
	closer chan func()
}

// close closes the list writer used to construct the HTTP service request and
// performs any necessary teardown operations, such as flushing a writer buffer.
// It ranges over the bldr.errs channel to return the first error encountered,
// if any. This method blocks until the errors and closers channels have been
// closed by the "build" method.
func (bldr *httpRequestBuilder) close() error {
	defer func() {
		for c := range bldr.closer {
			c()
		}
	}()

	var err error

	for cerr := range bldr.errs {
		if cerr != nil {
			err = fmt.Errorf("failed to build HTTP request: %w", cerr)
		}
	}

	return err
}

// setAuth sets the authentication round trip to use when building a request for
// an HTTP service. The "globalAuth" argument will be used if the builder's
// write request has no authentication set.
func (bldr *httpRequestBuilder) setAuth(globAuth *web.Auth) *httpRequestBuilder {
	if bldr.writeRequest != nil && bldr.writeRequest.Auth != nil {
		bldr.auth = bldr.writeRequest.Auth
	} else {
		bldr.auth = globAuth
	}

	return bldr
}

// setWriter sets the writer to use when building a request for an HTTP service.
// The "globalWriter" argument will be used if the builder's write request has
// no writers set.
func (bldr *httpRequestBuilder) setWriters(globWriters []*web.Writer) *httpRequestBuilder {
	if bldr.writeRequest != nil && len(bldr.writeRequest.Writers) > 0 {
		bldr.writers = bldr.writeRequest.GetWriters()
	} else {
		bldr.writers = globWriters
	}

	return bldr
}

// build builds a gidari HTTP request from the builder's write request and
// returns a channel that will be closed when the request has been built.
func (bldr *httpRequestBuilder) build(ctx context.Context) <-chan struct{} {
	bldr.errs = make(chan error, len(bldr.writers))
	bldr.closer = make(chan func(), len(bldr.writers))

	done := make(chan struct{}, 1)

	go func() {
		defer close(bldr.errs)
		defer close(bldr.closer)
		defer close(done)

		if bldr.gidariService == nil || bldr.gidariService.HTTP == nil {
			return
		}

		opts := make([]gidari.RequestOption, len(bldr.writers)+1)
		opts[0] = gidari.WithAuth(newAuthRoundTrip(bldr.auth)) //nolint:bodyclose

		wreq := bldr.writeRequest

		for i, writer := range bldr.writers {
			switch writer.GetType() {
			case web.WriteType_CSV:
				var table string
				if csv := wreq.GetCsv(); csv != nil {
					table = wreq.Csv.GetFile()
				}

				fileName, err := parseWriteRequestTable(wreq, table, "csv")
				if err != nil {
					bldr.errs <- err

					return
				}

				cfg := csvConfig{
					filename: fileName,
				}

				listW, c, err := newListWriterCSV(writer, cfg)
				if c != nil {
					bldr.closer <- c
				}

				if err != nil {
					bldr.errs <- err

					return
				}

				opts[i+1] = gidari.WithWriters(listW.glw)
			case web.WriteType_MONGO:
				var (
					collName string
					index    []string
				)

				if mongo := wreq.GetMongo(); mongo != nil {
					collName = wreq.Mongo.GetCollection()
					index = wreq.Mongo.GetIndex()
				}

				coll, err := parseWriteRequestTable(wreq, collName, "")
				if err != nil {
					bldr.errs <- err

					return
				}

				cfg := mongoConfig{
					pool:  bldr.mongoClientPool,
					coll:  coll,
					index: parseMongoIndexModel(index),
				}

				listW, c, err := addWebWriterMongo(ctx, writer, cfg)
				if c != nil {
					bldr.closer <- c
				}

				if err != nil {
					bldr.errs <- err

					return
				}

				opts[i+1] = gidari.WithWriters(listW.glw)
			}
		}

		if wreq == nil {
			return
		}

		// Create the HTTP Request.
		httpReq, _ := http.NewRequestWithContext(ctx, wreq.Method, wreq.Url, nil)

		// Add th query params to the HTTP request.
		query := httpReq.URL.Query()
		for key, value := range wreq.GetQueryParams() {
			query.Add(key, value)
		}

		httpReq.URL.RawQuery = query.Encode()

		// Add the request options to the service.
		gidReq := gidari.NewHTTPRequest(httpReq, opts...)
		bldr.gidariService.HTTP.Requests(gidReq)
	}()

	return done
}

// httpServiceBuilder is used to build a gidari HTTP service from a proto
// web.Request. The resulting service will be used to execute a bulk "fetch"
// operation to receive the data from the web and write it to storage.
type httpServiceBuilder struct {
	gidariService *gidari.Service
	req           *web.Request

	closers chan func() error
}

func newHTTPServiceBuilder(ctx context.Context, req *web.Request) (*httpServiceBuilder, error) {
	svc, err := gidari.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	bldr := &httpServiceBuilder{
		gidariService: svc,
		req:           req,
	}

	return bldr, nil
}

func (bldr *httpServiceBuilder) close() error {
	var err error

	for closer := range bldr.closers {
		if cerr := closer(); cerr != nil {
			err = fmt.Errorf("failed to HTTP service builder close: %w", cerr)
		}
	}

	return err
}

func (bldr *httpServiceBuilder) build(ctx context.Context) <-chan struct{} {
	svc := bldr.gidariService
	req := bldr.req

	mongoClientPool := newMongoClientPool()
	defer func() {
		if err := mongoClientPool.close(ctx); err != nil {
			panic(err)
		}
	}()

	writeRequests := req.GetRequests()

	bldr.closers = make(chan func() error, len(writeRequests))
	done := make(chan struct{}, len(writeRequests))

	go func() {
		defer close(bldr.closers)
		defer close(done)

		for _, wr := range writeRequests {
			reqBldr := &httpRequestBuilder{
				gidariService:   svc,
				writeRequest:    wr,
				mongoClientPool: mongoClientPool,
			}

			done <- <-reqBldr.setAuth(req.GetAuth()).
				setWriters(req.GetWriters()).
				build(ctx)

			bldr.closers <- reqBldr.close
		}
	}()

	return done
}

type webServer struct {
	web.UnimplementedWebServer
}

func newWebServer() *webServer {
	s := &webServer{}

	return s
}

// Write will build an HTTP Service to query the web APIs defined by the
// request, writing the results to the requested writers.
func (svr *webServer) Write(ctx context.Context, req *web.Request) (*web.Response, error) {
	bldr, err := newHTTPServiceBuilder(ctx, req)
	if err != nil {
		return nil, err
	}

	for range bldr.build(ctx) {
	}

	defer func() {
		if err := bldr.close(); err != nil {
			panic(err)
		}
	}()

	// Execute the service.
	if err := bldr.gidariService.HTTP.Upsert(ctx); err != nil {
		return nil, fmt.Errorf("failed to execute service: %w", err)
	}

	return &web.Response{}, nil
}

// AddWebServerCommand adds the web server command to the CLI.
func AddWebServerCommand(cmd *cobra.Command) {
	type webServerConfig struct {
		host string
		Port string
	}

	var cfg webServerConfig

	scmd := &cobra.Command{
		Use:   "start",
		Short: "Start the RPC server",
		Run: func(cmd *cobra.Command, args []string) {
			addr := net.JoinHostPort(cfg.host, cfg.Port)

			logrus.Infof("Starting RPC server on %s", addr)

			lis, err := net.Listen("tcp", addr)
			if err != nil {
				panic(err)
			}

			s := grpc.NewServer()
			web.RegisterWebServer(s, newWebServer())

			if err := s.Serve(lis); err != nil {
				panic(err)
			}
		},
	}

	flags := scmd.Flags()
	flags.StringVar(&cfg.host, "host", "localhost", "host to listen on")
	flags.StringVar(&cfg.Port, "port", "50051", "port to listen on")

	cmd.AddCommand(scmd)
}
