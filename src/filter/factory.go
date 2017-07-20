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
    FilterUpdateFlight = "UPDATE_FLIGHT"
    FilterResponseHead = "RESPONSE_HEAD" // 原响应头
    FilterVisitCount = "VISITCOUNT"
    FilterText = "FILTERTEXT"
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
    case FilterResponseHead:
        return newResponseHeaderFilter(), nil
    case FilterVisitCount:
        return newVisitCount(), nil
    case FilterText:
        return newTextFilter(), nil

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

        log.Printf("Filter <%s> added",f.Name())
        filters.PushBack(f)
    }

    return filters
}
