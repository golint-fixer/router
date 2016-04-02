package router

import (
	"github.com/nbio/st"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestRoutingHit(t *testing.T) {
	p := New()

	var ok bool
	p.Get("/foo/:name").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok = true
		t.Logf("%#v", r.URL.Query())
		st.Expect(t, r.URL.Query().Get(":name"), "keith")
	}))

	p.HandleHTTP(nil, newRequest("GET", "/foo/keith?a=b", nil), nil)
	if !ok {
		t.Error("handler not called")
	}
}

func TestRoutingMethodNotAllowed(t *testing.T) {
	p := New()

	var ok bool
	p.Post("/foo/:name").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok = true
	}))

	p.Put("/foo/:name").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok = true
	}))

	r := httptest.NewRecorder()
	p.HandleHTTP(r, newRequest("GET", "/foo/keith", nil), nil)

	if ok {
		t.Fatal("handler called when it should have not been allowed")
	}
	if r.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d; expected %d", r.Code, http.StatusMethodNotAllowed)
	}

	got := strings.Split(r.Header().Get("Allow"), ", ")
	sort.Strings(got)
	want := []string{"POST", "PUT"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got Allow header %v; want %v", got, want)
	}
}

// Check to make sure we don't pollute the Raw Query when we have no parameters
func TestNoParams(t *testing.T) {
	p := New()

	var ok bool
	p.Get("/foo/").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok = true
		t.Logf("%#v", r.URL.RawQuery)
		if r.URL.RawQuery != "" {
			t.Errorf("RawQuery was %q; should be empty", r.URL.RawQuery)
		}
	}))

	p.HandleHTTP(nil, newRequest("GET", "/foo/", nil), nil)
	if !ok {
		t.Error("handler not called")
	}
}

// Check to make sure we don't pollute the Raw Query when there are parameters but no pattern variables
func TestOnlyUserParams(t *testing.T) {
	p := New()

	var ok bool
	p.Get("/foo/").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok = true
		t.Logf("%#v", r.URL.RawQuery)
		if got, want := r.URL.RawQuery, "a=b"; got != want {
			t.Errorf("for RawQuery: got %q; want %q", got, want)
		}
	}))

	p.HandleHTTP(nil, newRequest("GET", "/foo/?a=b", nil), nil)
	if !ok {
		t.Error("handler not called")
	}
}

func TestImplicitRedirect(t *testing.T) {
	p := New()
	p.Get("/foo/").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	res := httptest.NewRecorder()
	p.HandleHTTP(res, newRequest("GET", "/foo", nil), nil)
	if res.Code != 200 {
		t.Errorf("got Code %d; want 200", res.Code)
	}
	if loc := res.Header().Get("Location"); loc != "" {
		t.Errorf("got %q; want %q", loc, "")
	}

	p = New()
	p.Get("/foo").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	p.Get("/foo/").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	res = httptest.NewRecorder()
	p.HandleHTTP(res, newRequest("GET", "/foo", nil), nil)
	if res.Code != 200 {
		t.Errorf("got %d; want Code 200", res.Code)
	}

	p = New()
	p.Get("/hello/:name/").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	res = httptest.NewRecorder()
	p.HandleHTTP(res, newRequest("GET", "/hello/bob?a=b#f", nil), nil)
	if res.Code != 200 {
		t.Errorf("got %d; want Code 200", res.Code)
	}
}

func TestNotFound(t *testing.T) {
	p := New()
	p.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(123)
	})
	p.Post("/bar").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	for _, path := range []string{"/foo", "/bar"} {
		res := httptest.NewRecorder()
		p.HandleHTTP(res, newRequest("GET", path, nil), nil)
		if res.Code != 123 {
			t.Errorf("for path %q: got code %d; want 123", path, res.Code)
		}
	}
}

func TestMethodPatch(t *testing.T) {
	p := New()
	p.Patch("/foo/bar").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// Test to see if we get a 405 Method Not Allowed errors from trying to
	// issue a GET request to a handler that only supports the PATCH method.
	res := httptest.NewRecorder()
	res.Code = http.StatusMethodNotAllowed
	p.HandleHTTP(res, newRequest("GET", "/foo/bar", nil), nil)
	if res.Code != http.StatusMethodNotAllowed {
		t.Errorf("got Code %d; want 405", res.Code)
	}

	// Now, test to see if we get a 200 OK from issuing a PATCH request to
	// the same handler.
	res = httptest.NewRecorder()
	p.HandleHTTP(res, newRequest("PATCH", "/foo/bar", nil), nil)
	if res.Code != http.StatusOK {
		t.Errorf("Expected code %d, got %d", http.StatusOK, res.Code)
	}
}

func BenchmarkPatternMatching(b *testing.B) {
	p := New()
	p.Get("/hello/:name").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		r := newRequest("GET", "/hello/blake", nil)
		b.StartTimer()
		p.HandleHTTP(nil, r, nil)
	}
}

func newRequest(method, urlStr string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		panic(err)
	}
	return req
}
