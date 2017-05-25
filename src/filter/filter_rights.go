package filter

import (
    "errors"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "bytes"
    "log"
    "github.com/valyala/fasthttp"
)

var (
    ErrRightsFailure = errors.New("没有权限")
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
    c.GetOriginRequestCtx().QueryArgs().Add("url", string(c.GetOriginRequestCtx().URI().Path()))

    params := make(map[string]interface{})
    c.GetOriginRequestCtx().QueryArgs().VisitAll(func(k []byte, v []byte) {
        params[string(k)] = string(v)
    })

    c.GetProxyOuterRequest().PostArgs().VisitAll(func(k []byte, v []byte) {
        params[string(k)] = string(v)
    })

    paramString, _ := json.Marshal(params)
    log.Println(string(paramString))

    //resp, err := http.Post("http://test.iairportcloud.com/guest-permission/get-user-permission", "application/json", bytes.NewReader(paramString))
    resp, err := http.Post("http://guest-permission/guest-permission/get-user-permission", "application/json", bytes.NewReader(paramString))
    //resp, err := http.Post("http://localhost:8080/get-user-permission", "application/json", bytes.NewReader(paramString))
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    var permission map[string]string
    body, _ := ioutil.ReadAll(resp.Body)
    json.Unmarshal(body, &permission)

    if permission[params["url"].(string)] == "false" {
        return fasthttp.StatusForbidden, ErrRightsFailure
    }
    c.GetProxyOuterRequest().PostArgs().Add("role-ids", permission["roleId"])
    return
}