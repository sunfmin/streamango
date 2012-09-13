package tests

import (
	"fmt"
	"github.com/paulbellamy/mango"
	"github.com/sunfmin/streamango"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func stream(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s", "Felix")
}

func Authenticated(env mango.Env, app mango.App) (status mango.Status, headers mango.Headers, body mango.Body) {
	if env.Request().FormValue("Password") != "nopassword" {
		body = "Not authorized."
		status = 403
		return
	}
	status, headers, body = app(env)
	return
}

type testcase struct {
	url      string
	expected string
}

var cases = []testcase{
	{"", "Not authorized."},
	{"?Password=nopassword", "Hello Felix"},
	{"?Password=nopassword2", "Not authorized."},
}

func TestStream(t *testing.T) {
	stack := new(mango.Stack)
	stack.Middleware(Authenticated)
	ts := httptest.NewServer(streamango.HandlerFunc(stack, stream))
	defer ts.Close()

	for _, tc := range cases {

		res, err := http.Get(ts.URL + tc.url)
		if err != nil {
			t.Fatal(err)
		}
		got, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != tc.expected {
			t.Errorf("got %q, want %q", string(got), tc.expected)
		}
	}

}
