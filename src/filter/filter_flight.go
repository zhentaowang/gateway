package filter

import (
    "strings"
    "fmt"
    "encoding/json"
)

/*
航班更新--系统间调用 需要放行
 */
type UpdateFlightFilter struct{
    BaseFilter
}

func newUpdateFlightFilter() Filter {
    return &UpdateFlightFilter{}
}

func (v UpdateFlightFilter) Name() string {
    return FilterUpdateFlight
}

/*func (v UpdateFlightFilter) Pre(c Context) (statusCode int, err error) {
    bodyStr := string(c.GetOriginRequestCtx().PostBody())
    fmt.Println("传入的参数")
    fmt.Printf(bodyStr)
    if(strings.Contains(bodyStr,"Notify")){
        bodyStr = string([]rune(bodyStr)[7:])
        var buffer bytes.Buffer
        buffer.WriteString("{\"notify\":\"")
        buffer.WriteString(strings.TrimSpace(bodyStr))
        buffer.WriteString("\"}")
        bodyStr = buffer.String()
        c.GetProxyOuterRequest().SetBody([]byte(bodyStr))
    }
    if(strings.Contains(bodyStr,"sign")){
        tcUpdate := make([]byte,0,2048)
        tcUpdate = append(c.GetOriginRequestCtx().PostBody()) // 获取到的是 byte[]
        str := "{"
        str += string(tcUpdate[:])
        str = strings.Replace(str,"data=","\"data\":",1)
        str = strings.Replace(str,"&sign=",",\"sign\":\"",1)
        str += "\"}"
        tcUpdate = []byte(str)
        c.GetProxyOuterRequest().SetBody(tcUpdate)
    }
    return
}*/

func (v UpdateFlightFilter) Pre(c Context) (statusCode int, err error) {
    bodyStr := string(c.GetOriginRequestCtx().PostBody())
    fmt.Println("传入的参数")
    fmt.Printf(bodyStr)
    params := make(map[string]string)
    params["Notify"] = strings.TrimSpace(bodyStr)
    p,_ := json.Marshal(&params)
    c.GetProxyOuterRequest().SetBody(p)
    return
}
