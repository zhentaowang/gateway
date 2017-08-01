package util

import (
	"github.com/Shopify/sarama"
	"encoding/json"
	"strings"
	"log"
)

type InfoCount struct {
	RequestUrl string	`json:"RequestUrl"`
	RequestContent string	`json:"RequestContent"`
	ResponseContent string	`json:"ResponseContent"`
	UsedTime int64	`json:"UsedTime"`
}

func SendToKafka(info *InfoCount)  {

	defer ErrHandle()
	jsons, errs := json.Marshal(info)
	if errs != nil {
		panic(errs)
	}

	//logger := log.New(os.Stderr, "[srama]", log.LstdFlags|log.Llongfile)

	conf := GetConfigCenterInstance()

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

			log.Println("Failed to produce message: ", err)

		}

		log.Printf("partition=%d, offset=%d\n", partition, offset)
	} else {
		log.Println(string(jsons))
	}

}
