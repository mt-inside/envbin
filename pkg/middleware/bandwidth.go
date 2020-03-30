package middleware

import (
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mxk/go-flowrate/flowrate"
	"log"
	"net/http"
)

type slowResponseWriter struct {
	rw   http.ResponseWriter
	fw   *flowrate.Writer
	oldBw   int64
}

func newSlowResponseWriter(rw http.ResponseWriter) slowResponseWriter {
	bandwidth:=data.GetBandwidth()
	fw := flowrate.NewWriter(rw, bandwidth)
	return slowResponseWriter{rw, fw, bandwidth}
	//defer fw.Close()
}

func (sr slowResponseWriter) Header() http.Header {
	return sr.rw.Header()
}

func (sr slowResponseWriter) Write(b []byte) (written int, err error) {
	if sr.oldBw != data.GetBandwidth() {
		// FIXME: not thread safe
		sr.oldBw = data.GetBandwidth()
		sr.fw.SetLimit(sr.oldBw)
		log.Println("adjusted writer bw to ", sr.oldBw)
	}
	written, err = sr.fw.Write(b)
	sr.rw.(http.Flusher).Flush()
	return 0, nil
}

func (sr slowResponseWriter) WriteHeader(statusCode int) {
	sr.rw.WriteHeader(statusCode)
}

func bandwidthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(newSlowResponseWriter(w), r)
	})
}