服务端使用

router.go
```go
import (
	"thrift-server/src/gen-go/server"
	"encoding/json"
)

type BusinessServiceImpl struct {
}

// 通过 BusinessServiceImpl 实现 IBusinessService 接口的 Send 方法，从而实现 IBusinessService 接口
func (msi *BusinessServiceImpl) Handle(operation string, paramJSON []byte) (r *server.Response, err error) {

	paramInput := make(map[string]interface{})
	paramInput["d"] = "a"
	json.Unmarshal(paramJSON, &paramInput)

	dd := make([]byte, 1000)
	dd, err  = json.Marshal(paramInput)

	switch operation {
	case "test":
		a := &server.Response{
			ResponeCode:  server.RESCODE__200,
			ResponseJSON: dd,
		}
		return a, nil
	default:
		return nil, nil
	}

}
```
main.go
```go

import (
	"thrift-server/src/business"
	"thrift-server/src/processor"
	"thrift-server/src/server"
	"servicego/src/router"
)

type BusinessServiceImpl struct {
}

func main() {

	const thrift_server_address = "0.0.0.0" // 0.0.0.0 表示监听所有端口
	const thrift_server_port = "8888"

	businessServiceMap := &business.BusinessServiceMap{
		ServiceMap: make(map[string]business.IBusinessService),
	}
	businessServiceMap.RegisterService("businessService", &router.BusinessServiceImpl{})

	wyunServiceImpl := processor.WyunServiceImpl{
		BusinessServiceMap: businessServiceMap,
	}

	server.StartServer(thrift_server_address, thrift_server_port, &wyunServiceImpl)

}
```
客户端使用
```go

import (
	"thrift-server/gen-go/server"
	"thrift-server"
	"log"
	"encoding/json"
)

func main() {

	const host = "localhost"
	const port = "8888"
	dict := make(map[string]string)
	dict["d"] = "a"

	service := client.Service //申明一个请求服务
	service.Host = host       // 初始化请求地址
	service.Port = port       // 初始化请求端口

	pooledClient, err := service.Pool.Get() //从连接池获取thrift client
	if err != nil {
		log.Println("Thrift pool get client error", err)
		return
	}

	defer service.Pool.Put(pooledClient, false)

	rawClient, ok := pooledClient.(*server.MyServiceClient)
	if !ok {
		log.Println("convert to raw client failed")
		return
	}

	//req := &server.Request{[]byte(dict), "name", "test"}

	req := server.NewRequest() //创建request
	req.ServiceName = "test_service"
	req.Operation = "test"

	paramInput := make(map[string]interface{})

	req.ParamJSON, _ = json.Marshal(paramInput)

	r, err := rawClient.Send(req) //发送request
}
```