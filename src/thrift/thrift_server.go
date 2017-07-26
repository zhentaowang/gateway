package thrift

import (
	"code.aliyun.com/wyunshare/thrift-server"
	"gateway/src/util"
)

func StartThriftServer()  {
	//thrift服务
	util.SetLogFlag()
	const thrift_server_address = "0.0.0.0" // 0.0.0.0 表示监听所有端口
	const thrift_server_port = "8889"

	thriftserver.StartSingleServer(thrift_server_address, thrift_server_port, "thriftServer", &BusinessServiceImpl{})
}
