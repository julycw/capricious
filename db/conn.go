package db

import (
	"labix.org/v2/mgo"
)

var conn IConnection = nil

//通过单例模式获取数据库连接
func GetConn() (*IConnection, error) {
	if conn == nil {
		//一般来说url从配置文件中读取
		url := "127.0.0.1"
		sess, err := mgo.Dial(url)

		if err != nil {
			sess.SetMode(mgo.Monotonic, true)
			return nil, err
		}

		conn = MongoDBConn{
			services: []string{url},
			session:  sess,
		}
	}

	return &conn, nil
}
