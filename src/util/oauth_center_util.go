package util

import (
	"net/http"
	"strings"
	"io/ioutil"
	"log"
	"github.com/bitly/go-simplejson"
)

func GetToken(username string , password string)  string{

	conf := GetConfigCenterInstance()
	grant_type := conf.ConfProperties["oauth_center"]["grant_type"]
	client_id := conf.ConfProperties["oauth_center"]["client_id"]
	client_secret := conf.ConfProperties["oauth_center"]["client_secret"]

	post := "{\"username\":"+"\""+username+"\""+"," +
		"\"password\":"+"\""+password+"\""+"," +
		"\"grant_type\":"+"\""+grant_type+"\""+"," +
		"\"client_id\":"+"\""+client_id+"\""+"," +
		"\"client_secret\":"+"\""+client_secret+"\""+"}"


	//resp, err := http.Post("https://front.zhiweicloud.com/oauth/access_token",
	//	"application/json",strings.NewReader(post))
	resp, err := http.Post("http://oauth-center.platform:443/oauth/access_token",
		"application/json",strings.NewReader(post))
	if err != nil {
		log.Printf("get access_token error",err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(err.Error())
	}

	js, err := simplejson.NewJson(body)
	if err != nil {
		log.Printf(err.Error())
	}

	js_token := js.Get("access_token")

	token, _ := js_token.String()

	return token
}

func GetPermission(username string,password string,permission string)  (bool,int){
	access_token := GetToken(username,password)
	resp, err := http.Get("http://oauth-center.platform:443/user/permission?access_token="+access_token+"&permission="+permission)
	//resp, err := http.Get("https://front.zhiweicloud.com/user/permission?access_token="+access_token+"&permission="+permission)
	if err != nil {
		log.Printf(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(err.Error())
	}

	js, err := simplejson.NewJson(body)
	if err != nil {
		log.Printf(err.Error())
	}

	js_allowed := js.Get("allowed")
	js_status := js.Get("status")

	allowed, _ := js_allowed.Bool()
	status, _ := js_status.Int()

	return allowed,status
}