package httptrace_test

import (
	"net/http"

	ht "github.com/DataDog/dd-trace-go/tracer/contrib/net/httptrace"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}

func Example() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	http.ListenAndServe(":8080", ht.NewTraceHandler(mux, "web-service", nil))
}
