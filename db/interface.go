package db

import ()

type ContextInterface interface {
	IsExist() bool

	//输入：id:主键，data:数据
	//输出：错误信息
	Update(data *DataStruct) error

	//输入：data:数据
	//输出：id、错误信息
	Insert(data *DataStruct) (string, error)

	//输入：id:主键
	//输出：数据
	Get(id string) (DataStruct, error)

	//输出：数据列表、数量
	GetAll() ([]DataStruct, int, error)

	//输入：id:主键
	//输出：错误信息
	Delete(id string) error

	//输出：数量
	Count() int
}

type Context struct {
	ContextInterface
	AppName     string
	ContextName string
}
