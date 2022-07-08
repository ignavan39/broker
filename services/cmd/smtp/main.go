package main

import (
	"broker/pkg/logger"
	smtpSrv "broker/smtp"
)

func main() {
	logger.Init()
	smtpSrv.Start()
}
