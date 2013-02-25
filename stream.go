package streamango

import (
	"github.com/paulbellamy/mango"
	"log"
	"net/http"
)

func streamerapp(streamer http.HandlerFunc) mango.App {
	return func(env mango.Env) (status mango.Status, headers mango.Headers, body mango.Body) {
		w := env["mango.writer"].(http.ResponseWriter)
		ifilter, _ := env["mango.bodyfilter"]

		env["streamango.streaming"] = true

		var bw http.ResponseWriter
		if ifilter != nil {
			var err error
			filter := ifilter.(BodyFilter)
			bw, err = newwriter(w, env, filter)
			if err != nil {
				log.Println("streamango FilteredFunc newwriter: ", err)
				return
			}
			defer filter.Flush()
		} else {
			bw = w
		}

		streamer(bw, env.Request().Request)
		return
	}
}

func HandlerFunc(stack *mango.Stack, streamer http.HandlerFunc) http.HandlerFunc {
	return FilteredFunc(stack, streamer, nil)
}
