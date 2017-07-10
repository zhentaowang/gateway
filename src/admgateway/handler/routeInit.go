package handler

import (
	"github.com/buaazp/fasthttprouter"
)

func InitRoute(router *fasthttprouter.Router) {


	router.GET("/", indexHandler)
	router.GET("/filter.html", ToService)
	router.GET("/service.html", ToFilter)
	router.GET("/api.html", indexHandler)
	router.GET("/delete", deleteHandler)
	router.POST("/gateway/admin/add_api", AddApi)
	router.POST("/gateway/admin/add_service", AddService)
	router.POST("/gateway/admin/add_filter", AddFilter)
	router.POST("/gateway/admin/query_api", QueryApi)
	router.POST("/gateway/admin/query_service", QueryService)
	router.POST("/gateway/admin/query_filter", QueryFilter)
	router.POST("/gateway/admin/modify_api", ModifyApi)
	router.POST("/gateway/admin/modify_service", ModifyService)
	router.POST("/gateway/admin/modify_filter", ModifyFilter)
	router.POST("/gateway/admin/delete_api", DeleteApi)
	router.POST("/gateway/admin/delete_service", DeleteService)
	router.POST("/gateway/admin/delete_filter", DeleteFilter)

}
