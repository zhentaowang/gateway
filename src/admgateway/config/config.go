package config


type MysqlConfig struct {

	MysqlHost string	`yaml:"MysqlHost"`
	MysqlPort string	`yaml:"MysqlPort"`
	MysqlUserName string	`yaml:"MysqlUserName"`
	MysqlPassword string	`yaml:"MysqlPassword"`
	MysqlDb string	`yaml:"MysqlDb"`
}
