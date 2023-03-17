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
	"sync"

	"github.com/alpstable/cli/api/service/web"
	"github.com/alpstable/csvpb"
	"github.com/alpstable/gidari"
	"github.com/alpstable/gidari/auth"
	"github.com/alpstable/mongopb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

// parseWriteRequestTable parses a WriteRequest and extracts the table name. If
// the table name is not specified in the WriteRequest, it is extracted from the
// path of the URL. If an extension is provided, it is appended to the table
// name. The final table name is returned along with a possible error if URL
// parsing fails.
func parseWriteRequestTable(wr *web.WriteRequest, ext string) (string, error) {
	table := wr.Table
	if table == "" {
		u, err := url.Parse(wr.Url)
		if err != nil {
			return "", err
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
func newListWriterCSV(ctx context.Context, writer *web.Writer, cfg csvConfig) (*listWriter, func(), error) {
	if cfg.filename == "" {
		return nil, nil, fmt.Errorf("filename is required to create a CSV writer")
	}

	database := writer.GetDatabase()
	if database == "" {
		return nil, nil, fmt.Errorf("database is required to create a CSV writer")
	}

	loc := filepath.Join(database, cfg.filename)

	out, err := os.Create(loc)
	if err != nil {
		return nil, nil, err
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

func (mcp *mongoClientPool) close() {
	mcp.mu.Lock()
	defer mcp.mu.Unlock()

	for _, client := range mcp.clients {
		client.Disconnect(context.Background())
	}
}

func (mcp *mongoClientPool) getClient(ctx context.Context, connStr string) (*mongo.Client, error) {
	mcp.mu.Lock()
	defer mcp.mu.Unlock()

	client, ok := mcp.clients[connStr]
	if ok {
		return client, nil
	}

	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	mcp.clients[connStr] = client

	return client, nil
}

type mongoConfig struct {
	pool *mongoClientPool
	coll string
}

// addWebWriterMongo will set the HTTP Request and writer on the service for
// transport to a MongoDB instance.
func addWebWriterMongo(ctx context.Context, w *web.Writer, cfg mongoConfig) (*listWriter, func(), error) {
	connStr := w.GetConnString()
	if connStr == "" {
		return nil, nil, fmt.Errorf("connection string is required to create a MongoDB writer")
	}

	client, err := cfg.pool.getClient(ctx, connStr)
	if err != nil {
		return nil, nil, err
	}

	database := w.GetDatabase()

	coll := client.Database(database).Collection(cfg.coll)
	return newListWriter(mongopb.NewListWriter(coll)), func() {}, nil
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
			err = cerr
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
func (bldr *httpRequestBuilder) build() <-chan struct{} {
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
		opts[0] = gidari.WithAuth(newAuthRoundTrip(bldr.auth))

		wreq := bldr.writeRequest
		ctx := context.Background()

		for i, writer := range bldr.writers {
			switch writer.GetType() {
			case web.WriteType_CSV:
				fn, err := parseWriteRequestTable(wreq, "csv")
				if err != nil {
					bldr.errs <- err

					return
				}

				cfg := csvConfig{
					filename: fn,
				}

				lw, c, err := newListWriterCSV(ctx, writer, cfg)
				if c != nil {
					bldr.closer <- c
				}

				if err != nil {
					bldr.errs <- err

					return
				}

				opts[i+1] = gidari.WithWriter(lw.glw)
			case web.WriteType_MONGO:
				coll, err := parseWriteRequestTable(wreq, "")
				if err != nil {
					bldr.errs <- err

					return
				}

				cfg := mongoConfig{
					pool: bldr.mongoClientPool,
					coll: coll,
				}

				lw, c, err := addWebWriterMongo(ctx, writer, cfg)
				if c != nil {
					bldr.closer <- c
				}

				if err != nil {
					bldr.errs <- err

					return
				}

				opts[i+1] = gidari.WithWriter(lw.glw)
			}
		}

		wr := bldr.writeRequest
		if wr == nil {
			return
		}

		// Create the HTTP Request.
		httpReq, _ := http.NewRequest(wr.Method, wr.Url, nil)

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
		return nil, err
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
			err = cerr
		}
	}

	return err
}

func (bldr *httpServiceBuilder) build() <-chan struct{} {
	svc := bldr.gidariService
	req := bldr.req

	mongoClientPool := newMongoClientPool()
	defer mongoClientPool.close()

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
				build()

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

	for range bldr.build() {
	}

	defer func() {
		if err := bldr.close(); err != nil {
			panic(err)
		}
	}()

	// Execute the service.
	if err := bldr.gidariService.HTTP.Upsert(ctx); err != nil {
		return nil, err
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
