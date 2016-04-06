package common

import (
	"fmt"
	"log"
	"os"
	"third/go-logging"
)

var Logger *logging.Logger
var MysqlLogger *log.Logger
var backend_info_leveld logging.LeveledBackend

func InitLogger(process_name, format_str string) (*logging.Logger, error) {

	if Logger != nil {
		return nil, nil
	}

	Logger = logging.MustGetLogger(process_name)
	//sql_log_fp, err := os.OpenFile(Config.LogDir+"/"+process_name+".log.mysql", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	sql_log_fp, err := logging.NewFileLogWriter(Config.LogDir+"/"+process_name+".log.mysql",true,2000)
	if err != nil {
		fmt.Println("open file[%s.mysql] failed[%s]", Config.LogFile, err)
		return nil, err
	}

	MysqlLogger = log.New(sql_log_fp, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	//info_log_fp, err := os.OpenFile(Config.LogDir+"/"+process_name+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	info_log_fp, err := logging.NewFileLogWriter(Config.LogDir+"/"+process_name+".log",true,2000)
	if err != nil {
		fmt.Println("open file[%s] failed[%s]", Config.LogFile, err)
		return nil, err
	}else {
		fmt.Println("open 11111111")
	}
	for i := 0; i < 10000000; i++ {
		info_log_fp.Write([]byte("adsfawefaewfawefawef"))
	}
	//err_log_fp, err := os.OpenFile(Config.LogDir+"/"+process_name+".log.wf", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	err_log_fp, err := logging.NewFileLogWriter(Config.LogDir+"/"+process_name+".log.wf",true,2000)
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
