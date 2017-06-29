package business

import "code.aliyun.com/wyunshare/thrift-server/gen-go/server"

type IBusinessService interface {
	Handle(operation string, paramJSON []byte) (r *server.Response, err error)
}
