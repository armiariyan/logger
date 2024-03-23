package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/armiariyan/logger"
)

func main() {
	// this needed to make sure log file reside in current path
	dir, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error get current path %s", err.Error())
		return
	}

	fileLocation := fmt.Sprintf("%s/%s", dir, "tmp")

	log := logger.SetupLoggerCombine(logger.Options{
		Name: "test",
		SysOptions: logger.OptionsLogger{
			Type: logger.File,
			OptionsFile: logger.OptionsFile{
				Stdout:       false,
				FileLocation: fileLocation + "/sys",
				FileMaxAge:   time.Hour,
				Mask:         false,
			},
		},
		TdrOptions: logger.OptionsLogger{
			Type: logger.File,
			OptionsFile: logger.OptionsFile{
				Stdout:       false,
				FileLocation: fileLocation + "/tdr",
				FileMaxAge:   time.Hour,
				Mask:         false,
			},
		},
	})

	defer func() {
		if err := log.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error closing logger %s\n", err)
		}
	}()

	h := &handler{
		log: log, // add log dependency as usual
	}

	mux := http.NewServeMux()
	mux.Handle("/", addLogContextOnReq(log, h.helloHandler()))

	log.Info(context.Background(), "Listening on :3000...")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Error(context.Background(), err.Error())
	}
}

func addLogContextOnReq(log logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Our middleware logic goes here...
		// add logging context such as tracing id per request, app name, route, etc
		threadID := fmt.Sprint(time.Now().UnixNano())
		ctxVal := logger.Context{
			ServiceName:    "my app",
			ServiceVersion: "1",
			ServicePort:    3000,
			ThreadID:       threadID,
			JourneyID:      "",
			Tag:            "xx",
			ReqMethod:      r.Method,
			ReqURI:         r.URL.Path,
			AdditionalData: nil,
		}

		ctx := logger.InjectCtx(context.Background(), ctxVal)
		r = r.WithContext(ctx)

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("error read body req in middleware: %v", err)))
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		// record response from handler and write into real writer
		respRecorder := httptest.NewRecorder()
		next.ServeHTTP(respRecorder, r)

		respBody, err := ioutil.ReadAll(respRecorder.Body)
		if err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("error read body req in middleware: %v", err)))
			return
		}

		// if we want, we can also add trace id in every response header,
		// and then easily query in kibana by "threadID"
		w.Header().Set("Correlation-ID", threadID)
		w.WriteHeader(respRecorder.Code)
		_, _ = bytes.NewReader(respBody).WriteTo(w) // put back response body

		// add TDR logger here
		log.TDR(ctx, logger.LogTdrModel{
			SrcIP:          r.RemoteAddr,
			RespTime:       time.Since(start).Milliseconds(),
			Path:           r.URL.Path,
			Header:         r.Header,
			Request:        string(reqBody),
			Response:       string(respBody),
			Error:          "",
			ThreadID:       threadID,
			AdditionalData: map[string]interface{}{},
			ResponseCode:   fmt.Sprint(respRecorder.Code),
			Method:         r.Method,
		})
	})
}

type handler struct {
	log logger.Logger
}

func (h *handler) helloHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		time.Sleep(2 * time.Millisecond) // mimicking request process

		// now, we can get logger as usual but with request scoped data in all log, such as trace id!
		h.log.Info(ctx, "Hey we got new request!")
		h.log.Error(ctx, "This if an error")

		_, _ = writer.Write([]byte("Hello world!"))
	}
}
