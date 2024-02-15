package filters

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type middleware func(next http.HandlerFunc) http.HandlerFunc
type ErrorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// validateJSON checks that the POST body contains a valid JSON object.
func validateJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonBytes, err := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(jsonBytes))
		if len(jsonBytes) > 0 {
			if !json.Valid(jsonBytes) || err != nil {
				errorBody, _ := json.Marshal(ErrorMessage{Message: "INVALID_JSON", Code: 0})

				w.WriteHeader(http.StatusForbidden)
				_, err := w.Write(errorBody)

				if err != nil {
					log.Printf("Error writing response body %s", err)
					return
				}
			}

			next.ServeHTTP(w, r)
		}
	}
}

func relay(w http.ResponseWriter, r *http.Request) {
	message := "Patroneos cannot receive fail2ban relay requests when running in filter mode. Please check your config."
	log.Printf("%s", message)

	errorBody, _ := json.Marshal(ErrorMessage{Message: message, Code: 403})

	w.WriteHeader(http.StatusForbidden)
	_, err := w.Write(errorBody)

	if err != nil {
		log.Printf("Error writing response body %s", err)
		return
	}
}

func chainMiddleware(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

func AddFilterHandlers(mux *http.ServeMux) {
	middlewareChain := chainMiddleware(
		validateJSON,
	)
	mux.HandleFunc("/", middlewareChain(relay))
}

func addLogHandlers() {

}
