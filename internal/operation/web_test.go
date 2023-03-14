package operation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/alpstable/cli/api/service/web"
	"github.com/alpstable/gidari/auth"
)

func newTestWriters(vol int, stub string) []*web.Writer {
	writers := make([]*web.Writer, vol)
	for i := 0; i < vol; i++ {
		writers[i] = &web.Writer{
			Database: fmt.Sprintf("db%d-%s", i, stub),
		}
	}

	return writers
}

func newTestWebRequests(vol int, auth *web.Auth) *web.Request {
	return &web.Request{
		Auth:    auth,
		Writers: newTestWriters(vol, "wr"),
	}
}

func newTestWebWriteRequests(vol int, auth *web.Auth) *web.WriteRequest {
	return &web.WriteRequest{
		Auth:    auth,
		Writers: newTestWriters(vol, "wwr"),
	}
}

func newTestCoinbaseAuth() *web.Auth {
	return &web.Auth{Auth: &web.Auth_Coinbase{}}
}

func newTestBasicAuth() *web.Auth {
	return &web.Auth{Auth: &web.Auth_Basic{}}
}

type ctorCfg struct {
	writeType  web.WriteType
	err        error
	waitTimeMS time.Duration
}

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

			_, closeLW, err := newListWriterCSV(ctx, tcase.filename, tcase.writer)
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

	t.Run("setWriters", func(t *testing.T) {})
	t.Run("build", func(t *testing.T) {})
}

func TestHTTPServiceBuilder(t *testing.T) {
	t.Run("close", func(t *testing.T) {})
	t.Run("build", func(t *testing.T) {})
}

//func newWebWriterSetter(cfg ctorCfg) webWriterSetter {
//	waitTime := cfg.waitTimeMS * time.Millisecond
//	err := cfg.err
//
//	return func(_ context.Context, w *web.Writer, get svcWR) (<-chan func(), <-chan error) {
//		closers := make(chan func(), 1)
//		errs := make(chan error, 1)
//
//		go func() {
//			defer close(closers)
//			defer close(errs)
//
//			if wt := waitTime; wt > 0 {
//				time.Sleep(waitTime)
//			}
//
//			closers <- func() {}
//			if err != nil {
//				errs <- err
//			}
//		}()
//
//		return closers, errs
//	}
//}
//
//func newTestCtor(req *web.Request, cfgs ...ctorCfg) []webWriterSetterCtor {
//	// cstors are the mock constructors that will write to close and error
//	// channels one message each.
//	cstors := []webWriterSetterCtor{}
//
//	// writeTypeSet will keep track of all the write type values we've used.
//	// There can only be one write type constructor per setWebWriter call.
//	writeTypeSet := make(map[web.WriteType]struct{})
//
//	// Set the write types on the mock requests along with any of the
//	// writerArgs.
//	for i, cfg := range cfgs {
//		if _, ok := writeTypeSet[cfg.writeType]; !ok {
//			writeTypeSet[cfg.writeType] = struct{}{}
//			req.Writers[i].Type = cfg.writeType
//		}
//
//		// Always set this data, it will test the robustness of the
//		// setWebWriter function.
//		wsetter := newWebWriterSetter(cfg)
//		cstors = append(cstors, func(web.WriteType) (webWriterSetter, bool) {
//			return wsetter, true
//		})
//	}
//
//	return cstors
//}
//
//func TestSetWebWriter(t *testing.T) {
//	t.Parallel()
//
//	var testerr = fmt.Errorf("test error")
//	var testerr2 = fmt.Errorf("test error 2")
//
//	for _, tcase := range []struct {
//		name         string
//		ctorCfgs     []ctorCfg
//		expectedErrs []error
//	}{
//		{
//			name: "none",
//		},
//		{
//			name: "one",
//			ctorCfgs: []ctorCfg{
//				{
//					writeType:  web.WriteType_CSV,
//					err:        nil,
//					waitTimeMS: 0,
//				},
//			},
//		},
//		{
//			name: "one with error",
//			ctorCfgs: []ctorCfg{
//				{
//					writeType:  web.WriteType_CSV,
//					err:        testerr,
//					waitTimeMS: 0,
//				},
//			},
//			expectedErrs: []error{testerr},
//		},
//		{
//			name: "one with wait",
//			ctorCfgs: []ctorCfg{
//				{
//					writeType:  web.WriteType_CSV,
//					err:        nil,
//					waitTimeMS: 100,
//				},
//			},
//		},
//		{
//			name: "one with wait and error",
//			ctorCfgs: []ctorCfg{
//				{
//					writeType:  web.WriteType_CSV,
//					err:        testerr,
//					waitTimeMS: 100,
//				},
//			},
//			expectedErrs: []error{testerr},
//		},
//		{
//			name: "two different write types",
//			ctorCfgs: []ctorCfg{
//				{
//					writeType:  web.WriteType_CSV,
//					err:        nil,
//					waitTimeMS: 0,
//				},
//				{
//					writeType:  web.WriteType_MONGO,
//					err:        nil,
//					waitTimeMS: 0,
//				},
//			},
//		},
//		{
//			name: "two same write types",
//			ctorCfgs: []ctorCfg{
//				{
//					writeType:  web.WriteType_CSV,
//					err:        nil,
//					waitTimeMS: 0,
//				},
//				{
//					writeType:  web.WriteType_CSV,
//					err:        nil,
//					waitTimeMS: 0,
//				},
//			},
//		},
//		{
//			name: "two same write types where last has error",
//			ctorCfgs: []ctorCfg{
//				{
//					writeType:  web.WriteType_CSV,
//					err:        nil,
//					waitTimeMS: 0,
//				},
//				{
//					writeType:  web.WriteType_CSV,
//					err:        testerr,
//					waitTimeMS: 0,
//				},
//			},
//		},
//		{
//			name: "two same write types where first has error",
//			ctorCfgs: []ctorCfg{
//				{
//					writeType:  web.WriteType_CSV,
//					err:        testerr,
//					waitTimeMS: 0,
//				},
//				{
//					writeType:  web.WriteType_CSV,
//					err:        nil,
//					waitTimeMS: 0,
//				},
//			},
//			expectedErrs: []error{testerr},
//		},
//		{
//			name: "two same write types where both have errors",
//			ctorCfgs: []ctorCfg{
//				{
//					writeType:  web.WriteType_CSV,
//					err:        testerr,
//					waitTimeMS: 0,
//				},
//				{
//					writeType:  web.WriteType_CSV,
//					err:        testerr2,
//					waitTimeMS: 0,
//				},
//			},
//			expectedErrs: []error{testerr},
//		},
//	} {
//		tcase := tcase
//		t.Run(tcase.name, func(t *testing.T) {
//			t.Parallel()
//
//			ctx := context.Background()
//			svc, _ := gidari.NewService(ctx)
//
//			// volSet is the number of ctorCfg values with unique
//			// write types.
//			volSet := make(map[web.WriteType]struct{})
//			for _, cfg := range tcase.ctorCfgs {
//				volSet[cfg.writeType] = struct{}{}
//			}
//
//			vol := len(volSet)
//
//			// Mock one write request with multiple writers, one for
//			// each write type defined in the test case.
//			writeReq := newTestWebWriteRequests(1, newTestCoinbaseAuth())
//			writeReq.Writers = newTestWriters(vol, "wwr")
//
//			req := newTestWebRequests(vol, newTestCoinbaseAuth())
//			req.Requests = []*web.WriteRequest{writeReq}
//
//			getter := newSvcWR(svc, newWebWriteRequest(newWebRequest(req), writeReq))
//
//			cstors := newTestCtor(req, tcase.ctorCfgs...)
//			closers, errs := setWebWriter(ctx, getter, cstors...)
//
//			closersCount := 0
//			for _ = range closers {
//				closersCount++
//			}
//
//			if closersCount != vol {
//				t.Fatalf("expected %d closers, got %d", volSet, closersCount)
//			}
//
//			for err := range errs {
//				expErr := tcase.expectedErrs[0]
//				if err != nil && !errors.Is(err, expErr) {
//					t.Fatalf("expected error %v, got %v", expErr, err)
//				}
//
//				tcase.expectedErrs = tcase.expectedErrs[1:]
//			}
//		})
//	}
//}
//
//func TestWebWriteRequest(t *testing.T) {
//	t.Parallel()
//
//	t.Run("writers", func(t *testing.T) {
//		for _, tcase := range []struct {
//			name   string
//			parent *web.Request
//			in     *web.WriteRequest
//			want   []*web.Writer
//		}{
//			{
//				name: "nil",
//				want: nil,
//			},
//			{
//				name:   "empty",
//				parent: &web.Request{},
//				want:   nil,
//			},
//			{
//				name:   "child=1 and parent=1",
//				parent: newTestWebRequests(1, nil),
//				in:     newTestWebWriteRequests(1, nil),
//				want:   newTestWriters(1, "wwr"),
//			},
//			{
//				name:   "child=0 and parent=1",
//				parent: newTestWebRequests(1, nil),
//				in:     newTestWebWriteRequests(0, nil),
//				want:   newTestWriters(1, "wr"),
//			},
//			{
//				name:   "child=1 and parent=0",
//				parent: newTestWebRequests(0, nil),
//				in:     newTestWebWriteRequests(1, nil),
//				want:   newTestWriters(1, "wwr"),
//			},
//			{
//				name:   "child=0 and parent=0",
//				parent: newTestWebRequests(0, nil),
//				in:     newTestWebWriteRequests(0, nil),
//				want:   nil,
//			},
//		} {
//			tcase := tcase
//			t.Run(tcase.name, func(t *testing.T) {
//				t.Parallel()
//
//				w := newWebWriteRequest(newWebRequest(tcase.parent), tcase.in)
//
//				actual := w.writerCfgs()
//				if !reflect.DeepEqual(actual, tcase.want) {
//					t.Errorf("expected %v, got %v", tcase.want, actual)
//				}
//			})
//		}
//	})
//
//	t.Run("auth", func(t *testing.T) {
//		t.Parallel()
//
//		for _, tcase := range []struct {
//			name   string
//			parent *web.Request
//			in     *web.WriteRequest
//			want   *web.Auth
//		}{
//			{
//				name: "nil",
//				want: nil,
//			},
//			{
//				name:   "empty",
//				parent: &web.Request{},
//				want:   nil,
//			},
//			{
//				name:   "child=1 and parent=1",
//				parent: newTestWebRequests(0, newTestBasicAuth()),
//				in:     newTestWebWriteRequests(0, newTestCoinbaseAuth()),
//				want:   newTestCoinbaseAuth(),
//			},
//			{
//				name:   "child=0 and parent=1",
//				parent: newTestWebRequests(0, newTestBasicAuth()),
//				in:     newTestWebWriteRequests(0, nil),
//				want:   newTestBasicAuth(),
//			},
//			{
//				name:   "child=1 and parent=0",
//				parent: newTestWebRequests(0, nil),
//				in:     newTestWebWriteRequests(0, newTestCoinbaseAuth()),
//				want:   newTestCoinbaseAuth(),
//			},
//		} {
//			tcase := tcase
//			t.Run(tcase.name, func(t *testing.T) {
//				t.Parallel()
//
//				w := newWebWriteRequest(newWebRequest(tcase.parent), tcase.in)
//
//				actual := w.auth()
//				if !reflect.DeepEqual(actual, tcase.want) {
//					t.Errorf("expected %v, got %v", tcase.want, actual)
//				}
//			})
//		}
//	})
//}
