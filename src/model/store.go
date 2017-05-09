package model

type Store interface {
    GetAPIs() ([]*API, error)
    GetServices() ([]*Service, error)
    GetFilters(apiId int) ([]string, error)
}
