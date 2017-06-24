package thriftserver

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
	"code.aliyun.com/wyunshare/thrift-server/pool"
	"code.aliyun.com/wyunshare/thrift-server/conf"
	"code.aliyun.com/wyunshare/thrift-server/gen-go/server"
)

func GetPool(hostPort string) (*pool.Pool) {
	configByte, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		log.Fatal(err)
	}

	conf.TConfig = conf.T{}
	err = yaml.Unmarshal(configByte, &conf.TConfig)
	if nil != err {
		log.Panic("thrift load conf error: ", err)
	}
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
