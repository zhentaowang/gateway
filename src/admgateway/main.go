package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gateway/src/admgateway/handler"
)




func main()  {
	handler.Run()
}