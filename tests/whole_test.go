package tests

import (
	"github.com/sunfmin/streamango"
	"github.com/paulbellamy/mango"
	"testing"
	"net/http"
	"fmt"
)

func stream(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s", "Felix")
}

func Authenticated(env mango.Env, app mango.App) (status mango.Status, headers mango.Headers, body mango.Body) {

}

func TestStream(t *testing.T) {

}
