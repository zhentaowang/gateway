package filter

import (
	"github.com/Shopify/sarama"
	"strings"
	"os"
	"log"
	"encoding/json"

	"conf_center"
)

type VisitCount struct {
	BaseFilter
	RequestUrl string	`json:"RequestUrl"`
	RequestContent string	`json:"RequestContent"`
	ResponseContent string	`json:"ResponseContent"`
	UsedTime int64	`json:"UsedTime"`
}

func newVisitCount()  Filter{
	return &VisitCount{}
}

func (vc *VisitCount) Name()  string{
	return FilterVisitCount
}

func (inf *VisitCount) Post(c Context)  (int, error){


	handleInfo := new(VisitCount)
	handleInfo.RequestUrl = string(c.GetOriginRequestCtx().Request.RequestURI())
	handleInfo.UsedTime = c.GetEndAt() - c.GetEndAt()
	handleInfo.ResponseContent = c.GetProxyResponse().String()
	handleInfo.RequestContent = c.GetOriginRequestCtx().Request.String()
	jsons, errs := json.Marshal(handleInfo)
	if errs != nil {
		panic(errs)
	}

	logger := log.New(os.Stderr, "[srama]", log.LstdFlags)

	conf := conf_center.New("/wyun/gateway")
	conf.Init()

	if conf.ConfProperties["kafka"]["kafka_host"] != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		producer, err := sarama.NewSyncProducer(strings.Split(conf.ConfProperties["kafka"]["kafka_host"], ","), config)
		if err != nil {
			panic(err)
		}

		defer producer.Close()

		msg := &sarama.ProducerMessage{}

		msg.Topic = conf.ConfProperties["kafka"]["kafka_topic"]
		msg.Partition = int32(-1)
		msg.Key = sarama.StringEncoder("info")
		msg.Value = sarama.ByteEncoder(string(jsons))

		partition, offset, err := producer.SendMessage(msg)

		if err != nil {

			logger.Println("Failed to produce message: ", err)

		}

		logger.Printf("partition=%d, offset=%d\n", partition, offset)
	} else {
		logger.Println(string(jsons))
	}

	return c.GetProxyResponse().StatusCode(), errs
}

func (inf *VisitCount) SetRequestUrl(url string)  {
	inf.RequestUrl = url
}

func (inf *VisitCount) SetRequestContent(content string)  {
	inf.RequestContent = content
}

func (inf *VisitCount) SetResponseContent(content string)  {
	inf.ResponseContent = content
}

func (inf *VisitCount) SetUsedTime(time int64)  {
	inf.UsedTime = time
}
