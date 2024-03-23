package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cast"

	"github.com/armiariyan/logger"
)

type Profile struct {
	Username           string          `json:"name" mask:"name"`
	Email              string          `json:"email" mask:"email"`
	PhoneNumber        string          `json:"phone" mask:"phone"`
	NIK                string          `json:"nik" mask:"any"`
	Base64FileByteMask []byte          `json:"base64byte_with_mask" mask:""`
	BirthTime          time.Time       `json:"birth_time"`
	FatherName         string          `json:"father_name"`
	Pin                string          `json:"pin" mask:"pin"`
	Base64File         string          `json:"file" mask:"base64"`
	ProfileStruct      ProfileStruct   `json:"profile_struct"`
	ProfileSlice       []ProfileStruct `json:"profile_slice"`
}

type ProfileStruct struct {
	Username       string `json:"name" mask:"name"`
	Email          string `json:"email" mask:"email"`
	PhoneNumber    string `json:"phone" mask:"phone"`
	NIK            string `json:"nik" mask:"any"`
	Base64FileByte []byte
	BirthTime      time.Time `json:"birth_time"`
	FatherName     string    `json:"father_name"`
	Pin            string    `json:"pin" mask:"pin"`
	Base64File     string    `json:"file" mask:"base64"`
}

func main() {
	// this needed to make sure log file reside in current path
	dir, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error get current path %s", err.Error())
		return
	}

	opt := logger.Options{
		Name: "Sample App",
		SysOptions: logger.OptionsLogger{
			Type: "file",
			OptionsFile: logger.OptionsFile{
				Stdout:       false,
				FileLocation: fmt.Sprintf("%s/%s/sys", dir, "tmp"),
				FileMaxAge:   time.Millisecond,
				Mask:         false,
			},
		},
		TdrOptions: logger.OptionsLogger{
			Type: "file",
			OptionsFile: logger.OptionsFile{
				Stdout:       false,
				FileLocation: fmt.Sprintf("%s/%s/tdr", dir, "tmp"),
				FileMaxAge:   time.Millisecond,
				Mask:         false,
			},
		},
	}

	// Data Testing
	request := Profile{
		Username:           "linkaja",
		Email:              "linkaja01@linkaja.com",
		PhoneNumber:        "08245678912",
		NIK:                "1234567890123456",
		Base64FileByteMask: []byte("tes data byte"),
		BirthTime:          time.Now(),
		FatherName:         "Bambang Pamungkas",
		Pin:                "123456",
		Base64File:         "data:file/pd;pdf,pdf====",
		ProfileStruct: ProfileStruct{
			Username:       "linkaja",
			Email:          "linkaja01@linkaja.com",
			PhoneNumber:    "08245678912",
			NIK:            "1234567890123456",
			Base64FileByte: []byte("tes data byte"),
			BirthTime:      time.Now(),
			FatherName:     "Bambang Pamungkas",
			Pin:            "123456",
			Base64File:     "data:file/pd;pdf,pdf====",
		},
		ProfileSlice: []ProfileStruct{
			{
				Username:       "linkaja",
				Email:          "linkaja01@linkaja.com",
				PhoneNumber:    "08245678912",
				NIK:            "1234567890123456",
				Base64FileByte: []byte("tes data byte"),
				BirthTime:      time.Now(),
				FatherName:     "Bambang Pamungkas",
				Pin:            "123456",
				Base64File:     "data:file/pd;pdf,pdf====",
			}, {
				Username:       "linkaja",
				Email:          "linkaja01@linkaja.com",
				PhoneNumber:    "08245678912",
				NIK:            "1234567890123456",
				Base64FileByte: []byte("tes data byte"),
				BirthTime:      time.Now(),
				FatherName:     "Bambang Pamungkas",
				Pin:            "123456",
				Base64File:     "data:file/pd;pdf,pdf====",
			},
		},
	}

	ctx := context.Background()
	log := logger.SetupLoggerCombine(opt)
	defer log.Close()

	log.TDR(ctx, logger.LogTdrModel{})
	log.TDR(ctx, logger.LogTdrModel{
		Path:    "/this-is-path",
		Request: request,
	})
	log.Error(ctx, "error message")
	log.Info(ctx, "THREAD", formatLogs(request)...)
}

func formatLogs(message ...interface{}) (logRecord []logger.Field) {
	for index, msg := range message {
		logRecord = append(logRecord, logger.ToField("_message_"+cast.ToString(index), msg))
	}
	return
}
