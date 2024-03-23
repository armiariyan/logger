package logger

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatAmount(amount int64) string {
	p := message.NewPrinter(language.Indonesian)
	return "Rp" + p.Sprintf("%d", amount)
}

func SanitizeName(customerName string) string {
	if customerName == "" || customerName == "null" {
		return "-"
	}
	if strings.HasPrefix(customerName, "null ") {
		customerName = customerName[5:]
	}
	if strings.HasSuffix(customerName, " null") {
		customerName = customerName[:len(customerName)-5]
	}
	return customerName
}

func SanitizePhoneNumber(phoneNumber string) string {
	if strings.HasPrefix(phoneNumber, "628") {
		return phoneNumber
	}
	if strings.HasPrefix(phoneNumber, "08") {
		return "628" + phoneNumber[2:]
	}
	if strings.HasPrefix(phoneNumber, "8") {
		return "628" + phoneNumber[1:]
	}
	return phoneNumber
}
