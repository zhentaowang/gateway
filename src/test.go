//package src
//
//import "io"
//
//func main()  {
//	out :=new(io.Writer)
//}
package main

import (
	"github.com/goinggo/mapstructure"
)

func main() {

	//conf := conf_center.New("gateway")
	//conf.Init()
	//jsonStr, _ := json.Marshal(conf.ConfProperties)
	//println(string(jsonStr))
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON where we're not quite sure of the
	// struct initially.
	input := map[string]interface{}{
		"name":   "Mitchell",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
		"sda": "tret",
		"yuoi": "dfsdf",
	}

	var result Person
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%#v", result)
	println( input["name"].(string))
}