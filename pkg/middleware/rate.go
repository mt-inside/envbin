package middleware

/*
* NB: Since it's not 1980, everything in the chain tries to buffer.
* We are buffer-free, honest. You can confirm this with netcat.
* To see the trickle feed with normal software, you need:
* - curl -N
* - http(ie) --stream # still line-buffers; you won't see character-by-character output
 */

//func Rate(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		log := r.Context().Value(CtxKeyLog).(logr.Logger)

//		next.ServeHTTP(newSlowResponseWriter(log, w), r)
//	})
//}

//type slowResponseWriter struct {
//	log   logr.Logger
//	rw    http.ResponseWriter
//	fw    *flowrate.Writer
//	oldBw int64
//}

//func newSlowResponseWriter(log logr.Logger, rw http.ResponseWriter) slowResponseWriter {
//	fw := flowrate.NewWriter(rw, 0)
//	return slowResponseWriter{log, rw, fw, 0}
//	//defer fw.Close()
//}

//func (sr slowResponseWriter) Header() http.Header {
//	return sr.rw.Header()
//}

//func (sr slowResponseWriter) Write(bs []byte) (written int, err error) {

//	monitorCtx, cancel := context.WithCancel(context.Background())
//	defer cancel()

//	go func() {
//		ticker := time.NewTicker(time.Second * 1)
//		for {
//			select {
//			case <-monitorCtx.Done():
//				return
//			case <-ticker.C: // by not printing until 1 second has elapsed, we don't bother with status updates for fast transfers
//				s := sr.fw.Status()
//				sr.log.V(1).Info("Slow transfer", "done", s.Progress, "remaining", s.TimeRem, "bytes/s", s.CurRate)
//			}
//		}

//	}()

//	sr.fw.SetTransferSize(int64(len(bs)))
//	defer sr.rw.(http.Flusher).Flush()
//	for i := 0; i < len(bs); i++ {
//		if sr.oldBw != viper.GetInt64("Rate") {
//			sr.oldBw = viper.GetInt64("Rate")
//			sr.fw.SetLimit(sr.oldBw)
//			sr.log.V(1).Info("adjusted writer bw", "new rate bytes/s", sr.oldBw)
//		}

//		sr.fw.Write(bs[i : i+1])
//		sr.rw.(http.Flusher).Flush()
//	}
//	written, err = 0, nil // TODO

//	return
//}

//func (sr slowResponseWriter) WriteHeader(statusCode int) {
//	sr.rw.WriteHeader(statusCode)
//}
