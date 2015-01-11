package db

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_MongoDBContext(t *testing.T) {
	Convey("Init Connection", t, func() {
		conn, err := GetConn()

		So(conn, ShouldNotEqual, nil)
		So(err, ShouldEqual, nil)

		context := conn.GetContext("test_app", "test_context")
		So(context, ShouldNotEqual, nil)

		count_before_insert := context.Count()
		// log.Println("Count before insert:", count_before_insert)

		var uuid string

		Convey("Test C/U/R/D", func() {
			data := NewDataStruct(map[string]interface{}{
				"name": "test",
				"age":  19,
			})
			uuid, _ = context.Insert(&data)

			//After insert, context count should GREATER THAN before
			So(count_before_insert+1, ShouldEqual, context.Count())
			So(len(uuid), ShouldEqual, 36)
			//Test Get
			data, _ = context.Get(uuid)
			So(data["name"], ShouldEqual, "test")
			So(data["age"], ShouldEqual, 19)

			//Test Update
			updated_data := NewDataStruct(map[string]interface{}{
				"name": "superman",
				"age":  25,
			})
			context.Update(uuid, &updated_data)
			So(count_before_insert+1, ShouldEqual, context.Count())
			data, _ = context.Get(uuid)
			So(data["name"], ShouldEqual, "superman")
			So(data["age"], ShouldEqual, 25)

			//Test GetAll
			_, count, _ := context.GetAll()
			So(count, ShouldEqual, count_before_insert+1)

			//Test Delete
			context.Delete(uuid)
			So(count_before_insert, ShouldEqual, context.Count())
		})
	})
}
