package filter

import (
	"strings"
)

/*
給龙腾的推送航班过滤器
 */
type DragonFilter struct{
	BaseFilter
}

func newDragonFilter() Filter {
	return &DragonFilter{}
}

func (v DragonFilter) Name() string {
	return FilterDragon
}

func (v DragonFilter) Pre(c Context) (statusCode int, err error) {
	tcUpdate := make([]byte,0,2048)
	tcUpdate = append(c.GetOriginRequestCtx().PostBody()) // 获取到的是 byte[]
	str := "{"
	str += string(tcUpdate[:])
	str = strings.Replace(str,"data=","\"data\":",1)
	str = strings.Replace(str,"&sign=",",\"sign\":\"",1)
	str += "\"}"
	tcUpdate = []byte(str)
	c.GetProxyOuterRequest().SetBody(tcUpdate)
	return
}
