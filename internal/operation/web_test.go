package operation

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/alpstable/cli/api/service/web"
	"github.com/alpstable/gidari"
	"github.com/alpstable/gidari/auth"
)

func TestParseWriteRequestTable(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name          string
		req           *web.WriteRequest
		ext           string
		expected      string
		expectedError bool
	}{
		{
			name: "table name is empty in WriteRequest w ext",
			req: &web.WriteRequest{
				Table: "",
				Url:   "https://example.com/path/to/table",
			},
			ext:      "csv",
			expected: "table.csv",
		},
		{
			name: "table name is empty in WriteRequest wo ext",
			req: &web.WriteRequest{
				Table: "",
				Url:   "https://example.com/path/to/table",
			},
			ext:      "",
			expected: "table",
		},
		{
			name: "table name is specified in WriteRequest w ext",
			req: &web.WriteRequest{
				Table: "mytable",
				Url:   "https://example.com/path/to/other/table",
			},
			ext:      "csv",
			expected: "mytable.csv",
		},
		{
			name: "table name is specified in WriteRequest wo ext",
			req: &web.WriteRequest{
				Table: "mytable",
				Url:   "https://example.com/path/to/other/table",
			},
			ext:      "",
			expected: "mytable",
		},
		{
			name: "URL parsing fails",
			req: &web.WriteRequest{
				Table: "",
				Url:   "::invalidurl",
			},
			ext:           "csv",
			expected:      "",
			expectedError: true,
		},
		{
			name: "extension is empty",
			req: &web.WriteRequest{
				Table: "",
				Url:   "https://example.com/path/to/table",
			},
			ext:      "",
			expected: "table",
		},
	} {
		tcase := tcase
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			table, err := parseWriteRequestTable(tcase.req, tcase.ext)

			if tcase.expectedError && err == nil {
				t.Fatalf("expected error, but got nil")
			}

			if table != tcase.expected {
				t.Errorf("Expected table name: %s, but got: %s", tcase.expected, table)
			}
		})
	}
}

func TestNewAuthRoundTrip(t *testing.T) {
	t.Parallel()

	newTestCoinbaseAuth := func(key, secret, pp string) *web.Auth {
		return &web.Auth{Auth: &web.Auth_Coinbase{
			Coinbase: &web.CoinbaseAuth{
				Key:        key,
				Secret:     secret,
				Passphrase: pp,
			},
		}}
	}

	for _, tcase := range []struct {
		name          string
		wa            *web.Auth
		expected      func(*http.Request) (*http.Response, error)
		expectedError bool
	}{
		{
			name:          "coinbase auth credentials provided",
			wa:            newTestCoinbaseAuth("mykey", "mysecret", "mypassphrase"),
			expected:      auth.NewCoinbaseRoundTrip("mykey", "mysecret", "mypass"),
			expectedError: false,
		},
		{
			name:          "coinbase auth credentials not provided",
			wa:            &web.Auth{},
			expected:      nil,
			expectedError: false,
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			roundTripFunc := newAuthRoundTrip(tcase.wa)

			if roundTripFunc == nil && tcase.expected != nil {
				t.Errorf("Expected a non-nil round trip function, but got nil")
			}

			if roundTripFunc != nil && tcase.expected == nil {
				t.Errorf("Expected a nil round trip function, but got a non-nil function")
			}

			actual := reflect.ValueOf(roundTripFunc).Pointer()
			expected := reflect.ValueOf(tcase.expected).Pointer()

			if actual != expected {
				t.Errorf("Expected round trip function does not match actual function")
			}
		})
	}
}

func TestNewListWriterCSV(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name          string
		filename      string
		writer        *web.Writer
		expectedError bool
	}{
		{
			name:     "valid csv writer",
			filename: "test.csv",
			writer: &web.Writer{
				Type:     web.WriteType_CSV,
				Database: "some/file/path",
			},
			expectedError: false,
		},
		{
			name:     "valid csv writer with no filename",
			filename: "",
			writer: &web.Writer{
				Type:     web.WriteType_CSV,
				Database: "some/file/path",
			},
			expectedError: true,
		},
		{
			name:     "valid csv writer with no database",
			filename: "test.csv",
			writer: &web.Writer{
				Type:     web.WriteType_CSV,
				Database: "",
			},
			expectedError: true,
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			_, closeLW, err := newListWriterCSV(ctx, tcase.writer, csvConfig{
				filename: tcase.filename,
			})

			if tcase.expectedError && err == nil {
				t.Fatalf("expected error, but got nil")
			}

			if closeLW != nil {
				defer closeLW()
			}
		})
	}
}

func TestHTTPRequestBuilder(t *testing.T) {
	t.Parallel()

	var testError1 = errors.New("test error 1")
	var testError2 = errors.New("test error 2")

	t.Run("close", func(t *testing.T) {
		t.Parallel()

		for _, tcase := range []struct {
			name          string
			closers       []func()
			errs          []error
			expectedError error
		}{
			{
				name: "no closers",
			},
			{
				name: "one closer",
				closers: []func(){
					func() {},
				},
			},
			{
				name: "multiple closers",
				closers: []func(){
					func() {},
					func() {},
					func() {},
				},
			},
			{
				name: "one closer with error",
				closers: []func(){
					func() {},
				},
				errs: []error{
					testError1,
				},
				expectedError: testError1,
			},
			{
				name: "multiple closers with errors",
				closers: []func(){
					func() {},
					func() {},
					func() {},
				},
				errs: []error{
					testError1,
					testError2,
				},
				expectedError: testError2,
			},
		} {
			tcase := tcase

			t.Run(tcase.name, func(t *testing.T) {
				t.Parallel()

				h := &httpRequestBuilder{
					closer: make(chan func(), len(tcase.closers)),
					errs:   make(chan error, len(tcase.errs)),
				}

				closeCounter := 0

				go func() {
					defer close(h.closer)
					defer close(h.errs)

					for _ = range tcase.closers {
						h.closer <- func() {
							closeCounter++
						}
					}

					for _, err := range tcase.errs {
						if err != nil {
							h.errs <- err
						}
					}
				}()

				if err := h.close(); err != nil {
					if !errors.Is(err, tcase.expectedError) {
						t.Fatalf("unexpected error: %v", err)
					}
				}

				if closeCounter != len(tcase.closers) {
					t.Fatalf("expected %d closers to be called, but got %d", len(tcase.closers), closeCounter)
				}
			})
		}
	})

	t.Run("setAuth", func(t *testing.T) {
		t.Parallel()

		for _, tcase := range []struct {
			name                        string
			globalAuth, localAuth, want *web.Auth
		}{
			{
				name: "no auth",
				want: nil,
			},
			{
				name:       "global auth",
				globalAuth: &web.Auth{Auth: &web.Auth_Coinbase{}},
				want:       &web.Auth{Auth: &web.Auth_Coinbase{}},
			},
			{
				name:      "local auth",
				localAuth: &web.Auth{Auth: &web.Auth_Coinbase{}},
				want:      &web.Auth{Auth: &web.Auth_Coinbase{}},
			},
			{
				name:       "global and local auth",
				globalAuth: &web.Auth{Auth: &web.Auth_Coinbase{}},
				localAuth:  &web.Auth{Auth: &web.Auth_Basic{}},
				want:       &web.Auth{Auth: &web.Auth_Basic{}},
			},
		} {
			tcase := tcase

			t.Run(tcase.name, func(t *testing.T) {
				t.Parallel()

				h := &httpRequestBuilder{
					writeRequest: &web.WriteRequest{},
				}

				if tcase.localAuth != nil {
					h.writeRequest.Auth = tcase.localAuth
				}

				h.setAuth(tcase.globalAuth)

				if !reflect.DeepEqual(h.auth, tcase.want) {
					t.Fatalf("unexpected auth: %v", h.auth)
				}
			})
		}
	})

	t.Run("setWriters", func(t *testing.T) {
		t.Parallel()

		for _, tcase := range []struct {
			name                              string
			globalWriters, localWriters, want []*web.Writer
		}{
			{
				name: "no writers",
				want: nil,
			},
			{
				name: "global writers",
				globalWriters: []*web.Writer{
					{Type: web.WriteType_MONGO},
				},
				want: []*web.Writer{
					{Type: web.WriteType_MONGO},
				},
			},
			{
				name: "local writers",
				localWriters: []*web.Writer{
					{Type: web.WriteType_CSV},
				},
				want: []*web.Writer{
					{Type: web.WriteType_CSV},
				},
			},
			{
				name: "global and local writers",
				globalWriters: []*web.Writer{
					{Type: web.WriteType_MONGO},
				},
				localWriters: []*web.Writer{
					{Type: web.WriteType_CSV},
				},
				want: []*web.Writer{
					{Type: web.WriteType_CSV},
				},
			},
		} {
			tcase := tcase

			t.Run(tcase.name, func(t *testing.T) {
				t.Parallel()

				h := &httpRequestBuilder{
					writeRequest: &web.WriteRequest{},
				}

				if tcase.localWriters != nil {
					h.writeRequest.Writers = tcase.localWriters
				}

				h.setWriters(tcase.globalWriters)

				if !reflect.DeepEqual(h.writers, tcase.want) {
					t.Fatalf("unexpected writers: %v", h.writers)
				}
			})
		}
	})

	t.Run("build", func(t *testing.T) {
		t.Parallel()

		for _, tcase := range []struct {
			name          string
			wreq          *web.WriteRequest
			expectedError error
		}{
			{
				name: "no service",
				wreq: nil,
			},
			{
				name: "no write request",
				wreq: nil,
			},
			{
				name: "no writers",
				wreq: &web.WriteRequest{},
			},
		} {
			tcase := tcase

			t.Run(tcase.name, func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()

				// Create a gidari service and add it to the
				// builder.
				svc, _ := gidari.NewService(ctx)
				bldr := &httpRequestBuilder{
					gidariService: svc,
					writeRequest:  tcase.wreq,
				}

				<-bldr.build()
				if !errors.Is(bldr.close(), tcase.expectedError) {
					t.Fatalf("unexpected error: %v", bldr.close())
				}

				if tcase.expectedError != nil {
					return
				}
			})
		}
	})
}
