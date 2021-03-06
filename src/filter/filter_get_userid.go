package filter

import (
    "net/http"
    "gateway/src/util"
    "errors"
)

type VisitOauth struct {
    BaseFilter
}

func newVisitOauth() Filter {
    return &VisitOauth{}
}

func (f VisitOauth) Name() string {
    return FilterVisitOauth
}

func (v VisitOauth) Pre(c Context) (statusCode int, err error) {

    if access_token := c.GetProxyOuterRequest().URI().QueryArgs().Peek("access_token");access_token != nil {
        ok , err := util.CheckToken(c.GetProxyOuterRequest(),c.GetProxyResponse())
        if ok {
            return http.StatusOK,nil
        }

        if err == nil {
            err = errors.New("token认证失败")
        }

        return http.StatusForbidden,err
    }

    return http.StatusOK,nil

}