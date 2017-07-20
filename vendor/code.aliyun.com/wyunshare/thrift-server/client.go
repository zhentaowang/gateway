package thriftserver

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"os"
	"time"
	"code.aliyun.com/wyunshare/thrift-server/pool"
	"code.aliyun.com/wyunshare/thrift-server/conf"
	"code.aliyun.com/wyunshare/thrift-server/gen-go/server"
	"strconv"
	"code.aliyun.com/wyunshare/thrift-server/utils"
)

func GetPool(hostPort string) (*pool.Pool) {

	cf := utils.GetConfigCenterInstance("thrift-server")

	conf.TConfig = conf.T{}

	conf.TConfig.MaxConnDuration, _= strconv.Atoi(cf.ConfProperties["go"]["max.conn.duration"])
	conf.TConfig.MaxConns, _= strconv.Atoi(cf.ConfProperties["go"]["max.conn.duration"])
	conf.TConfig.MaxIdle, _= strconv.Atoi(cf.ConfProperties["go"]["max.idle"])
	conf.TConfig.MaxIdleConnDuration, _= strconv.Atoi(cf.ConfProperties["go"]["max.conn.duration"])
	conf.TConfig.MaxResponseBodySize, _= strconv.Atoi(cf.ConfProperties["go"]["max.response.body.size"])
	conf.TConfig.ReadTimeout, _ = strconv.Atoi(cf.ConfProperties["go"]["read.timeout"])
	conf.TConfig.ReadTimeout, _ = strconv.Atoi(cf.ConfProperties["go"]["write.timeout"])
	conf.TConfig.ReadBufferSize, _ = strconv.Atoi(cf.ConfProperties["go"]["read.buffer.size"])
	conf.TConfig.WriteBufferSize, _ = strconv.Atoi(cf.ConfProperties["go"]["write.buffer.size"])

	// client
    return &pool.Pool{
		Dial: func() (interface{}, error) {

			transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
			protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

			transport, err := thrift.NewTSocket(hostPort)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error resolving address:", err)
				os.Exit(1)
			}

			useTransport := transportFactory.GetTransport(transport)
			client := server.NewMyServiceClientFactory(useTransport, protocolFactory)
			if err := transport.Open(); err != nil {
				fmt.Fprintln(os.Stderr, "Error opening socket to"+hostPort + " ", err)
				os.Exit(1)
			}
			return client, nil
		},
		Close: func(v interface{}) error {
			v.(*server.MyServiceClient).Transport.Close()
			return nil
		},
		MaxActive:   conf.TConfig.MaxConns,
		MaxIdle:     conf.TConfig.MaxIdle,
		IdleTimeout: time.Duration(conf.TConfig.MaxIdleConnDuration),
	}
}
