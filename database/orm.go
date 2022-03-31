package database

var (
	DefaultOrm = NewOrm()
)

type Orm interface {
	Init(...Option) error
	Connect() error
	Get() interface{}
	String() string
}
