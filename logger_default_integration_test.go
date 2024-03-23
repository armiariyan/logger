package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	test "github.com/gogo/protobuf/test/example"
	"github.com/stretchr/testify/assert"
)

type protoMessage struct {
	Text string
}

func (p *protoMessage) Reset()        {}
func (p *protoMessage) ProtoMessage() {}
func (p *protoMessage) String() string {
	return proto.CompactTextString(p)
}

// testAssertionLogger hold value written by logger
type testAssertionLogger struct {
	actualData []byte
}

func (a *testAssertionLogger) Write(p []byte) (n int, err error) {
	a.actualData = p
	return len(p), nil
}

func (a *testAssertionLogger) GetActualData() []byte {
	return a.actualData
}

func (a *testAssertionLogger) Close() error {
	return nil
}

// preparation
var (
	object = Object{
		FirstName:   "Bias",
		LastName:    "Tegaralaga",
		PIN:         "123456",
		FullName:    "Bias Tegaralaga",
		PhoneNumber: "6281260695203",
		Address:     "Bandung",
	}

	mapString = map[string]string{
		"firstName": object.FirstName,
		"lastName":  object.LastName,
	}
	mapInteger = map[int]int{1: 1, 2: 2}

	mapObject = map[string]Object{
		"one": object,
		"two": object,
	}

	mapObjectPointer = map[string]*Object{
		"one": &object,
		"two": &object,
	}
	sliceString  = []string{object.FirstName, object.LastName}
	sliceInteger = []int{1, 2}

	sliceObject = []Object{object, object}

	sliceObjectPointer = []*Object{&object, &object}

	integer = 1
	float   = 123.456
)

var complexObject = Main{
	Object:             object,
	MapString:          mapString,
	MapInteger:         mapInteger,
	MapObject:          mapObject,
	MapObjectPointer:   mapObjectPointer,
	SliceString:        sliceString,
	SliceInteger:       sliceInteger,
	SliceObject:        sliceObject,
	SliceObjectPointer: sliceObjectPointer,
	Integer:            integer,
	Float:              float,
	Boolean:            true,
}

// data
const message = "log message"

var ctxValue = Context{
	ServiceName:    "my service name",
	ServiceVersion: "1",
	ServicePort:    8000,
	ThreadID:       fmt.Sprint(time.Now().UnixNano()),
	JourneyID:      fmt.Sprint(time.Now().UnixNano()),
	ChainID:        fmt.Sprint(time.Now().UnixNano()),
	Tag:            "my-tag",
	ReqMethod:      "POST",
	ReqURI:         "/my-uri",
	AdditionalData: map[string]interface{}{
		"foo": "bar",
	},
}

var ctx = InjectCtx(context.Background(), ctxValue)

var fields = []Field{
	ToField("string", "string"),
	ToField("int", 0),
	ToField("int32", 32),
	ToField("int64", 64),
	ToField("slice_string", sliceString),
	ToField("slice_int", sliceInteger),
	ToField("slice_object", sliceObject),
	ToField("slice_object_ptr", sliceObjectPointer),
	ToField("map_string", mapString),
	ToField("complex_object", complexObject),
}

// MessageAndFields holds value message and field above
type MessageAndFields struct {
	// required
	LogType string `json:"logType"`
	Level   string `json:"level"`
	Message string `json:"message"`

	// inline
	String         string            `json:"string"`
	Int            int               `json:"int"`
	Int32          int32             `json:"int32"`
	Int64          int64             `json:"int64"`
	SliceString    []string          `json:"slice_string"`
	SliceInt       []int             `json:"slice_int"`
	SliceObject    []Object          `json:"slice_object"`
	SliceObjectPtr []*Object         `json:"slice_object_ptr"`
	MapString      map[string]string `json:"map_string"`
	ComplexObject  Main              `json:"complex_object"`
}

func TestDefaultLogger_Integration(t *testing.T) {

	writer := &testAssertionLogger{}
	log, err := newLogger(
		WithLevel(DebugLevel),
		WithCustomWriter(writer),
	)

	assert.NotNil(t, log)
	assert.NoError(t, err)

	t.Run("not panic sys", func(t *testing.T) {
		type F func(ctx context.Context, message string, fields ...Field)

		type TestCase struct {
			Func  F
			Level string
		}

		testCases := []TestCase{
			{
				Func:  log.Debug,
				Level: "debug",
			},
			{
				Func:  log.Info,
				Level: "info",
			},
			{
				Func:  log.Warn,
				Level: "warn",
			},
			{
				Func:  log.Error,
				Level: "error",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.Level, func(t *testing.T) {
				tc.Func(ctx, message, fields...) // call func

				// get data from log
				actualString := writer.GetActualData()

				// ensure context exactly the same, except for additional data
				var logCtxData Context
				err = json.Unmarshal(actualString, &logCtxData)
				assert.NoError(t, err)
				assert.EqualValues(t, logCtxData, ctxValue)

				// ensure message and fields exactly the same
				var logMsgAndFields MessageAndFields
				err = json.Unmarshal(actualString, &logMsgAndFields)
				assert.NoError(t, err)
				assert.EqualValues(t, LogTypeSYS, logMsgAndFields.LogType)
				assert.EqualValues(t, tc.Level, logMsgAndFields.Level) // assert expected level
				assert.EqualValues(t, message, logMsgAndFields.Message)
				assert.EqualValues(t, "string", logMsgAndFields.String)
				assert.EqualValues(t, 0, logMsgAndFields.Int)
				assert.EqualValues(t, int32(32), logMsgAndFields.Int32)
				assert.EqualValues(t, int64(64), logMsgAndFields.Int64)
				assert.EqualValues(t, sliceString, logMsgAndFields.SliceString)
				assert.EqualValues(t, sliceInteger, logMsgAndFields.SliceInt)
				assert.EqualValues(t, sliceObject, logMsgAndFields.SliceObject)
				assert.EqualValues(t, sliceObjectPointer, logMsgAndFields.SliceObjectPtr)
				assert.EqualValues(t, complexObject, logMsgAndFields.ComplexObject)
			})
		}

	})

	t.Run("panic sys", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("level panic must be panic\n")
			}
		}()

		log.Panic(ctx, message, fields...)

		// get data from log
		actualString := writer.GetActualData()

		// ensure context exactly the same, except for additional data
		var logCtxData Context
		err = json.Unmarshal(actualString, &logCtxData)
		assert.NoError(t, err)
		assert.EqualValues(t, logCtxData, ctxValue)

		// ensure message and fields exactly the same
		var logMsgAndFields MessageAndFields
		err = json.Unmarshal(actualString, &logMsgAndFields)
		assert.NoError(t, err)
		assert.EqualValues(t, LogTypeSYS, logMsgAndFields.LogType)
		assert.EqualValues(t, "panic", logMsgAndFields.Level) // assert expected level
		assert.EqualValues(t, message, logMsgAndFields.Message)
		assert.EqualValues(t, "string", logMsgAndFields.String)
		assert.EqualValues(t, 0, logMsgAndFields.Int)
		assert.EqualValues(t, int32(32), logMsgAndFields.Int32)
		assert.EqualValues(t, int64(64), logMsgAndFields.Int64)
		assert.EqualValues(t, sliceString, logMsgAndFields.SliceString)
		assert.EqualValues(t, sliceInteger, logMsgAndFields.SliceInt)
		assert.EqualValues(t, sliceObject, logMsgAndFields.SliceObject)
		assert.EqualValues(t, sliceObjectPointer, logMsgAndFields.SliceObjectPtr)
		assert.EqualValues(t, complexObject, logMsgAndFields.ComplexObject)
	})

}

func TestDefaultLoggerTDR_Integration(t *testing.T) {

	t.Run("unmasked: nil interface", func(t *testing.T) {
		writer := &testAssertionLogger{}
		log, err := newLogger(
			WithLevel(DebugLevel),
			WithCustomWriter(writer),
		)

		assert.NotNil(t, log)
		assert.NoError(t, err)

		tdrData := LogTdrModel{
			AppName:        "nama service",
			AppVersion:     "versi",
			ThreadID:       "thread id",
			JourneyID:      "journey id",
			ChainID:        "chain id",
			Path:           "path",
			Method:         "method",
			IP:             "IP",
			Port:           1234,
			SrcIP:          "source ipd",
			RespTime:       12,
			ResponseCode:   "response code",
			Header:         nil,
			Request:        nil,
			Response:       nil,
			Error:          "error if any",
			AdditionalData: nil,
		}

		log.TDR(ctx, tdrData)

		// get data from log
		actualString := writer.GetActualData()

		// ensure context exactly the same, except for additional data
		var logCtxData Context
		err = json.Unmarshal(actualString, &logCtxData)
		assert.NoError(t, err)
		assert.EqualValues(t, logCtxData, ctxValue)

		var logTdrData LogTdrModel
		err = json.Unmarshal(actualString, &logTdrData)
		assert.NoError(t, err)
		assert.EqualValues(t, tdrData.AppName, logTdrData.AppName)
		assert.EqualValues(t, tdrData.AppVersion, logTdrData.AppVersion)
		assert.EqualValues(t, tdrData.ThreadID, logTdrData.ThreadID)
		assert.EqualValues(t, tdrData.JourneyID, logTdrData.JourneyID)
		assert.EqualValues(t, tdrData.ChainID, logTdrData.ChainID)
		assert.EqualValues(t, tdrData.Path, logTdrData.Path)
		assert.EqualValues(t, tdrData.Method, logTdrData.Method)
		assert.EqualValues(t, tdrData.IP, logTdrData.IP)
		assert.EqualValues(t, tdrData.SrcIP, logTdrData.SrcIP)
		assert.EqualValues(t, tdrData.RespTime, logTdrData.RespTime)
		assert.EqualValues(t, tdrData.ResponseCode, logTdrData.ResponseCode)

		// nil interface will convert into  map[string]interface{}{} by default
		nilDefault := map[string]interface{}{}
		assert.EqualValues(t, nilDefault, logTdrData.Header)
		assert.EqualValues(t, nilDefault, logTdrData.Request)
		assert.EqualValues(t, nilDefault, logTdrData.Response)

		assert.EqualValues(t, tdrData.Error, logTdrData.Error)
		assert.EqualValues(t, nilDefault, logTdrData.AdditionalData)
		t.Logf("%s\n", actualString)
	})

	t.Run("unmasked: contain proto message non valid json", func(t *testing.T) {
		writer := &testAssertionLogger{}
		log, err := newLogger(
			WithLevel(DebugLevel),
			WithCustomWriter(writer),
		)

		protoMsg := test.A{}
		var _ proto.Message = (*test.A)(nil)

		assert.NotNil(t, log)
		assert.NoError(t, err)

		tdrData := LogTdrModel{
			AppName:        "nama service",
			AppVersion:     "versi",
			ThreadID:       "thread id",
			JourneyID:      "journey id",
			ChainID:        "chain id",
			Path:           "path",
			Method:         "method",
			IP:             "IP",
			Port:           1234,
			SrcIP:          "source ipd",
			RespTime:       12,
			ResponseCode:   "response code",
			Header:         protoMsg,
			Request:        protoMsg,
			Response:       protoMsg,
			Error:          "error if any",
			AdditionalData: protoMsg,
		}

		log.TDR(ctx, tdrData)

		// get data from log
		actualString := writer.GetActualData()

		// ensure context exactly the same, except for additional data
		var logCtxData Context
		err = json.Unmarshal(actualString, &logCtxData)
		assert.NoError(t, err)
		assert.EqualValues(t, logCtxData, ctxValue)

		var logTdrData LogTdrModel
		err = json.Unmarshal(actualString, &logTdrData)
		assert.NoError(t, err)
		assert.EqualValues(t, tdrData.AppName, logTdrData.AppName)
		assert.EqualValues(t, tdrData.AppVersion, logTdrData.AppVersion)
		assert.EqualValues(t, tdrData.ThreadID, logTdrData.ThreadID)
		assert.EqualValues(t, tdrData.JourneyID, logTdrData.JourneyID)
		assert.EqualValues(t, tdrData.ChainID, logTdrData.ChainID)
		assert.EqualValues(t, tdrData.Path, logTdrData.Path)
		assert.EqualValues(t, tdrData.Method, logTdrData.Method)
		assert.EqualValues(t, tdrData.IP, logTdrData.IP)
		assert.EqualValues(t, tdrData.SrcIP, logTdrData.SrcIP)
		assert.EqualValues(t, tdrData.RespTime, logTdrData.RespTime)
		assert.EqualValues(t, tdrData.ResponseCode, logTdrData.ResponseCode)

		protoBytes, err := json.Marshal(protoMsg)
		assert.NotEmpty(t, protoBytes)
		assert.NoError(t, err)

		protoMap := map[string]interface{}{}
		err = json.Unmarshal(protoBytes, &protoMap)
		assert.NoError(t, err)

		assert.EqualValues(t, protoMap, logTdrData.Header)
		assert.EqualValues(t, protoMap, logTdrData.Request)
		assert.EqualValues(t, protoMap, logTdrData.Response)

		assert.EqualValues(t, tdrData.Error, logTdrData.Error)
		assert.EqualValues(t, protoMap, logTdrData.AdditionalData)
	})

	t.Run("unmasked: interface contain non-valid json string", func(t *testing.T) {
		writer := &testAssertionLogger{}
		log, err := newLogger(
			WithLevel(DebugLevel),
			WithCustomWriter(writer),
		)

		assert.NotNil(t, log)
		assert.NoError(t, err)

		tdrData := LogTdrModel{
			AppName:        "nama service",
			AppVersion:     "versi",
			ThreadID:       "thread id",
			JourneyID:      "journey id",
			ChainID:        "chain id",
			Path:           "path",
			Method:         "method",
			IP:             "IP",
			Port:           1234,
			SrcIP:          "source ipd",
			RespTime:       12,
			ResponseCode:   "response code",
			Header:         "00:Success",
			Request:        "00:Success",
			Response:       "00:Success",
			Error:          "error if any",
			AdditionalData: map[string]interface{}{},
		}

		log.TDR(ctx, tdrData)

		// get data from log
		actualString := writer.GetActualData()

		// ensure context exactly the same, except for additional data
		var logCtxData Context
		err = json.Unmarshal(actualString, &logCtxData)
		assert.NoError(t, err)
		assert.EqualValues(t, logCtxData, ctxValue)

		var logTdrData LogTdrModel
		err = json.Unmarshal(actualString, &logTdrData)
		assert.NoError(t, err)
		assert.EqualValues(t, tdrData, logTdrData)
	})

	t.Run("unmasked: interface contain valid json string", func(t *testing.T) {
		writer := &testAssertionLogger{}
		log, err := newLogger(
			WithLevel(DebugLevel),
			WithCustomWriter(writer),
		)

		assert.NotNil(t, log)
		assert.NoError(t, err)

		jsonString := `{"status":"00","message":"Request being process","data":{"mobile":"6281297191466","title":"INFO","type":"TRANSACTION_SUCCESS","pushStatus":"Success"}}`

		tdrData := LogTdrModel{
			AppName:        "nama service",
			AppVersion:     "versi",
			ThreadID:       "thread id",
			JourneyID:      "journey id",
			ChainID:        "chain id",
			Path:           "path",
			Method:         "method",
			IP:             "IP",
			Port:           1234,
			SrcIP:          "source ipd",
			RespTime:       12,
			ResponseCode:   "response code",
			Header:         jsonString,
			Request:        jsonString,
			Response:       jsonString,
			Error:          "error if any",
			AdditionalData: map[string]interface{}{},
		}

		log.TDR(ctx, tdrData)

		// get data from log
		actualString := writer.GetActualData()

		// ensure context exactly the same, except for additional data
		var logCtxData Context
		err = json.Unmarshal(actualString, &logCtxData)
		assert.NoError(t, err)
		assert.EqualValues(t, logCtxData, ctxValue)

		var logTdrData LogTdrModel
		err = json.Unmarshal(actualString, &logTdrData)
		assert.NoError(t, err)
		assert.EqualValues(t, tdrData.AppName, logTdrData.AppName)
		assert.EqualValues(t, tdrData.AppVersion, logTdrData.AppVersion)
		assert.EqualValues(t, tdrData.ThreadID, logTdrData.ThreadID)
		assert.EqualValues(t, tdrData.JourneyID, logTdrData.JourneyID)
		assert.EqualValues(t, tdrData.ChainID, logTdrData.ChainID)
		assert.EqualValues(t, tdrData.Path, logTdrData.Path)
		assert.EqualValues(t, tdrData.Method, logTdrData.Method)
		assert.EqualValues(t, tdrData.IP, logTdrData.IP)
		assert.EqualValues(t, tdrData.SrcIP, logTdrData.SrcIP)
		assert.EqualValues(t, tdrData.RespTime, logTdrData.RespTime)
		assert.EqualValues(t, tdrData.ResponseCode, logTdrData.ResponseCode)

		// special handle, valid json string will be converted as json object when marshal
		// so, unmarshal will return map
		jsonMap := map[string]interface{}{}
		err = json.Unmarshal([]byte(jsonString), &jsonMap)
		assert.NoError(t, err)

		// header also converted to map if possible, but it never masked
		assert.EqualValues(t, jsonMap, logTdrData.Header)
		assert.EqualValues(t, jsonMap, logTdrData.Request)
		assert.EqualValues(t, jsonMap, logTdrData.Response)

		assert.EqualValues(t, tdrData.Error, logTdrData.Error)
		assert.EqualValues(t, tdrData.AdditionalData, logTdrData.AdditionalData)
	})

	t.Run("masked: interface contain valid json string", func(t *testing.T) {
		writer := &testAssertionLogger{}
		log, err := newLogger(
			MaskEnabled(),
			WithLevel(DebugLevel),
			WithCustomWriter(writer),
		)

		assert.NotNil(t, log)
		assert.NoError(t, err)

		jsonString := `{"status":"00","message":"Request being process","data":{"mobile":"6281297191466","title":"INFO","type":"TRANSACTION_SUCCESS","pushStatus":"Success"}}`

		tdrData := LogTdrModel{
			AppName:        "nama service",
			AppVersion:     "versi",
			ThreadID:       "thread id",
			JourneyID:      "journey id",
			ChainID:        "chain id",
			Path:           "path",
			Method:         "method",
			IP:             "IP",
			Port:           1234,
			SrcIP:          "source ipd",
			RespTime:       12,
			ResponseCode:   "response code",
			Header:         jsonString,
			Request:        jsonString,
			Response:       jsonString,
			Error:          "error if any",
			AdditionalData: map[string]interface{}{},
		}

		log.TDR(ctx, tdrData)

		// get data from log
		actualString := writer.GetActualData()

		// ensure context exactly the same, except for additional data
		var logCtxData Context
		err = json.Unmarshal(actualString, &logCtxData)
		assert.NoError(t, err)
		assert.EqualValues(t, logCtxData, ctxValue)

		var logTdrData LogTdrModel
		err = json.Unmarshal(actualString, &logTdrData)
		assert.NoError(t, err)
		assert.EqualValues(t, tdrData.AppName, logTdrData.AppName)
		assert.EqualValues(t, tdrData.AppVersion, logTdrData.AppVersion)
		assert.EqualValues(t, tdrData.ThreadID, logTdrData.ThreadID)
		assert.EqualValues(t, tdrData.JourneyID, logTdrData.JourneyID)
		assert.EqualValues(t, tdrData.ChainID, logTdrData.ChainID)
		assert.EqualValues(t, tdrData.Path, logTdrData.Path)
		assert.EqualValues(t, tdrData.Method, logTdrData.Method)
		assert.EqualValues(t, tdrData.IP, logTdrData.IP)
		assert.EqualValues(t, tdrData.SrcIP, logTdrData.SrcIP)
		assert.EqualValues(t, tdrData.RespTime, logTdrData.RespTime)
		assert.EqualValues(t, tdrData.ResponseCode, logTdrData.ResponseCode)

		// for valid json string, it will never be masked since Map is not supported by masking function
		jsonMap := map[string]interface{}{}
		err = json.Unmarshal([]byte(jsonString), &jsonMap)
		assert.NoError(t, err)

		// header also converted to map if possible, but it never masked
		assert.EqualValues(t, jsonMap, logTdrData.Header)
		assert.EqualValues(t, jsonMap, logTdrData.Request)
		assert.EqualValues(t, jsonMap, logTdrData.Response)

		assert.EqualValues(t, tdrData.Error, logTdrData.Error)
		assert.EqualValues(t, tdrData.AdditionalData, logTdrData.AdditionalData)
	})

	t.Run("masked: struct masking", func(t *testing.T) {
		writer := &testAssertionLogger{}
		log, err := newLogger(
			MaskEnabled(),
			WithLevel(DebugLevel),
			WithCustomWriter(writer),
		)

		assert.NotNil(t, log)
		assert.NoError(t, err)

		type maskedData struct {
			Pin    string `json:"pin"    mask:"pin"`
			Name   string `json:"name"   mask:"name"`
			Phone  string `json:"phone"  mask:"phone"`
			Any    string `json:"any"    mask:"any"`
			Base64 string `json:"base64" mask:"base64"`
			Email  string `json:"email"  mask:"email"`
		}

		jsonData := maskedData{
			Pin:    "123456",
			Name:   "John Doe",
			Phone:  "6281297191466",
			Any:    "any string",
			Base64: "eyJzdGF0dXMiOiIwMCIsIm1lc3NhZ2UiOiJSZXF1ZXN0IGJlaW5nIHByb2Nlc3MiLCJkYXRhIjp7Im1vYmlsZSI6IjYyODEyOTcxOTE0NjYiLCJ0aXRsZSI6IklORk8iLCJ0eXBlIjoiVFJBTlNBQ1RJT05fU1VDQ0VTUyIsInB1c2hTdGF0dXMiOiJTdWNjZXNzIn19",
			Email:  "john.doe@example.com",
		}

		tdrData := LogTdrModel{
			AppName:        "nama service",
			AppVersion:     "versi",
			ThreadID:       "thread id",
			JourneyID:      "journey id",
			ChainID:        "chain id",
			Path:           "path",
			Method:         "method",
			IP:             "IP",
			Port:           1234,
			SrcIP:          "source ipd",
			RespTime:       12,
			ResponseCode:   "response code",
			Header:         jsonData,
			Request:        jsonData,
			Response:       jsonData,
			Error:          "error if any",
			AdditionalData: map[string]interface{}{},
		}

		log.TDR(ctx, tdrData)

		// get data from log
		actualString := writer.GetActualData()

		// ensure context exactly the same, except for additional data
		var logCtxData Context
		err = json.Unmarshal(actualString, &logCtxData)
		assert.NoError(t, err)
		assert.EqualValues(t, logCtxData, ctxValue)

		var logTdrData LogTdrModel
		err = json.Unmarshal(actualString, &logTdrData)
		assert.NoError(t, err)
		assert.EqualValues(t, tdrData.AppName, logTdrData.AppName)
		assert.EqualValues(t, tdrData.AppVersion, logTdrData.AppVersion)
		assert.EqualValues(t, tdrData.ThreadID, logTdrData.ThreadID)
		assert.EqualValues(t, tdrData.JourneyID, logTdrData.JourneyID)
		assert.EqualValues(t, tdrData.ChainID, logTdrData.ChainID)
		assert.EqualValues(t, tdrData.Path, logTdrData.Path)
		assert.EqualValues(t, tdrData.Method, logTdrData.Method)
		assert.EqualValues(t, tdrData.IP, logTdrData.IP)
		assert.EqualValues(t, tdrData.SrcIP, logTdrData.SrcIP)
		assert.EqualValues(t, tdrData.RespTime, logTdrData.RespTime)
		assert.EqualValues(t, tdrData.ResponseCode, logTdrData.ResponseCode)

		// header also ALWAYS converted to map if possible, but it never masked
		header := map[string]interface{}{}
		if str, err := json.Marshal(jsonData); err == nil {
			err = json.Unmarshal(str, &header)
			assert.NoError(t, err)
		}

		assert.EqualValues(t, header, logTdrData.Header)

		// try to mask json data and compare to what returned by log
		maskedJsonData := masking(jsonData)
		if convert, ok := maskedJsonData.(reflect.Value); ok {
			maskedJsonData = convert.Interface()
		}

		maskedJsonBytes, err := json.Marshal(maskedJsonData)
		assert.NotEmpty(t, maskedJsonBytes)
		assert.NoError(t, err)

		expectedMaskedMap := map[string]interface{}{}
		err = json.Unmarshal(maskedJsonBytes, &expectedMaskedMap)
		assert.NoError(t, err)

		assert.EqualValues(t, expectedMaskedMap, logTdrData.Request)
		assert.EqualValues(t, expectedMaskedMap, logTdrData.Response)

		assert.EqualValues(t, tdrData.Error, logTdrData.Error)
		assert.EqualValues(t, tdrData.AdditionalData, logTdrData.AdditionalData)

		// t.Logf("%s", actualString)
	})

	t.Run("unmasked: text proto message", func(t *testing.T) {
		writer := &testAssertionLogger{}
		log, err := newLogger(
			WithLevel(DebugLevel),
			WithCustomWriter(writer),
		)

		assert.NotNil(t, log)
		assert.NoError(t, err)

		pm := protoMessage{Text: "lorem ipsum sit dolor amet"}

		tdrData := LogTdrModel{
			Request: pm,
		}

		log.TDR(ctx, tdrData)

		// get data from log
		actualString := writer.GetActualData()

		// ensure context exactly the same, except for additional data
		var logCtxData Context
		err = json.Unmarshal(actualString, &logCtxData)
		assert.NoError(t, err)
		assert.EqualValues(t, logCtxData, ctxValue)

		var logTdrData LogTdrModel
		err = json.Unmarshal(actualString, &logTdrData)
		assert.NoError(t, err)

		b, err := json.Marshal(pm)
		assert.NoError(t, err)

		var data interface{}
		err = json.Unmarshal(b, &data)
		assert.NoError(t, err)

		assert.EqualValues(t, data, logTdrData.Request)
	})

	t.Run("unmasked: json proto message", func(t *testing.T) {
		writer := &testAssertionLogger{}
		log, err := newLogger(
			WithLevel(DebugLevel),
			WithCustomWriter(writer),
		)

		assert.NotNil(t, log)
		assert.NoError(t, err)

		pm := protoMessage{Text: `{"json":"lorem ipsum sit dolor amet"}`}

		tdrData := LogTdrModel{
			Request: pm,
		}

		log.TDR(ctx, tdrData)

		// get data from log
		actualString := writer.GetActualData()

		// ensure context exactly the same, except for additional data
		var logCtxData Context
		err = json.Unmarshal(actualString, &logCtxData)
		assert.NoError(t, err)
		assert.EqualValues(t, logCtxData, ctxValue)

		var logTdrData LogTdrModel
		err = json.Unmarshal(actualString, &logTdrData)
		assert.NoError(t, err)

		expectedMaskedMap := map[string]interface{}{}

		b, err := json.Marshal(pm)
		assert.NoError(t, err)

		err = json.Unmarshal(b, &expectedMaskedMap)
		assert.NoError(t, err)

		assert.EqualValues(t, expectedMaskedMap, logTdrData.Request)
	})
}
