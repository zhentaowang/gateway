package handdler


import (

	"github.com/go-xorm/xorm"
	"admgateway/util"
)



var Engine *xorm.Engine





type Api struct {
	ApiId       int	`xorm:"pk"`
	Name        string
	Uri         string
	Method      string
	NeedLogin   int
	DisplayName string
	Status      int
	ServiceId   int
	Mock        string
	Desc        string
	Filters     string
}

type Service struct {

	ServiceId int	`xorm:"pk"`
	Namespace string
	Name      string
	Desc      string
	Port      string
	Protocol  string
}


type Filter struct {

	FilterId int	`xorm:"pk"`
	ApiId    int
	Name     string
	Seq      int
}


func init()  {
	var err error
	var MysqlUrl string = util.MysqlUrl
	Engine, err = xorm.NewEngine("mysql", MysqlUrl)
	CheckErr(err)
}



func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}



func MInsertApi(data *Api, filter_seq int) {

	
	affected, err := Engine.Insert(data)

	result := MQueryApi(data)
	var filter_data = Filter{0, result[len(result)-1].ApiId, data.Filters, filter_seq}
	MInsertFilter(&filter_data)

	CheckErr(err)
	println(affected )
}


func MInsertService(data *Service)  {
	
	affected, err := Engine.Insert(data)

	CheckErr(err)
	println(affected )
}


func MInsertFilter(data *Filter)  {
	
	affected, err := Engine.Insert(data)
	var ApiData = Api{0, "", "", "", 0, "", 0, 0, "", "", data.Name}
	affected, err = Engine.Where("api_id = ?",data.ApiId).Cols("filters").Update(ApiData)

	CheckErr(err)
	println(affected )
}

func MQueryService(data *Service) []Service {


	AllService := make([]Service, 0)
	err := Engine.Find(&AllService, data)
	CheckErr(err)

	return AllService

}



func MQueryApi(data *Api) []Api {


	AllApi := make([]Api, 0)
	err := Engine.Find(&AllApi, data)

	CheckErr(err)

	return AllApi
}



func MQueryFilter(data *Filter) []Filter {


	AllFilter := make([]Filter, 0, 10)
	err := Engine.Find(&AllFilter, data)

	CheckErr(err)

	return AllFilter
}



func MModifyApi(data *Api) {

	affected, err := Engine.Id(data.ApiId).Update(data)
	CheckErr(err)
	println(affected)

}

func MModifyService(data *Service)  {

	affected, err := Engine.Id(data.ServiceId).Update(data)
	CheckErr(err)
	println(affected)

}

func MModifyFilter(data *Filter) {

	affected, err := Engine.Id(data.FilterId).Update(data)
	CheckErr(err)
	println(affected)

}

func MDeleteApi(data *Api) {

	affected, err := Engine.Delete(data)

	CheckErr(err)
	println(affected)

	var FilterData = Filter{0, data.ApiId, "", 0}
	affected, err = Engine.Delete(FilterData)

	CheckErr(err)
	println(affected)
}

func MDeleteService(data *Service) {

	affected, err := Engine.Delete(data)
	var ApiData = Api{0, "", "", "", 0, "", data.ServiceId, 0, "", "", ""}
	affected, err = Engine.Delete(ApiData)

	CheckErr(err)
	println(affected)
}

func MDeleteFilter(data *Filter)  {

	has, err := Engine.Where("filter_id=?", data.FilterId).Get(data)
	affected, err := Engine.Delete(data)

	var ApiData = Api{0, "", "", "", 0, "", 0, 0, "", "", ""}
	affected, err = Engine.Where("filters = ?", data.Name).Cols("filters").Update(ApiData)

	CheckErr(err)
	println(affected, has)

}