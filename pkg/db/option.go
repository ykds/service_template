package db

type DriverType string

const (
	Mysql DriverType = "mysql"
)

type Option struct {
	Driver   DriverType `json:"driver" yaml:"driver"`
	Host     string     `json:"host" yaml:"host"`
	Port     string     `json:"port" yaml:"port"`
	User     string     `json:"user" yaml:"user"`
	Password string     `json:"password" yaml:"password"`
	DbName   string     `json:"db_name" yaml:"db_name"`
}

func DefaultOption() Option {
	return Option{
		Driver:   Mysql,
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "",
		DbName:   "",
	}
}
