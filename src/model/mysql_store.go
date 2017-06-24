package model

import (
    "database/sql"
    "log"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "encoding/json"
)

type MysqlStore struct {
    DB *sql.DB
}

func NewMysqlStore(host string, username string, password string, dbName string) *MysqlStore {
    m := &MysqlStore{}
    m.dbInit(host, username, password, dbName)
    return m
}

func (m *MysqlStore) dbInit(host string, username string, password string, dbName string) {
    var err error
    m.DB, err = sql.Open("mysql", username + ":" + password + "@tcp(" + host + ")/" + dbName) //返回一个连接池，不是单个连接
    if err != nil {
        panic(err)
    }
    m.DB.SetMaxOpenConns(100) //最大连接数
    m.DB.SetMaxIdleConns(50)  //最大闲置数
    m.DB.Ping()
}

func (m *MysqlStore) GetAPIs() ([]*API, error) {
    rows, err := m.DB.Query("select api_id, name, uri, method, service_id, status, need_login, mock from api")
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    defer rows.Close()
    var value []*API
    for rows.Next() {
        api := new(API)
        var mockStr []byte
        rows.Scan(&api.APIId, &api.Name, &api.URI, &api.Method, &api.ServiceId, &api.Status, &api.NeedLogin, &mockStr)
        mock := new(Mock)
        err := json.Unmarshal(mockStr, mock)
        if err != nil {
            log.Fatal(err)
            return nil, err
        }

        api.filterNames, _ = m.GetFilters(api.APIId)
        api.Mock = mock
        value = append(value, api)
    }

    return value, nil
}

func (m *MysqlStore) GetServices() ([]*Service, error) {
    rows, err := m.DB.Query("select service_id, namespace, name, port, protocol from service")
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    defer rows.Close()
    var value []*Service
    for rows.Next() {
        service := new(Service)
        rows.Scan(&service.ServiceId, &service.Namespace, &service.Name, &service.Port, &service.Protocol)
        value = append(value, service)
    }

    return value, nil
}

/*
apiId 为-1的时候表示系统的filter
 */
func (m *MysqlStore) GetFilters(apiId int) ([]string, error) {
    var query = fmt.Sprintf("select name from filter where api_id=%d order by seq", apiId)
    rows, err := m.DB.Query(query)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    defer rows.Close()
    var value []string
    for rows.Next() {
        var filter string
        rows.Scan(&filter)
        value = append(value, filter)
    }

    return value, nil
}