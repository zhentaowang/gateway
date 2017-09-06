package util

import (
	"net/http"
	"strings"
	"io/ioutil"
	"log"
	"github.com/bitly/go-simplejson"
	"strconv"
	"github.com/valyala/fasthttp"
	"errors"
	"encoding/json"
)

type OauthCenter struct {
	UserName string
}

func GetToken(username string , password string)  string{

	defer ErrHandle()
	conf := GetConfigCenterInstance()
	grant_type := conf.ConfProperties["oauth_center"]["grant_type"]
	client_id := conf.ConfProperties["oauth_center"]["client_id"]
	client_secret := conf.ConfProperties["oauth_center"]["client_secret"]

	post := "{\"username\":"+"\""+username+"\""+"," +
		"\"password\":"+"\""+password+"\""+"," +
		"\"grant_type\":"+"\""+grant_type+"\""+"," +
		"\"client_id\":"+"\""+client_id+"\""+"," +
		"\"client_secret\":"+"\""+client_secret+"\""+"}"

	resp, err := http.Post(conf.ConfProperties["oauth_center"]["oauth_addr"]+"/oauth/access_token",
		"application/json",strings.NewReader(post))
	if err != nil {
		log.Panic("get access_token error",err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err.Error())
	}

	js, err := simplejson.NewJson(body)
	if err != nil {
		log.Println(err.Error())
	}

	token := ""
	if js != nil {
		js_token := js.Get("access_token")

		token, _ = js_token.String()
	}

	return token
}

func GetPermission(username string,password string,permission string)  (bool,int){
	defer ErrHandle()
	conf := GetConfigCenterInstance()
	access_token := GetToken(username,password)
	if access_token == "" {
		return false,0
	}
	resp, err := http.Get(conf.ConfProperties["oauth_center"]["oauth_addr"]+"/user/permission?access_token="+access_token+"&permission="+permission)
	if err != nil {
		log.Panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err.Error())
	}

	js, err := simplejson.NewJson(body)
	if err != nil {
		log.Println(err.Error())
	}

	js_allowed := js.Get("allowed")
	js_status := js.Get("status")

	allowed, _ := js_allowed.Bool()
	status, _ := js_status.Int()


	return allowed,status
}


func  CheckToken(req *fasthttp.Request , response *fasthttp.Response) (bool, error) {
	accessToken := req.URI().QueryArgs().Peek("access_token")
	if nil ==  accessToken{
		return false, errors.New("No access token")
	}
	//res, err := h.fastHTTPClient.Do(outReq, config.TConfig.OauthHost
	conf := GetConfigCenterInstance()
	res, err := http.Get(conf.ConfProperties["oauth_center"]["oauth_addr"] +"/user/getUser?access_token="+ string(accessToken))

	if res==nil {
		res = new(http.Response)
	}

	response.SetStatusCode(res.StatusCode)

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
		body, _ := ioutil.ReadAll(res.Body)
		response.AppendBody(body)
		return false, err
	}
}
