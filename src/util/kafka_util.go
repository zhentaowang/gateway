package util

import (
	"github.com/Shopify/sarama"
	"encoding/json"
	"strings"
	"log"
	"fmt"
	"runtime/debug"
)

type InfoCount struct {
	RequestUrl string	`json:"RequestUrl"`
	RequestContent string	`json:"RequestContent"`
	ResponseContent string	`json:"ResponseContent"`
	UsedTime int64	`json:"UsedTime"`
}


func SendToKafka(info *InfoCount , topic string)  {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR!! ",err)
			v := fmt.Sprintf("ERROR!!\n%s--\n  stack \n%s", err,string(debug.Stack()))
			SendToDingDing(v)
		}
	}()

	jsons, errs := json.Marshal(info)
	if errs != nil {
		log.Panic(errs)
	}

	//logger := log.New(os.Stderr, "[srama]", log.LstdFlags|log.Llongfile)

	conf := GetConfigCenterInstance()

	if conf.ConfProperties["kafka"]["kafka_host"] != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		producer, err := sarama.NewSyncProducer(strings.Split(conf.ConfProperties["kafka"]["kafka_host"], ","), config)
		if err != nil {
			log.Panic(err)
		}

		defer producer.Close()

		msg := &sarama.ProducerMessage{}

		msg.Topic = conf.ConfProperties["kafka"][topic]
		msg.Partition = int32(-1)
		msg.Key = sarama.StringEncoder("info")
		msg.Value = sarama.ByteEncoder(string(jsons))

		partition, offset, err := producer.SendMessage(msg)

		if err != nil {

			log.Println("Failed to produce message: ", err)

		}

		log.Printf("partition=%d, offset=%d\n", partition, offset)
	} else {
		log.Println(string(jsons))
	}

}
