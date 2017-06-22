package util

import (
	"github.com/Shopify/sarama"
	"strings"
	"os"
	"log"
	"encoding/json"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"gateway/src/config"
)

type HttpHandlInfo struct {
	RequestUrl string	`json:"RequestUrl"`
	RequestContent string	`json:"RequestContent"`
	ResponseContent string	`json:"ResponseContent"`
	UsedTime int64	`json:"UsedTime"`
}


func (inf *HttpHandlInfo) VisitCount()  {

	logger := log.New(os.Stderr, "[srama]", log.LstdFlags)

	configByte, err := ioutil.ReadFile("config.yml")
	conf := new(config.KafkaConfig)
	err = yaml.Unmarshal(configByte, &conf)
	if err != nil {
		panic(err)
	}


	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(strings.Split(conf.KafkaHost, ","), config)
	if err != nil {
		panic(err)
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{}
	jsons, errs := json.Marshal(inf)
	if errs != nil {
		panic(errs)
	}

	msg.Topic = conf.KafkaTopic
	msg.Partition = int32(-1)
	msg.Key = sarama.StringEncoder("info")
	msg.Value = sarama.ByteEncoder(string(jsons))

	partition, offset, err := producer.SendMessage(msg)


	if err != nil {

		logger.Println("Failed to produce message: ", err)

	}

	logger.Printf("partition=%d, offset=%d\n", partition, offset)
}

func (inf *HttpHandlInfo) SetRequestUrl(url string)  {
	inf.RequestUrl = url
}

func (inf *HttpHandlInfo) SetRequestContent(content string)  {
	inf.RequestContent = content
}

func (inf *HttpHandlInfo) SetResponseContent(content string)  {
	inf.ResponseContent = content
}

func (inf *HttpHandlInfo) SetUsedTime(time int64)  {
	inf.UsedTime = time
}
