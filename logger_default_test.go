package logger

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type CoordinateJSON struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

type LocationJSON struct {
	Address    string         `json:"string"`
	Coordinate CoordinateJSON `json:"coordinates"`
}

type NameJSON struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type DataJSON struct {
	Say      string       `json:"say"`
	Name     NameJSON     `json:"name"`
	Location LocationJSON `json:"location"`
}

type ResponseJSON struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    DataJSON `json:"data"`
}

func GenerateResponseJson() ResponseJSON {
	return ResponseJSON{
		Success: true,
		Message: "Success",
		Data: DataJSON{
			Say: "something",
			Name: NameJSON{
				First: "John",
				Last:  "Doe",
			},
			Location: LocationJSON{
				Address: "no where",
				Coordinate: CoordinateJSON{
					Latitude:  123,
					Longitude: 321,
				},
			},
		},
	}
}

func GenerateLogTDR(dataResp interface{}) LogTdrModel {
	if dataResp == nil {
		dataResp = `{"success": true, "message": null, "data": {"say": "Hello, World!", "name": {"first": "Bias", "last": "Tegaralaga"}, "location": {"address": "Jakarta", "coordinates": {"lat": 0.0, "lng": 0.1}}}}`
	}

	return LogTdrModel{
		AppName:    "testing",
		AppVersion: "",
		IP:         "127.0.0.1",
		Port:       80,
		SrcIP:      "0.0.0.0",
		RespTime:   17,
		Path:       "/v1/check/health",
		Header: http.Header{
			"a": []string{"a"},
			"b": []string{"b"},
		},
		Request:        `{"action": "hello"}`,
		Response:       dataResp,
		Error:          "",
		ThreadID:       "",
		AdditionalData: nil,
		ResponseCode:   "",
		Method:         "",
	}
}

func TestNewLogger(t *testing.T) {
	t.Run("Options error", func(t *testing.T) {
		optErr := func() Option {
			return func(_ *defaultLogger) error {
				return fmt.Errorf("error")
			}
		}

		log, err := newLogger(optErr())
		assert.Nil(t, log)
		assert.Error(t, err)
	})

	t.Run("Without any options should return success", func(t *testing.T) {
		log, err := newLogger()
		assert.NotNil(t, log)
		assert.NoError(t, err)
	})

	t.Run("Return noop logger", func(t *testing.T) {
		log, err := newLogger(OptNoop())
		assert.NotNil(t, log)
		assert.NoError(t, err)
	})
}

func loggerInstanceFile() Logger {
	dir, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error get current path %s", err.Error())
		return nil
	}

	fileLocation := fmt.Sprintf("%s/%s/test", dir, "tmp")

	fileConfig := &OptionsFile{
		Stdout:       false,
		FileLocation: fileLocation + "/sys",
		FileMaxAge:   time.Hour,
		Mask:         false,
	}

	loggerInstance, err := newLogger(WithFileOutput(fileConfig))
	if err != nil {
		panic(err)
	}

	return loggerInstance
}

func cleanDir() {
	dir, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error get current path %s", err.Error())
		return
	}

	fileLocation := fmt.Sprintf("%s/%s/test", dir, "tmp")

	err = os.RemoveAll(fileLocation)
	if err != nil {
		panic(err)
	}
}

func TestDefaultLogger_TDR(t *testing.T) {
	ctx := InjectCtx(context.Background(), Context{
		ServiceName: "xxx",
		ThreadID:    fmt.Sprint(time.Now().UnixNano()),
		Tag:         "mylog",
		ReqMethod:   "GET",
		ReqURI:      "/",
	})

	log := loggerInstanceFile()
	defer func() {
		cleanDir()
		err := log.Close()
		assert.NoError(t, err)
	}()
	for i := 0; i < 100; i++ {
		log.TDR(ctx, GenerateLogTDR(nil))
	}
}

func BenchmarkDefaultLogger_TDR(b *testing.B) {
	ctx := InjectCtx(context.Background(), Context{
		ServiceName: "xxx",
		ThreadID:    fmt.Sprint(time.Now().UnixNano()),
		Tag:         "mylog",
		ReqMethod:   "GET",
		ReqURI:      "/",
	})

	log := loggerInstanceFile()
	defer func() {
		cleanDir()
		err := log.Close()
		assert.NoError(b, err)
	}()

	for i := 0; i < b.N; i++ {
		log.TDR(ctx, GenerateLogTDR(nil))
	}
}
