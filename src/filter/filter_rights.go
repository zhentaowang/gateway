package filter

import (
    "errors"
    "net/http"
    "io/ioutil"
    "log"
    "code.aliyun.com/wyunshare/thrift-server/gen-go/server"
    "code.aliyun.com/wyunshare/thrift-server"
    "github.com/bitly/go-simplejson"
    "gateway/src/util"
)

/*
验证用户权限的filter
 */
type RightsFilter struct{
    BaseFilter
}

func newRightsFilter() Filter {
    return &RightsFilter{}
}

func (v RightsFilter) Name() string {
    return FilterRights
}

// Pre pre filter, before proxy request
func (v RightsFilter) Pre(c Context) (statusCode int, err error) {
    // 检查用户是否有权限
    conf := util.GetConfigCenterInstance()
    accessToken := c.GetProxyOuterRequest().URI().QueryArgs().Peek("access_token")

    log.Println("accessToken="+string(accessToken))
    res, err := http.Get(conf.ConfProperties["oauth_center"]["oauth_addr"]+"/user/getUser?access_token="+ string(accessToken))
    body, _ := ioutil.ReadAll(res.Body)
    log.Println("body="+string(body))

    thriftreq := server.NewRequest()
    thriftreq.ServiceName = "PermissionValidate"
    thriftreq.Operation = "validate"
    thriftreq.ParamJSON = body

    Pool := thriftserver.GetPool(conf.ConfProperties["oauth_center"]["permission_thrift"])
    pooledClient, err := Pool.Get()
    if err != nil {
        log.Println("Thrift pool get client error")
        return 500,err
    }

    defer Pool.Put(pooledClient, false)
    rawClient, ok := pooledClient.(*server.MyServiceClient)
    if !ok {
        log.Println("convert to raw client failed")
        return 503,errors.New("convert to raw client failed")
    }

    thriftres, err := rawClient.Send(thriftreq)
    if err != nil {
        log.Println("处理thrift请求失败  ")
        return http.StatusInternalServerError,err
    }

    js, err := simplejson.NewJson(thriftres.ResponseJSON)
    thriftresValue := js.Get("success")
    Checked ,_:= thriftresValue.Bool()

    if Checked {
        return http.StatusOK,nil
    }

    return http.StatusForbidden,errors.New("you don't have permission")

}