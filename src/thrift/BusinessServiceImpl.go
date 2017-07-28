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
)

type BusinessServiceImpl struct {
}

// 通过 BusinessServiceImpl 实现 IBusinessService 接口的 Send 方法，从而实现 IBusinessService 接口
func (msi *BusinessServiceImpl) Handle(operation string, paramJSON []byte) (*server.Response, error) {

	defer util.ErrHandle()
	var pooled *pool.Pool
	addr := strings.Split(operation, "/")

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
	if err != nil {
		log.Println(err)
	}

	return res, err
}
