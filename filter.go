package streamango

import (
	"github.com/paulbellamy/mango"
	"net/http"
)

type filteredResponseWriter struct {
	filter BodyFilter
	rw     http.ResponseWriter
}

func (w *filteredResponseWriter) Header() (h http.Header) {
	h = w.rw.Header()
	h.Del("Content-Length")
	return h
}

func (w *filteredResponseWriter) Write(p []byte) (int, error) {
	return w.filter.Write(p)
}

func (w *filteredResponseWriter) WriteHeader(h int) {
	w.rw.WriteHeader(h)
}

func newwriter(w http.ResponseWriter, env mango.Env, filter BodyFilter) (brw *filteredResponseWriter, err error) {
	brw = &filteredResponseWriter{filter: filter, rw: w}
	err = filter.Init(env, w)
	return
}

type BodyFilter interface {
	Init(env mango.Env, w http.ResponseWriter) (err error)
	Write(p []byte) (int, error)
	Flush() error
}

func FilteredFunc(stack *mango.Stack, streamer http.HandlerFunc, filter BodyFilter) http.HandlerFunc {
	compiled_app := stack.Compile(streamerapp(streamer))
	return func(w http.ResponseWriter, r *http.Request) {
		env := make(map[string]interface{})
		env["mango.request"] = &mango.Request{r}
		env["mango.bodyfilter"] = filter
		env["mango.writer"] = w

		status, headers, body := compiled_app(env)
		_, streaming := env["streamango.streaming"]
		// streaming, so don't need to do
		if streaming {
			return
		}

		for key, values := range headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(int(status))
		w.Write([]byte(body))
	}
}
