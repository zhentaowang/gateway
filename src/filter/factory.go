package filter

import (
    "strings"
    "errors"
    "container/list"
    "log"
)

var (
    // ErrUnknownFilter unknown filter error
    ErrUnknownFilter = errors.New("unknow filter")
)

const (
    // FilterRights 验证用户是否有权限访问资源
    FilterRights = "RIGHTS"
    FilterCORS = "CORS"
    // 给龙腾推送航班信息的参数过滤器
    FilterUpdateFlight = "SPLIT_SIGN"
)

func NewFilter(filterName string) (Filter, error) {
    input := strings.ToUpper(filterName)

    switch input {
    case FilterRights:
        return newRightsFilter(), nil
    case FilterCORS:
        return newCORSFilter(), nil
    case FilterUpdateFlight:
        return newUpdateFlightFilter(), nil
    default:
        return nil, ErrUnknownFilter
    }
}

func NewFilters(filterNames []string) (*list.List) {
    var filters = list.New()
    if filterNames == nil || len(filterNames) == 0 {
        return filters
    }

    for _, filterName := range filterNames {
        f, err := NewFilter(filterName)
        if nil != err {
            log.Panicf("Proxy unknow filter <%+v>", filterName)
        }

        log.Println(f)
        filters.PushBack(f)
    }

    return filters
}
