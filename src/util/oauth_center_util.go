package util

import (
	"net/http"
	"strings"
	"io/ioutil"
	"log"
	"github.com/bitly/go-simplejson"
)

func GetToken(username string , password string)  string{
	post := "{\"username\":"+username+"," +
		"\"password\":"+password+"," +
		"\"grant_type\":\"password\"," +
		"\"client_id\":\"gateway\" ," +
		"\"client_secret\":\"A9B1-4D6F3F3044B1\"}"

	resp, err := http.Post("oauth-center.platform:443/oauth/access_token",
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
	resp, err := http.Get("oauth-center.platform:443/user/permission?access_token="+access_token+"&permission="+permission)
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