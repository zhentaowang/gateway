package filter

import (
	"gateway/src/util"
)

type VisitCount struct {
	BaseFilter
}

func newVisitCount()  Filter{
	return &VisitCount{}
}

func (vc *VisitCount) Name()  string{
	return FilterVisitCount
}

func (inf *VisitCount) Post(c Context)  (int, error){


	handleInfo := new(util.InfoCount)
	handleInfo.RequestUrl = string(c.GetOriginRequestCtx().Request.RequestURI())
	handleInfo.UsedTime = c.GetEndAt() - c.GetStartAt()
	handleInfo.ResponseContent = c.GetProxyResponse().String()
	handleInfo.RequestContent = c.GetOriginRequestCtx().Request.String()

	util.SendToKafka(handleInfo,"kafka_topic")

	return c.GetProxyResponse().StatusCode(), nil
}
