package common

import (
	"fmt"
	"third/raven-go"
)

var SentryClient *raven.Client

func InitSentryClient() error {
	var err error
	fmt.Println(Config.SentryUrl)
	SentryClient, err = raven.NewWithTags(Config.SentryUrl, map[string]string{"blast": "test"})
	if nil != err {
		Logger.Error("init sentry client err")
		return err
	}
	return nil
}

func InitSentryClientWithUrl(url string) error {
	var err error
	fmt.Println(url)
	SentryClient, err = raven.NewWithTags(url, map[string]string{"blast": "test"})
	if nil != err {
		Logger.Error("init sentry client err")
		return err
	}
	return nil
}

func trace() *raven.Stacktrace {
	return raven.NewStacktrace(2, 2, nil)
}

func CheckError(err error) {
	if nil == SentryClient {
		return
	}
	var in_err error
	packet := raven.NewPacket(err.Error(), raven.NewException(err, trace()))

	eventID, ch := SentryClient.Capture(packet, nil)
	in_err = <-ch
	message := fmt.Sprintf("Error event with id %s,%v", eventID, in_err)
	Logger.Error(message)
}
