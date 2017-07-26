package thriftserver

import (
    "fmt"
    "os"
    "git.apache.org/thrift.git/lib/go/thrift"
    "code.aliyun.com/wyunshare/thrift-server/gen-go/server"
    "code.aliyun.com/wyunshare/thrift-server/processor"
    "code.aliyun.com/wyunshare/thrift-server/business"
    "log"
    "gateway/src/util"
)

func StartServer(address string, port string, handler server.MyService) {

    util.SetLogFlag()
    processor := server.NewMyServiceProcessor(handler)
    serverTransport, err := thrift.NewTServerSocket(address + ":" + port)
    transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
    protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
    if err != nil {
        fmt.Println("Error!", err)
        os.Exit(1)
    }

    server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
    log.Printf("thrift server in %s", address + ":" + port)
    server.Serve()
}

func StartSingleServer(address string, port string, serviceName string, businessService business.IBusinessService) {
    businessServiceMap := &business.BusinessServiceMap{
        ServiceMap: make(map[string]business.IBusinessService),
    }
    businessServiceMap.RegisterService(serviceName, businessService)

    wyunServiceImpl := processor.WyunServiceImpl{
        BusinessServiceMap: businessServiceMap,
    }

    StartServer(address, port, &wyunServiceImpl)
}