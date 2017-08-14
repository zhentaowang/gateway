package util

import (
	"log"
	"time"
	"os"
	"path/filepath"
	"sync"
	"fmt"
	"runtime/debug"
	"net/http"
	"strings"
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
	log.Println("初始化日志文件")

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
                errInfo := new(InfoCount)

		v := fmt.Sprintf("ERROR!!\n%s--\n  stack \n%s", err,string(debug.Stack()))
		errInfo.ResponseContent = v

		fmt.Println("ERROR!! ",err)
		debug.PrintStack()

                SendToKafka(errInfo,"gatewayErr")

                SendToDingDing(v)
	}
}

func SendToDingDing(v string)  {
    host, _ := os.Hostname()

    conf := GetConfigCenterInstance()

    post := "{\"msgtype\": \"markdown\",\"markdown\": { \"title\": \"gateway错误详解\",\"text\":\"### <font color=red>详细信息</font>\n"+
            "<font color=green> host:"+host+"</font>\n <p><code>"+ v+"</p></code>\"}}"
    http.Post(conf.ConfProperties["kafka"]["DingDing"],
        "application/json",
        strings.NewReader(post))
}