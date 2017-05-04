package filter

import (
    "strings"
    "errors"
)

var (
    // ErrUnknownFilter unknown filter error
    ErrUnknownFilter = errors.New("unknow filter")
)

const (
    // FilterRights 验证用户是否有权限访问资源
    FilterRights = "RIGHTS"
)

func NewFilter(filterName string) (Filter, error) {
    input := strings.ToUpper(filterName)

    switch input {
    case FilterRights:
        return newRightsFilter(), nil
    default:
        return nil, ErrUnknownFilter
    }
}
