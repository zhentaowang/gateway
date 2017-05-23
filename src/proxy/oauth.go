package proxy

import (
    "github.com/valyala/fasthttp"
    "errors"
    "gateway/src/config"
    "gateway/src/model"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "log"
)

func (h *HttpProxy) CheckToken(req *fasthttp.Request , result *model.RouteResult) (bool, error) {
    accessToken := req.URI().QueryArgs().Peek("access_token")
    if nil ==  accessToken{
        return false, errors.New("No access token")
    }

    //res, err := h.fastHTTPClient.Do(outReq, config.TConfig.OauthHost
    res, err := http.Get(config.TConfig.OauthHost + string(accessToken))
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
        //req.Header.Add("User-Id", strconv.Itoa(int(oauthResult["user_id"].(float64))))
        req.PostArgs().Add("user_id", oauthResult["user_id"].(string))
        req.PostArgs().Add("client_id", oauthResult["client_id"].(string))

        return true, nil
    } else {
        return false, err
    }
}
