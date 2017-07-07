package thrift

import (
	"code.aliyun.com/wyunshare/thrift-server/business"
	"code.aliyun.com/wyunshare/thrift-server/processor"
	"code.aliyun.com/wyunshare/thrift-server"
)

func StartThriftServer()  {
	//thrift服务
	const thrift_server_address = "0.0.0.0" // 0.0.0.0 表示监听所有端口
	const thrift_server_port = "8889"


	businessServiceMap := &business.BusinessServiceMap{
		ServiceMap: make(map[string]business.IBusinessService),
	}
	businessServiceMap.RegisterService("businessService", BusinessServiceImpl{})

	wyunServiceImpl := processor.WyunServiceImpl{
		BusinessServiceMap: businessServiceMap,
	}

	thriftserver.StartServer(thrift_server_address, thrift_server_port, &wyunServiceImpl)

}
