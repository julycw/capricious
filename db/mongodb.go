package db

import (
	"github.com/julycw/capricious/uuid"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

//数据库连接
type MongoDBConn struct {
	services []string
	session  *mgo.Session
}

func (conn MongoDBConn) GetContext(appName, contextName string) IContext {
	ctx := MongoDBContext{
		conn: &conn,
	}
	ctx.AppName = appName
	ctx.ContextName = contextName

	return ctx
}

//数据对象实体
type MongoDBContext struct {
	conn *MongoDBConn
	Context
}

func (ctx MongoDBContext) IsExist() bool {
	return ctx.Count() > 0
}

func (ctx MongoDBContext) Insert(data *DataStruct) (string, error) {
	clearData(data)

	bson_data := make(bson.M, 0)
	for key, value := range *data {
		bson_data[key] = value
	}

	bson_data["_uuid"] = uuid.New()
	bson_data["_insert_at"] = time.Now()

	err := ctx.conn.session.DB(ctx.AppName).C(ctx.ContextName).Insert(&bson_data)
	if err != nil {
		return "", err
	}
	return bson_data["_uuid"].(string), nil
}

func (ctx MongoDBContext) Update(id string, data *DataStruct) error {
	clearData(data)

	if prev_data, err := ctx.Get(id); err != nil {
		return err
	} else {
		for key, value := range *data {
			prev_data[key] = value
		}

		prev_data["_update_at"] = time.Now()
		err := ctx.conn.session.DB(ctx.AppName).C(ctx.ContextName).Update(bson.M{
			"_uuid": id,
		}, prev_data)

		if err != nil {
			return err
		}
	}

	return nil
}

func (ctx MongoDBContext) Get(id string) (DataStruct, error) {
	var result DataStruct
	err := ctx.conn.session.DB(ctx.AppName).C(ctx.ContextName).Find(bson.M{
		"_uuid": id,
	}).One(&result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ctx MongoDBContext) GetAll() ([]DataStruct, int, error) {
	var results []DataStruct = make([]DataStruct, 0)
	err := ctx.conn.session.DB(ctx.AppName).C(ctx.ContextName).Find(nil).All(&results)

	if err != nil {
		return nil, 0, err
	}
	return results, len(results), nil
}

func (ctx MongoDBContext) Delete(id string) error {
	err := ctx.conn.session.DB(ctx.AppName).C(ctx.ContextName).Remove(bson.M{
		"_uuid": id,
	})

	if err != nil {
		return err
	}
	return nil
}

func (ctx MongoDBContext) Count() int {
	count, _ := ctx.conn.session.DB(ctx.AppName).C(ctx.ContextName).Count()
	return count
}
