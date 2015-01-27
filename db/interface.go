package db

type IContext interface {
	IsExist() bool

	Update(id string, data *DataStruct) error

	Insert(data *DataStruct) (string, error)

	Get(id string) (DataStruct, error)

	GetAll() ([]DataStruct, int, error)

	Delete(id string) error

	Count() int
}

type IConnection interface {
	GetContext(appName, contextName string) IContext
}

type Connection struct {
	IConnection
}

type Context struct {
	IContext
	AppName     string
	ContextName string
}
