package conf_center

import "encoding/json"

type DataChangeHandler interface {
	OnConfCreated(createdConf string, newProperties map[string]string)
	OnConfDeleted(createdConf string, oldProperties map[string]string)
	OnConfUpdated(createdConf string, oldProperties map[string]string, newProperties map[string]string)
}

type LogDataChangeHandler struct {

}

func (logDataChangeHandler LogDataChangeHandler)OnConfCreated(createdConf string, newProperties map[string]string) {
	println("Created " + createdConf + ", data :" + jsonStr(newProperties))
}

func (logDataChangeHandler LogDataChangeHandler)OnConfDeleted(createdConf string, oldProperties map[string]string) {
	println("Delete " + createdConf + ", data :" + jsonStr(oldProperties))
}

func (logDataChangeHandler LogDataChangeHandler)OnConfUpdated(createdConf string, oldProperties map[string]string, newProperties map[string]string){
	println("Updated " + createdConf + ", data : from" + jsonStr(oldProperties) + ", to:" + jsonStr(newProperties))
}

func jsonStr(aMap map[string]string) string {
	bytes, _ := json.Marshal(aMap)
	return string(bytes)
}
