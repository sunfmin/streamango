package streamango

import (
	"github.com/paulbellamy/mango"
	"net/http"
)

func streamerapp(streamer http.HandlerFunc) mango.App {
	return func(env mango.Env) (status mango.Status, headers mango.Headers, body mango.Body) {
		w := env["mango.writer"].(http.ResponseWriter)
		env["streamango.streaming"] = true
		streamer(w, env.Request().Request)
		return
	}
}

func HandlerFunc(stack *mango.Stack, streamer http.HandlerFunc) http.HandlerFunc {
	compiled_app := stack.Compile(streamerapp(streamer))
	return func(w http.ResponseWriter, r *http.Request) {
		env := make(map[string]interface{})
		env["mango.request"] = &mango.Request{r}
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
