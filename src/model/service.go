package model

import (
    "encoding/base64"
    "fmt"
    "code.aliyun.com/wyunshare/thrift-server/pool"
    "code.aliyun.com/wyunshare/thrift-server"
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
    s.Pool = thriftserver.GetPool(s.GetHost())
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
