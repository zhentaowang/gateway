package handler


import (

	"github.com/go-xorm/xorm"
	"gateway/src/util"
	"strconv"
	"log"
)



var Engine *xorm.Engine



/**
映射数据库api表
 */
type Api struct {
	ApiId       int	`xorm:"pk"`
	Name        string
	Uri         string
	Method      string
	NeedLogin   int
	DisplayName string
	Status      int
	ServiceId   int
	ServiceProviderName	string
	Mock        string
	Desc        string
	Filters     string
}

/**
映射数据库service表
 */
type Service struct {

	ServiceId int	`xorm:"pk"`
	Namespace string
	Name      string
	Desc      string
	Port      string
	Protocol  string
}

/**
映射数据库filter表
 */
type Filter struct {

	FilterId int	`xorm:"pk"`
	ApiId    int
	Name     string
	Seq      int
}

type LoginData struct {
	name	string
	password	string
}

/**
初始化数据库连接
 */
func init()  {

	conf := util.GetConfigCenterInstance()

	var MysqlUrl string = conf.ConfProperties["jdbc"]["db_username"] + ":" + conf.ConfProperties["jdbc"]["db_password"] + "@tcp(" + conf.ConfProperties["jdbc"]["db_host"] + ")/" +
		conf.ConfProperties["jdbc"]["db_name"] + "?charset=utf8"

	Engine, _ = xorm.NewEngine("mysql", MysqlUrl)

}


/**
检查错误，如果错误不为空，打印错误
 */
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}



func MInsertApi(data *Api, filter_seq int) {

	api := new(Api)
	total, err := Engine.Where("api_id =?", data.ApiId).Count(api)
	CheckErr(err)

	if total >0 {
		MDeleteApi(data)
	}
	affected, err := Engine.Insert(data)
	CheckErr(err)
	log.Println("  insert api " +data.Uri+ " success "+strconv.FormatInt(affected,10))

	util.SetData()
}

func MutiInsertApi(data []*Api)  {

	for _,d := range data {
		service := new(Service)
		api := new(Api)
		_ ,err := Engine.Get(service)
		d.ServiceId = service.ServiceId

		total, err := Engine.Where("uri =?", d.Uri).Count(api)

		CheckErr(err)

		if total == 0 {
			affected, err := Engine.Insert(d)

			CheckErr(err)
			log.Println("  insert api " +d.Uri+ " success "+strconv.FormatInt(affected,10))
		}
	}

	util.SetData()
}


func MInsertService(data *Service)  {

	service := new(Service)

	total, err := Engine.Where("service_id =?", data.ServiceId).Count(service)
	CheckErr(err)

	if total >0 {
		MDeleteService(data)
	}
	
	affected, err := Engine.Insert(data)
	CheckErr(err)
	log.Println("  insert service " +data.Namespace+"."+ data.Name +":"+data.Port+ " success "+strconv.FormatInt(affected,10))

	util.SetData()
}


func MInsertFilter(data *Filter)  {

	filter := new(Filter)
	total, err := Engine.Where("filter_id =?", data.FilterId).Count(filter)
	CheckErr(err)

	if total >0 {
		MDeleteFilter(data)
	}

	affected, err := Engine.Insert(data)
	CheckErr(err)

	log.Println("  insert filter " + data.Name + " success " + strconv.FormatInt(affected,10))
	util.SetData()
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
	util.SetData()

}

func MModifyService(data *Service)  {

	affected, err := Engine.Id(data.ServiceId).Update(data)
	CheckErr(err)
	println(affected)
	util.SetData()

}

func MModifyFilter(data *Filter) {

	affected, err := Engine.Id(data.FilterId).Update(data)
	CheckErr(err)
	println(affected)
	util.SetData()

}

func MDeleteApi(data *Api) {

	api := new(Api)
	affected, err := Engine.Id(data.ApiId).Delete(api)

	CheckErr(err)
	log.Println("  have deleted api " + data.Uri +"  id="+strconv.Itoa(data.ApiId)+ "  success " + strconv.FormatInt(affected,10))

	util.SetData()
}

func MDeleteService(data *Service) {

	service := new(Service)
	affected, err := Engine.Id(data.ServiceId).Delete(service)
	CheckErr(err)
	log.Println("  have deleted service " +data.Namespace+"."+ data.Name +":"+data.Port+"  id="+strconv.Itoa(data.ServiceId)+ " success " + strconv.FormatInt(affected,10))

	var ApiData = Api{0, "", "", "", 0, "", 0,data.ServiceId,"" ,"", "", ""}
	affected, err = Engine.Delete(ApiData)
	CheckErr(err)
	log.Println("  have deleted api where ServiceId =  " + strconv.Itoa(data.ServiceId)+ " success " + strconv.FormatInt(affected,10))

	util.SetData()
}

func MDeleteFilter(data *Filter)  {

	filter := new(Filter)
	affected, err := Engine.Id(data.FilterId).Delete(filter)
	CheckErr(err)
	log.Println("  haved deleted filter " + data.Name+"  id="+strconv.Itoa(data.FilterId) + "  success " + strconv.FormatInt(affected,10))

	util.SetData()

}