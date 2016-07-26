package common

import (
	"fmt"
	"log"
	"third/go-logging"
	"third/raven-go"
)

var Logger *logging.Logger
var MysqlLogger *log.Logger
var backend_info_leveld logging.LeveledBackend

func InitLogger(process_name string) (*logging.Logger, error) {

	if Logger != nil {
		return nil, nil
	}

	//format_str := "%{color}%{level}:[%{time:2006-01-02 15:04:05.000}][goroutine:%{goroutinecount}][%{shortfile}]%{color:reset}[%{message}]"
	format_str := "%{color}%{level:.4s}:%{time:2006-01-02 15:04:05.000}[%{id:03x}][%{goroutineid}/%{goroutinecount}] %{shortfile}%{color:reset} %{message}"
	Logger = logging.MustGetLogger(process_name)

	sql_log_fp, err := logging.NewFileLogWriter(Config.LogDir+"/"+process_name+".log.mysql", false, 1024*1024*1024)
	if err != nil {
		fmt.Println("open file[%s.mysql] failed[%s]", Config.LogFile, err)
		return nil, err
	}

	MysqlLogger = log.New(sql_log_fp, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	info_log_fp, err := logging.NewFileLogWriter(Config.LogDir+"/"+process_name+".log", false, 1024*1024*1024)
	if err != nil {
		fmt.Println("open file[%s] failed[%s]", Config.LogFile, err)
		return nil, err
	}

	err_log_fp, err := logging.NewFileLogWriter(Config.LogDir+"/"+process_name+".log.wf", false, 1024*1024*1024)
	if err != nil {
		fmt.Println("open file[%s.wf] failed[%s]", Config.LogFile, err)
		return nil, err
	}

	backend_info := logging.NewLogBackend(info_log_fp, "", 0)
	backend_err := logging.NewLogBackend(err_log_fp, "", 0)
	format := logging.MustStringFormatter(format_str)
	backend_info_formatter := logging.NewBackendFormatter(backend_info, format)
	backend_err_formatter := logging.NewBackendFormatter(backend_err, format)

	backend_info_leveld = logging.AddModuleLevel(backend_info_formatter)
	switch Config.LogLevel {
	case "ERROR":
		backend_info_leveld.SetLevel(logging.ERROR, "")
	case "WARNING":
		backend_info_leveld.SetLevel(logging.WARNING, "")
	case "NOTICE":
		backend_info_leveld.SetLevel(logging.NOTICE, "")
	case "INFO":
		backend_info_leveld.SetLevel(logging.INFO, "")
	case "DEBUG":
		backend_info_leveld.SetLevel(logging.DEBUG, "")
	default:
		backend_info_leveld.SetLevel(logging.ERROR, "")
	}

	backend_err_leveld := logging.AddModuleLevel(backend_err_formatter)
	backend_err_leveld.SetLevel(logging.WARNING, "")

	//add sentry log author:yuanxiang
	sentry_client, err := raven.NewWithTags(Config.SentryUrl, map[string]string{"servicename": "servicename"})
	if nil != err {
		log.Fatalf("init sentry client err")
		return nil, err
	}
	sentry_err := logging.NewSentryBackend(sentry_client, logging.ERROR)
	sentry_formatter := logging.NewBackendFormatter(sentry_err, format)
	sentry_err_leveld := logging.AddModuleLevel(sentry_formatter)
	sentry_err_leveld.SetLevel(logging.ERROR, "")

	logging.SetBackend(backend_info_leveld, backend_err_leveld, sentry_err_leveld)

	return Logger, err
}

func InitLogger1(process_name, format_str string) (*logging.Logger, error) {

	if Logger != nil {
		return nil, nil
	}

	Logger = logging.MustGetLogger(process_name)

	info_log_fp, err := logging.NewFileLogWriter(Config.LogDir+"/"+process_name+".log", false, 1024*1024*1024)
	if err != nil {
		fmt.Println("open file[%s] failed[%s]", Config.LogFile, err)
		return nil, err
	}

	err_log_fp, err := logging.NewFileLogWriter(Config.LogDir+"/"+process_name+".log.wf", false, 1024*1024*1024)
	if err != nil {
		fmt.Println("open file[%s.wf] failed[%s]", Config.LogFile, err)
		return nil, err
	}

	backend_info := logging.NewLogBackend(info_log_fp, "", 0)
	backend_err := logging.NewLogBackend(err_log_fp, "", 0)
	format := logging.MustStringFormatter(format_str)
	backend_info_formatter := logging.NewBackendFormatter(backend_info, format)
	backend_err_formatter := logging.NewBackendFormatter(backend_err, format)

	backend_info_leveld = logging.AddModuleLevel(backend_info_formatter)
	switch Config.LogLevel {
	case "ERROR":
		backend_info_leveld.SetLevel(logging.ERROR, "")
	case "WARNING":
		backend_info_leveld.SetLevel(logging.WARNING, "")
	case "NOTICE":
		backend_info_leveld.SetLevel(logging.NOTICE, "")
	case "INFO":
		backend_info_leveld.SetLevel(logging.INFO, "")
	case "DEBUG":
		backend_info_leveld.SetLevel(logging.DEBUG, "")
	default:
		backend_info_leveld.SetLevel(logging.ERROR, "")
	}

	backend_err_leveld := logging.AddModuleLevel(backend_err_formatter)
	backend_err_leveld.SetLevel(logging.WARNING, "")

	logging.SetBackend(backend_info_leveld, backend_err_leveld)

	return Logger, err
}

func ChangeLogLevel(LogLevel string) {
	switch LogLevel {
	case "ERROR":
		backend_info_leveld.SetLevel(logging.ERROR, "")
	case "WARNING":
		backend_info_leveld.SetLevel(logging.WARNING, "")
	case "NOTICE":
		backend_info_leveld.SetLevel(logging.NOTICE, "")
	case "INFO":
		backend_info_leveld.SetLevel(logging.INFO, "")
	case "DEBUG":
		backend_info_leveld.SetLevel(logging.DEBUG, "")
	default:
		backend_info_leveld.SetLevel(logging.ERROR, "")
	}
}

// If use followwing functions, call UseCommonLogger first
func SetExtraCalldepth(depth int) {
	if nil != Logger {
		Logger.ExtraCalldepth = depth
	}
}

func Debugf(format string, v ...interface{}) {
	if nil != Logger {
		Logger.Debug(format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	if nil != Logger {
		Logger.Info(format, v...)
	}
}

func Noticef(format string, v ...interface{}) {
	if nil != Logger {
		Logger.Notice(format, v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if nil != Logger {
		Logger.Warning(format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if nil != Logger {
		Logger.Error(format, v...)
	}
}

// by liudan
var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func ColorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func ColorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}

func ColorForReset() string {
	return reset
}
