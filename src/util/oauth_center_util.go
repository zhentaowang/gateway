package util

import (
	"net/http"
	"strings"
	"io/ioutil"
	"log"
	"github.com/bitly/go-simplejson"
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
