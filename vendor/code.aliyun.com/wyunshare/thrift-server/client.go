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
	"gateway/src/util"
)

func GetPool(hostPort string) (*pool.Pool) {

	cf := util.GetConfigCenterInstance()

	conf.TConfig = conf.T{}

	conf.TConfig.MaxConnDuration, _= strconv.Atoi(cf.ConfProperties["jdbc"]["max_conn_duration"])
	conf.TConfig.MaxConns, _= strconv.Atoi(cf.ConfProperties["jdbc"]["max_conns"])
	conf.TConfig.MaxIdle, _= strconv.Atoi(cf.ConfProperties["jdbc"]["max_idle"])
	conf.TConfig.MaxIdleConnDuration, _= strconv.Atoi(cf.ConfProperties["jdbc"]["max_idle_conn_duration"])
	conf.TConfig.MaxResponseBodySize, _= strconv.Atoi(cf.ConfProperties["jdbc"]["max_response_body_size"])
	conf.TConfig.ReadTimeout, _ = strconv.Atoi(cf.ConfProperties["jdbc"]["read_timeout"])
	conf.TConfig.ReadTimeout, _ = strconv.Atoi(cf.ConfProperties["jdbc"]["write_timeout"])
	conf.TConfig.ReadBufferSize, _ = strconv.Atoi(cf.ConfProperties["jdbc"]["read_buffer_size"])
	conf.TConfig.WriteBufferSize, _ = strconv.Atoi(cf.ConfProperties["jdbc"]["write_buffer_size"])

	// client
    return &pool.Pool{
		Dial: func() (interface{}, error) {

			transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
			protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

			transport, err := thrift.NewTSocket(hostPort)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error resolving address:", err)
				//os.Exit(1)
			}

			useTransport := transportFactory.GetTransport(transport)
			client := server.NewMyServiceClientFactory(useTransport, protocolFactory)
			if err := transport.Open(); err != nil {
				fmt.Fprintln(os.Stderr, "Error opening socket to"+hostPort + " ", err)
				//os.Exit(1)
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
