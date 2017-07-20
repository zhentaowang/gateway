package conf_center

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"strings"
	"os"
)

type AppProperties struct {
	AppPath            string
	AppName            string
	conn               *zk.Conn
	zkServers          []string
	secretKey          string
	dataChangeHandlers []DataChangeHandler
	ConfProperties     map[string]map[string]string
}

func New(appName string) AppProperties {
	appProperties := AppProperties{AppName: appName}
	appProperties.AppPath = os.Getenv("zk_root")+ appName
	appProperties.zkServers = []string{os.Getenv("zk_servers")}
	appProperties.secretKey = os.Getenv("secret_key")
	return appProperties
}

func NewWithArgs(zkRootPath string,appName string,zkServer []string,secretKey string) AppProperties {
	appProperties := AppProperties{AppName: appName}
	appProperties.AppPath = zkRootPath + appName
	appProperties.zkServers = zkServer
	appProperties.secretKey = secretKey
	return appProperties
}

func (appProperties *AppProperties)Init() {
	conn, _, _ := zk.Connect(appProperties.zkServers, time.Second * 10)
	appProperties.conn = conn
	appProperties.ConfProperties = appProperties.loadProperties(true)
}

func (appProperties *AppProperties)loadProperties(watchChild bool) map[string]map[string]string {
	confProperties := map[string]map[string]string{}
	children, _, event, _ := appProperties.conn.ChildrenW(appProperties.AppPath)
	go appProperties.watch(event)
	for _, child := range children {
		childData := appProperties.getData(appProperties.AppPath + "/" + child, watchChild)
		confProperties[child] = childData
	}
	return confProperties
}

func (appProperties *AppProperties)getData(path string, watch bool) map[string]string {
	bytes, _, event, _ := appProperties.conn.GetW(path)
	if watch {
		go appProperties.watch(event)
	}
	return appProperties.extractData(string(bytes))
}

func (appProperties *AppProperties)extractData(data string) map[string]string {
	Try(func() {
		plaintext, err := AesDecrypt(data, appProperties.secretKey)
		if err != nil {
			println("Decrypt error.")
		} else {
			data = plaintext
		}
	}, func(e interface{}) {
		println("Not enrypt data.")
	})
	properties := map[string]string{}
	splits := strings.Split(data, "\n")
	for _, split := range splits {
		index := strings.Index(split, "=")
		name := split[0:index]
		value := split[index + 1:]
		properties[name] = value
	}
	return properties
}

func (appProperties *AppProperties) RegisterDataChangeHandler(dataChangeHandler DataChangeHandler) {
	appProperties.dataChangeHandlers = append(appProperties.dataChangeHandlers, dataChangeHandler)
}

func (appProperties *AppProperties)watch(event <- chan zk.Event) {
	e := <-event
	eventType := e.Type
	path := e.Path
	if e.Err != nil {
		return
	}
	if eventType == zk.EventNodeChildrenChanged {
		newConfProperties := appProperties.loadProperties(false)
		createdProperties := MapMinus(newConfProperties, appProperties.ConfProperties)
		deletedProperties := MapMinus(appProperties.ConfProperties, newConfProperties)
		if len(createdProperties) > 0 {
			appProperties.doConfCreated(createdProperties, newConfProperties)
		}
		if len(deletedProperties) > 0 {
			appProperties.doConfDeleted(deletedProperties, newConfProperties)
		}
	}
	if eventType == zk.EventNodeDataChanged {
		newData := appProperties.getData(path, true)
		index := strings.LastIndex(path, "/")
		confName := path[index + 1 :]
		appProperties.doConfUpdated(confName, appProperties.ConfProperties[confName], newData)
	}
}

func (appProperties *AppProperties) doConfCreated(createdConfs []string, newConfProperties map[string]map[string]string) {
	// watch new node
	for _, createdConf := range createdConfs {
		_, _, event, _ := appProperties.conn.ExistsW(appProperties.AppPath + "/" + createdConf)
		go appProperties.watch(event)
		for _, dataChangeHandler := range appProperties.dataChangeHandlers {
			dataChangeHandler.OnConfCreated(createdConf, newConfProperties[createdConf])
		}
	}
	appProperties.ConfProperties = newConfProperties
}

func (appProperties *AppProperties) doConfDeleted(deletedConfs []string, newConfProperties map[string]map[string]string) {
	for _, deletedConf := range deletedConfs {
		for _, dataChangeHandler := range appProperties.dataChangeHandlers {
			dataChangeHandler.OnConfDeleted(deletedConf, appProperties.ConfProperties[deletedConf])
		}
	}
	appProperties.ConfProperties = newConfProperties
}

func (appProperties *AppProperties) doConfUpdated(updatedConf string, oldProperties map[string]string, newProperties map[string]string) {
	for _, dataChangeHandler := range appProperties.dataChangeHandlers {
		dataChangeHandler.OnConfUpdated(updatedConf, oldProperties, newProperties)
	}
	appProperties.ConfProperties[updatedConf] = newProperties
}
