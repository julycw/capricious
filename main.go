package main

import (
	"fmt"
	"github.com/Unknwon/macaron"
	"github.com/julycw/capricious/capricious"
)

const (
	objectNameField string = "object"
	objectIDField          = "id"
)

func main() {
	m := macaron.Classic()

	m.Use(capricious.Capricious(capricious.Options{
		URLPrefix:       "/objects/",
		ObjectNameField: objectNameField,
		ObjectIDField:   objectIDField,
	}))

	//匹配   /objects/book [get]
	m.Get(fmt.Sprintf("/objects/:%s:string", objectNameField))
	//匹配   /objects/book/1 [get]
	m.Get(fmt.Sprintf("/objects/:%s:string/:%s", objectNameField, objectIDField))
	//匹配   /objects/book [post]
	m.Post(fmt.Sprintf("/objects/:%s:string", objectNameField))
	//匹配   /objects/book/1 [put]
	m.Put(fmt.Sprintf("/objects/:%s:string/:%s", objectNameField, objectIDField))
	//匹配   /objects/book/1 [delete]
	m.Delete(fmt.Sprintf("/objects/:%s:string/:%s", objectNameField, objectIDField))

	m.Run()
}
