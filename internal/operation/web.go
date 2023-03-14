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

	"github.com/alpstable/cli/api/service/web"
	"github.com/alpstable/csvpb"
	"github.com/alpstable/gidari"
	"github.com/alpstable/gidari/auth"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

// newAuthRoundTrip returns a function that can be used as a RoundTripper in an
// HTTP client, with authentication provided by the given web.Auth object.
func newAuthRoundTrip(wa *web.Auth) func(*http.Request) (*http.Response, error) {
	if cbauth := wa.GetCoinbase(); cbauth != nil {
		return auth.NewCoinbaseRoundTrip(cbauth.Key, cbauth.Secret, cbauth.Passphrase)
	}

	return nil
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
func newListWriterCSV(ctx context.Context, filename string, writer *web.Writer) (gidari.ListWriter, func(), error) {
	if filename == "" {
		return nil, nil, fmt.Errorf("filename is required to create a CSV writer")
	}

	database := writer.GetDatabase()
	if database == "" {
		return nil, nil, fmt.Errorf("database is required to create a CSV writer")
	}

	loc := filepath.Join(database, filename)

	out, err := os.Create(loc)
	if err != nil {
		return nil, nil, err
	}

	// Create a new CSV Writer for the output file.
	csvWriter := csv.NewWriter(out)
	return csvpb.NewListWriter(csvWriter), func() {
		csvWriter.Flush()
		fmt.Println("flush is closed.")

		if err := out.Close(); err != nil {
			panic(err)
		}
	}, nil
}

// httpRequestBuilder will build a gidari HTTP request from a web.WriteRequest.
type httpRequestBuilder struct {
	gidariService *gidari.Service
	writeRequest  *web.WriteRequest

	writers []*web.Writer
	auth    *web.Auth

	errs   chan error
	closer chan func()
}

func newHTTPRequestBuilder(svc *gidari.Service, writeRequest *web.WriteRequest) *httpRequestBuilder {
	return &httpRequestBuilder{
		gidariService: svc,
		writeRequest:  writeRequest,
	}
}

// close closes the list writer used to construct the HTTP service request and
// performs any necessary teardown operations, such as flushing a writer buffer.
// It ranges over the bldr.errs channel to return the first error encountered,
// if any. This method blocks until the errors and closers channels have been
// closed by the "build" method.
func (bldr *httpRequestBuilder) close() error {
	defer func() {
		for c := range bldr.closer {
			fmt.Println("closer is closed.")
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
func (bldr *httpRequestBuilder) setAuth(globalAuth *web.Auth) *httpRequestBuilder {
	if bldr.writeRequest != nil && bldr.writeRequest.Auth != nil {
		bldr.auth = bldr.writeRequest.Auth
	} else {
		bldr.auth = globalAuth
	}

	return bldr
}

func (bldr *httpRequestBuilder) setWriters(def []*web.Writer) *httpRequestBuilder {
	if bldr.writeRequest != nil && len(bldr.writeRequest.Writers) > 0 {
		def = bldr.writeRequest.GetWriters()
	}

	bldr.writers = def

	return bldr
}

func (bldr *httpRequestBuilder) build() <-chan struct{} {
	bldr.errs = make(chan error, len(bldr.writers))
	bldr.closer = make(chan func(), len(bldr.writers))

	done := make(chan struct{}, 1)

	go func() {
		defer close(bldr.errs)
		defer close(bldr.closer)
		defer close(done)

		opts := make([]gidari.RequestOption, len(bldr.writers)+1)
		opts[0] = gidari.WithAuth(newAuthRoundTrip(bldr.auth))

		for i, writer := range bldr.writers {
			fmt.Println("i, type: ", i, writer.GetType())
			switch writer.GetType() {
			case web.WriteType_CSV:
				fn, err := parseWriteRequestTable(bldr.writeRequest, "csv")
				if err != nil {

				}

				lw, c, err := newListWriterCSV(context.Background(), fn, writer)
				if c != nil {
					bldr.closer <- c
				}

				if err != nil {
					fmt.Println("err: ", err)
					bldr.errs <- err

					return
				}

				opts[i+1] = gidari.WithWriter(lw)
			}
		}

		wr := bldr.writeRequest

		// Create the HTTP Request.
		httpReq, _ := http.NewRequest(wr.Method, wr.Url, nil)

		// Add the request options to the service.
		gidReq := gidari.NewHTTPRequest(httpReq, opts...)
		bldr.gidariService.HTTP.Requests(gidReq)

		fmt.Println("build is done.")
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
			fmt.Println("error closing", cerr)
			err = cerr
		}

		fmt.Println("post close")
	}

	return err
}

func (bldr *httpServiceBuilder) build() <-chan struct{} {
	svc := bldr.gidariService
	req := bldr.req

	writeRequests := req.GetRequests()

	bldr.closers = make(chan func() error, len(writeRequests))
	done := make(chan struct{}, len(writeRequests))

	go func() {
		defer close(bldr.closers)
		defer close(done)

		for _, wr := range writeRequests {
			reqBldr := newHTTPRequestBuilder(svc, wr)
			done <- <-reqBldr.setAuth(req.GetAuth()).
				setWriters(req.GetWriters()).
				build()

			bldr.closers <- reqBldr.close
			fmt.Println("added")
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

	fmt.Println("before execute")

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

//type webServer struct {
//	web.UnimplementedWebServer
//}
//
//func newWebServer() *webServer {
//	s := &webServer{}
//
//	return s
//}
//
//// webRequest wraps a web.Request.
//type webRequest struct {
//	in           *web.Request
//	mongoClients map[string]*mongo.Client
//}
//
//func newWebRequest(in *web.Request) *webRequest {
//	return &webRequest{
//		in:           in,
//		mongoClients: make(map[string]*mongo.Client),
//	}
//}
//
//// webWriteRequest wraps a web.WriteRequest and it's parent request.
//type webWriteRequest struct {
//	parent *webRequest
//	in     *web.WriteRequest
//}
//
//// newWebWriteRequest will construct an object to handle an individeual write
//// request for a web server request.
//func newWebWriteRequest(p *webRequest, in *web.WriteRequest) *webWriteRequest {
//	return &webWriteRequest{
//		parent: p,
//		in:     in,
//	}
//}
//
//// writerCfgs will return the list of writerCfgs from the request. If there are
//// no request-specific writerCfgs, then global writerCfgs are used.
//func (req *webWriteRequest) writerCfgs() []*web.Writer {
//	if req.in != nil && len(req.in.Writers) > 0 {
//		return req.in.GetWriters()
//	}
//
//	if req.parent.in != nil && len(req.parent.in.Writers) > 0 {
//		return req.parent.in.GetWriters()
//	}
//
//	return nil
//}
//
//// auth will return the auth from the request. If there is no request-specific
//// auth, then global auth is used.
//func (req *webWriteRequest) auth() *web.Auth {
//	if req.in != nil && req.in.Auth != nil {
//		return req.in.GetAuth()
//	}
//
//	if req.parent.in != nil && req.parent.in.Auth != nil {
//		return req.parent.in.GetAuth()
//	}
//
//	return nil
//}
//
//// table will return the name of the data object to write to. If there is no
//// request-specific table, then the table is derived from the last part of the
//// URL.
//func (req *webWriteRequest) table(ext string) (string, error) {
//	table := req.in.Table
//	if table == "" {
//		u, err := url.Parse(req.in.Url)
//		if err != nil {
//			return "", err
//		}
//
//		table = path.Base(u.Path)
//	}
//
//	if ext != "" {
//		table = fmt.Sprintf("%s.%s", table, ext)
//	}
//
//	return table, nil
//}
//
//type svcWR func() (*gidari.Service, *webWriteRequest)
//
//func newSvcWR(svc *gidari.Service, req *webWriteRequest) svcWR {
//	return func() (*gidari.Service, *webWriteRequest) {
//		return svc, req
//	}
//}
//
//// addWebWriterCSV will set the HTTP Request and writer on the service for
//// transport to local storage.
//func addWebWriterCSV(ctx context.Context, w *web.Writer, get svcWR) (<-chan func(), <-chan error) {
//	svc, req := get()
//
//	errs := make(chan error, 1)
//	closers := make(chan func(), 1)
//
//	table, err := req.table("csv")
//	if err != nil {
//		errs <- err
//
//		return closers, errs
//	}
//
//	reqAuth := req.auth()
//
//	loc := filepath.Join(w.GetDatabase(), table)
//
//	out, err := os.Create(loc)
//	if err != nil {
//		errs <- err
//
//		return closers, errs
//	}
//
//	closers <- func() {
//		if err := out.Close(); err != nil {
//			panic(err)
//		}
//	}
//
//	// Create a new CSV Writer for the output file.
//	csvWriter := csv.NewWriter(out)
//	csvListWriter := csvpb.NewListWriter(csvWriter)
//
//	// Create the request options.
//	gidWriter := gidari.WithWriter(csvListWriter)
//
//	gidAuth := newAuthFromWeb(reqAuth)
//	if gidAuth == nil {
//		gidAuth = newAuthFromWeb(reqAuth)
//	}
//
//	// Create the HTTP Request.
//	httpReq, err := http.NewRequest(req.in.Method, req.in.Url, nil)
//	if err != nil {
//		errs <- err
//
//		return closers, errs
//	}
//
//	// Add the request options to the service.
//	gidReq := gidari.NewHTTPRequest(httpReq, gidWriter, gidAuth)
//	svc.HTTP.Requests(gidReq)
//
//	return closers, errs
//}
//
//// addWebWriterMongo will set the HTTP Request and writer on the service for
//// transport to a MongoDB instance.
//func addWebWriterMongo(ctx context.Context, w *web.Writer, get svcWR) (<-chan func(), <-chan error) {
//	svc, req := get()
//	connStr := w.GetConnString()
//
//	errs := make(chan error, 1)
//	closers := make(chan func(), 1)
//
//	// Get the client from the map.
//	client, ok := req.parent.mongoClients[connStr]
//	if !ok {
//		// Create a new client.
//		var err error
//		client, err = mongo.NewClient(options.Client().ApplyURI(connStr))
//		if err != nil {
//			errs <- err
//
//			return closers, errs
//		}
//
//		// Connect to the client.
//		if err := client.Connect(ctx); err != nil {
//			errs <- err
//
//			return closers, errs
//		}
//
//		closers <- func() {
//			if err := client.Disconnect(ctx); err != nil {
//				panic(err)
//			}
//		}
//
//		// Ping the client.
//		if err := client.Ping(ctx, nil); err != nil {
//			errs <- err
//
//			return closers, errs
//		}
//	}
//
//	database := w.GetDatabase()
//
//	collection, err := req.table("")
//	if err != nil {
//		errs <- err
//
//		return closers, errs
//	}
//
//	if collection == "" {
//		u, err := url.Parse(req.in.Url)
//		if err != nil {
//			errs <- err
//
//			return closers, errs
//		}
//
//		collection = path.Base(u.Path)
//	}
//
//	coll := client.Database(database).Collection(collection)
//	listWriter := mongopb.NewListWriter(coll)
//
//	// Create the request options.
//	gidWriter := gidari.WithWriter(listWriter)
//
//	auth := req.auth()
//
//	gidAuth := newAuthFromWeb(auth)
//	if gidAuth == nil {
//		gidAuth = newAuthFromWeb(auth)
//	}
//
//	// Create the HTTP Request.
//	httpReq, err := http.NewRequest(req.in.Method, req.in.Url, nil)
//	if err != nil {
//		errs <- err
//
//		return closers, errs
//	}
//
//	// Add the request options to the service.
//	gidReq := gidari.NewHTTPRequest(httpReq, gidWriter, gidAuth)
//	svc.HTTP.Requests(gidReq)
//
//	return closers, errs
//}
//
//type webWriterSetter func(context.Context, *web.Writer, svcWR) (<-chan func(), <-chan error)
//
//// webWriterSetterCtor is a function that will return a function that can
//// be usd to set conditions for writting an HTTP to storage.
//type webWriterSetterCtor func(web.WriteType) (webWriterSetter, bool)
//
//func newWebWriterSetterCSV(writeType web.WriteType) (webWriterSetter, bool) {
//	if writeType != web.WriteType_CSV {
//		return nil, false
//	}
//
//	return addWebWriterCSV, true
//}
//
//func newWebWriterSetterMongo(writeType web.WriteType) (webWriterSetter, bool) {
//	if writeType != web.WriteType_MONGO {
//		return nil, false
//	}
//
//	return addWebWriterMongo, true
//}
//
//// setWebWriter will set the HTTP Request and writer options on the gidari
//// service for a web-to-storage transport. This function takes a "get" function
//// that returns the web request and the service to write to.
////
//// Note that the webWriterCtor functions ARE NOT the writing configuration,
//// they are functions that return a writer setter. The writer configurations are
//// defined on the request returend by the "get" function. For this reason, if
//// there are myltiple constructors of the same type, the first one will be used.
//func setWebWriter(ctx context.Context, get svcWR, ctors ...webWriterSetterCtor) (<-chan func(), <-chan error) {
//	_, req := get()
//	writerCfgs := req.writerCfgs()
//
//	closers := make(chan func(), len(writerCfgs))
//	errCh := make(chan error, len(writerCfgs))
//
//	// Set the writers asynchronously so that we can collect closers. We
//	// do this to ensure that we can gracefully close any open resources
//	// if an error occurs mid-way through the process.
//	go func() {
//		defer close(closers)
//		defer close(errCh)
//
//		for _, cfg := range writerCfgs {
//			var setWriter webWriterSetter
//
//			// Find the first constructor for the writer's WriteType
//			// and use it to create the writer.
//			for _, newWriter := range ctors {
//				var ok bool
//
//				setWriter, ok = newWriter(cfg.Type)
//				if !ok {
//					continue
//				}
//
//				// Once we find a writer that matches the type,
//				// we can break out of the loop.
//				break
//			}
//
//			// If no writer was found, then we can't write to this
//			// type of storage.
//			if setWriter == nil {
//				errCh <- fmt.Errorf("web writer not found")
//				closers <- nil
//
//				return
//			}
//
//			// Set the writer on the gidari service.
//			wclosers, werrs := setWriter(ctx, cfg, get)
//
//			// Even though we iterate over the wclosers channel, we
//			// only expect one closer to be returned.
//			for closer := range wclosers {
//				closers <- closer
//			}
//
//			for err := range werrs {
//				errCh <- err
//			}
//		}
//	}()
//
//	return closers, errCh
//}
//
//// svcR is a function that wraps a gidari service and a web request. This is
//// used to hand off "paired" data to other functions.
//type svcR func() (*gidari.Service, *webRequest)
//
//func newSvcR(svc *gidari.Service, req *webRequest) svcR {
//	return func() (*gidari.Service, *webRequest) {
//		return svc, req
//	}
//}
//
//// setWebRequest will set the web request's write requests onto a gidari
//// service for bulk transport.
//func setWebRequest(ctx context.Context, get svcR) (<-chan func(), <-chan error) {
//	svc, req := get()
//	writeRequests := req.in.GetRequests()
//
//	errs := make(chan error, len(writeRequests))
//	closers := make(chan func(), len(writeRequests))
//
//	go func() {
//		defer close(errs)
//		defer close(closers)
//
//		for _, wr := range writeRequests {
//			get := newSvcWR(svc, newWebWriteRequest(req, wr))
//
//			// Set all relevant web writers on the gidari service.
//			rclosers, rerr := setWebWriter(ctx, get,
//				newWebWriterSetterCSV,
//				newWebWriterSetterMongo,
//			)
//
//			for rcloser := range rclosers {
//				closers <- rcloser
//			}
//
//			for err := range rerr {
//				errs <- err
//			}
//		}
//	}()
//
//	return closers, errs
//}
