package httptrace

import (
	"net/http"
	"strconv"

	"github.com/DataDog/dd-trace-go/tracer"
	"github.com/DataDog/dd-trace-go/tracer/ext"
)

// TraceHandler is a handler that traces all incoming requests.
// It implements the Handler interface.
type TraceHandler struct {
	*tracer.Tracer
	http.Handler
	service string
}

// NewTraceHandler allocates and returns a new TraceHandler.
func NewTraceHandler(h http.Handler, service string, t *tracer.Tracer) *TraceHandler {
	if t == nil {
		t = tracer.DefaultTracer
	}
	t.SetServiceInfo(service, "net/http", ext.AppTypeWeb)
	return &TraceHandler{t, h, service}
}

// ServeHTTP creates a new span for each incoming request
// and pass them through the underlying handler.
func (h *TraceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// bail out if tracing isn't enabled
	if !h.Tracer.Enabled() {
		h.Handler.ServeHTTP(w, r)
		return
	}

	// create a new span
	resource := r.Method + " " + r.URL.Path
	span := h.Tracer.NewRootSpan("http.request", h.service, resource)
	defer span.Finish()

	span.Type = ext.HTTPType
	span.SetMeta(ext.HTTPMethod, r.Method)
	span.SetMeta(ext.HTTPURL, r.URL.Path)

	// pass the span through the request context
	ctx := span.Context(r.Context())
	tracedRequest := r.WithContext(ctx)

	// trace the response to get the status code
	tracedWriter := newTracedResponseWriter(w, span)

	// run the request
	h.Handler.ServeHTTP(tracedWriter, tracedRequest)
}

// tracedResponseWriter is a small wrapper around an http response writer that will
// intercept and store the status of a request.
// It implements the ResponseWriter interface.
type tracedResponseWriter struct {
	http.ResponseWriter
	span   *tracer.Span
	status int
}

func newTracedResponseWriter(w http.ResponseWriter, span *tracer.Span) *tracedResponseWriter {
	return &tracedResponseWriter{w, span, 0}
}

func (w *tracedResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func (w *tracedResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
	w.span.SetMeta(ext.HTTPCode, strconv.Itoa(status))
	if status >= 500 && status < 600 {
		w.span.Error = 1
	}
}
