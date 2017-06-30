package proxy

import (
    "github.com/valyala/fasthttp"
    "errors"
    "gateway/src/model"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "log"
    "strconv"
    "conf_center"
)

func (h *HttpProxy) CheckToken(req *fasthttp.Request , result *model.RouteResult) (bool, error) {
    accessToken := req.URI().QueryArgs().Peek("access_token")
    if nil ==  accessToken{
        return false, errors.New("No access token")
    }
    //res, err := h.fastHTTPClient.Do(outReq, config.TConfig.OauthHost
    conf := conf_center.New("gateway")
    conf.Init()
    res, err := http.Get(conf.ConfProperties["jdbc"]["oauth_host"] + string(accessToken))
    result.Res = &fasthttp.Response{}

    result.Res.SetStatusCode(res.StatusCode)

    if err != nil {
        log.Println(err)
        return false, err
    }
    defer res.Body.Close()

    if res.StatusCode == 200 {
        var oauthResult map[string]interface{}
        body, _ := ioutil.ReadAll(res.Body)
        json.Unmarshal(body, &oauthResult)

        // 设置user_id
        //req.Header.Add("user_id", strconv.Itoa(int(oauthResult["user_id"].(float64))))
        //req.PostArgs().Add("user_id", oauthResult["user_id"].(string))
        //req.PostArgs().Add("client_id", oauthResult["client_id"].(string))
        if clientId, ok :=  oauthResult["client_id"].(string); ok{
            req.PostArgs().Add("client_id", clientId)
        }

        if userId, ok := oauthResult["user_id"].(string); ok {
            req.PostArgs().Add("user_id", userId)
        } else if userId, ok := oauthResult["user_id"].(float64); ok {
            req.PostArgs().Add("user_id", strconv.Itoa(int(userId)))
        }

        return true, nil
    } else {
        return false, err
    }
}
