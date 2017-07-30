package util

import (
	"log"
	"time"
	"os"
	"path/filepath"
	"sync"
	"fmt"
)

var creLog sync.Once
var Lg log.Logger

var comLog sync.Once
var ComLog *CommonLog

type CommonLog struct {
	Info string
}


func (cg *CommonLog) Error ()  string{
	return cg.Info
}


func (cg *CommonLog) Init()  {

	pwd, _ := os.Getwd()

	lf := filepath.Join(pwd, "log",time.Now().Format("20060102") + ".log")

	logFile, err := os.OpenFile(lf,os.O_APPEND|os.O_CREATE,0666);
	if err != nil {
		log.Println("创建日志文件失败 " + err.Error());
	}

	Lg = *log.New(logFile, "", log.LstdFlags|log.Llongfile)

	log.Println("创建日志文件 " + lf + "成功！")
}


func (cg *CommonLog) GetLogger() (log.Logger,error) {

	creLog.Do(func() {
		cg.Init()
	})

	pwd, _ := os.Getwd()
	lf := filepath.Join(pwd, "log",time.Now().Format("20060102") + ".log")
	_, err := os.Stat(lf)

	if err == nil {
		log.Println("创建日志文件 "+lf+" 成功！")
		return Lg,nil
	}

	if os.IsNotExist(err) {
		cg.Init()

		return Lg,nil
	}

	Err := new(CommonLog)
	Err.Info = "创建日志文件" + lf + "失败"

	return log.Logger{},Err
}

func GetCommonLog()  log.Logger{

	comLog.Do(func() {
		ComLog = new(CommonLog)
	})
	Lg,Err := ComLog.GetLogger()
	if Err != nil {
		log.Println(Err);
	}

	return Lg
}

func SetLogFlag()  {
	log.SetFlags(log.LstdFlags|log.Llongfile)
}

func ErrHandle()  {
	if err := recover(); err != nil {
		fmt.Println("ERROR!! ",err)
	}
}