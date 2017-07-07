package thrift

import (
	"code.aliyun.com/wyunshare/thrift-server/gen-go/server"
	"code.aliyun.com/wyunshare/thrift-server"
	"log"
	"strings"
	"bytes"
	"github.com/go-xorm/xorm"
	"conf_center"
	_ "github.com/go-sql-driver/mysql"
)

type BusinessServiceImpl struct {
}

// 通过 BusinessServiceImpl 实现 IBusinessService 接口的 Send 方法，从而实现 IBusinessService 接口
func (msi *BusinessServiceImpl) Handle(request *server.Request) (r *server.Response, err error) {
	addr := strings.Split(request.Operation, "/")

	buffer := bytes.NewBufferString("")

	for _, s := range addr[1:len(addr)-1] {
		buffer.WriteString("/"+s)
	}

	conf := conf_center.New("gateway")
	conf.Init()

	var MysqlUrl string = conf.ConfProperties["jdbc"]["db_username"] + ":" + conf.ConfProperties["jdbc"]["db_password"] + "@tcp(" + conf.ConfProperties["jdbc"]["db_host"] + ")/" +
		conf.ConfProperties["jdbc"]["db_name"] + "?charset=utf8"

	//print("MysqlUrl    " + MysqlUrl)
	Engine, _ := xorm.NewEngine("mysql", MysqlUrl)

	sql := "select service.name,service.namespace,service.port from service,api where api.service_id = service.service_id and api.uri=?"
	results, err := Engine.Query(sql,buffer.String())

	pooled:= thriftserver.GetPool(string(results["name"])+"."+string(results["namespace"])+  ":"+string(results["port"]))
	client, err := pooled.Get()
	if err != nil {
		log.Println("Thrift pool get client error", err)
		return
	}

	defer pooled.Put(client, false)

	rawClient, ok := client.(*server.MyServiceClient)
	if !ok {
		log.Println("convert to raw client failed")
		return
	}

	res, err := rawClient.Send(request)
	if err != nil {
		log.Println(err)
		return
	}

	return res, err
}
