package model

import (
    "encoding/base64"
    "gateway/src/thrift/gen-go/server"
    "git.apache.org/thrift.git/lib/go/thrift"
    "fmt"
    "time"
    "gateway/src/config"
    pool "gateway/src/thrift"
    "log"
)

type Service struct {
    ServiceId int    `json:"service_id"`
    Namespace string `json:"namespace"`
    Name string      `json:"name"`
    Port string      `json:"port"`
    Protocol string  `json:"protocol"`
    Pool *pool.Pool         `json:"-"`
}

func (s *Service) init(r *RouteTable) error {

    //s.Pool = thrift_client_pool.NewChannelClientPool(maxIdle, 0, servers, 0, time.Duration(timeoutMs)*time.Millisecond,
    //    func(openedSocket thrift.TTransport) thrift_client_pool.Client {
    //        transport := r.transportFactory.GetTransport(openedSocket)
    //        return server.NewMyServiceClientFactory(transport, r.protocolFactory)
    //    },
    //)

    // client
    s.Pool = &pool.Pool{
        Dial: func() (interface{}, error) {
            sock, err := thrift.NewTSocket(s.GetHost())  // client端不设置超时
            if err != nil {
                log.Printf("thrift.NewTSocketTimeout(%s) error(%v)", s.GetHost(), err)
                return nil, err
            }
            tf := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
            client := server.NewMyServiceClientFactory(tf.GetTransport(sock), r.protocolFactory)
            if err = client.Transport.Open(); err != nil {
                log.Printf("client.Transport.Open() error(%v)", err)
                return nil, err
            }
            return client, nil
        },
        Close: func(v interface{}) error {
            v.(*server.MyServiceClient).Transport.Close()
            return nil
        },
        MaxActive:   config.TConfig.MaxConns,
        MaxIdle:     config.TConfig.MaxIdle,
        IdleTimeout: time.Duration(config.TConfig.MaxIdleConnDuration),
    }

    return nil
}

func (s *Service) getKey() string {
    key := fmt.Sprintf("%s-%s-%s", s.Namespace, s.Name, s.Port)
    return base64.RawURLEncoding.EncodeToString([]byte(key))
}

func (s *Service) GetHost() string {
    if len(s.Namespace) > 0 {
        return fmt.Sprintf("%s.%s:%s",s.Name, s.Namespace,  s.Port)
    } else {
        return fmt.Sprintf("%s:%s", s.Name, s.Port)
    }
}
