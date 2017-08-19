package thrift

import (
	"code.aliyun.com/wyunshare/thrift-server/gen-go/server"
	"code.aliyun.com/wyunshare/thrift-server"
	"log"
	"strings"
	"bytes"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	"code.aliyun.com/wyunshare/thrift-server/pool"
	"gateway/src/util"
	"time"
	"strconv"
)

type BusinessServiceImpl struct {
}

// 通过 BusinessServiceImpl 实现 IBusinessService 接口的 Send 方法，从而实现 IBusinessService 接口
func (msi *BusinessServiceImpl) Handle(operation string, paramJSON []byte) (*server.Response, error) {

	defer util.ErrHandle()

	log.Println("处理thrift,operation="+operation+"  paramJSON="+string(paramJSON))

	if operation == "" {
		log.Panic("Operation is null!")
	}

	HandleInfo := new(util.InfoCount)

	HandleInfo.RequestUrl = operation
	HandleInfo.RequestContent = string(paramJSON)

	startTime := time.Now().UnixNano()

	var pooled *pool.Pool
	addr := strings.Split(operation, "/")
	if len(addr)<3 {
		log.Panic("Operation格式为/uri/operation")
	}

	buffer := bytes.NewBufferString("")

	for _, s := range addr[1:len(addr)-1] {
		buffer.WriteString("/"+s)
	}

	conf := util.GetConfigCenterInstance()

	var MysqlUrl string = conf.ConfProperties["jdbc"]["db_username"] + ":" + conf.ConfProperties["jdbc"]["db_password"] + "@tcp(" + conf.ConfProperties["jdbc"]["db_host"] + ")/" +
		conf.ConfProperties["jdbc"]["db_name"] + "?charset=utf8"

	//print("MysqlUrl    " + MysqlUrl)
	Engine, _ := xorm.NewEngine("mysql", MysqlUrl)

	sql := "select service.name,service.namespace,service.port ,api.service_provider_name from service,api where api.service_id = service.service_id and api.uri=?"
	results, err := Engine.Query(sql,buffer.String())
	if err != nil {
		log.Panic("thrift从数据库获取Service失败 , Operation格式为/uri/operation ",err)
	}

	if len(results)!=0 {

		if len(string(results[0]["namespace"]))==0 {
			pooled = thriftserver.GetPool(string(results[0]["name"])  + ":" + string(results[0]["port"]))
		} else {
			pooled = thriftserver.GetPool(string(results[0]["name"]) + "." + string(results[0]["namespace"]) + ":" + string(results[0]["port"]))
		}
		client, err := pooled.Get()
		if err != nil {
			log.Panic("Thrift pool get client error", err)
		}

		defer pooled.Put(client, false)

		rawClient, ok := client.(*server.MyServiceClient)
		if !ok {
			log.Panic("convert to raw client failed")
		}

		req := server.NewRequest()

		req.ServiceName = string(results[0]["service_provider_name"])
		req.Operation = operation
		req.ParamJSON = paramJSON

		res, err := rawClient.Send(req)

		endTime := time.Now().UnixNano()
		HandleInfo.UsedTime = endTime - startTime

		if res != nil {
                    HandleInfo.ResponseContent = "ResponseCode="+strconv.FormatInt(int64(res.ResponeCode),10)+"  content="+string(res.ResponseJSON)
                    log.Println("结束处理thrift,response="+res.String())
		} else {
                    HandleInfo.ResponseContent = "返回了空结果"
                    log.Println("结束处理thrift,response=空")
                }

		util.SendToKafka(HandleInfo,"kafka_topic")

		if err != nil {
			log.Println(err)
		}

		return res, err
	}

	log.Println("没有查询到thrift相关服务")

	return nil,nil
}
