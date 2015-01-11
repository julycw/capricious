package capricious

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/macaron"
	"github.com/julycw/capricious/db"
	"log"
	"reflect"
	"strings"
)

type Options struct {
	URLPrefix       string
	ObjectNameField string
	ObjectIDField   string
}

type capricious struct {
	URLPrefix       string
	ObjectNameField string
	ObjectIDField   string

	ObjectName string
	ObjectID   string

	ContentType string
	AccessToken string
	APPName     string

	dbConn *db.MongoDBConn
}

var MaxMemory = int64(1024 * 1024 * 10)

func prepareOptions(options []Options) Options {
	var op Options
	if len(options) > 0 {
		op = options[0]
	}

	if len(op.URLPrefix) == 0 {
		op.URLPrefix = "/objects/"
	} else if op.URLPrefix[len(op.URLPrefix)-1] != '/' {
		op.URLPrefix += "/"
	}

	if len(op.ObjectIDField) == 0 {
		op.ObjectIDField = "id"
	} else {
		op.ObjectIDField = strings.TrimLeft(op.ObjectIDField, ":")
	}

	if len(op.ObjectNameField) == 0 {
		op.ObjectNameField = "object"
	} else {
		op.ObjectNameField = strings.TrimLeft(op.ObjectNameField, ":")
	}

	return op
}

func newCapricious(options Options) *capricious {
	conn, err := db.GetConn()
	if err != nil {
		panic("init connection failed!")
	}
	return &capricious{
		URLPrefix:       options.URLPrefix,
		ObjectNameField: options.ObjectNameField,
		ObjectIDField:   options.ObjectIDField,
		dbConn:          conn,
	}
}

func (this *capricious) Get(ctx *macaron.Context) {
	if this.ObjectID == "" {
		results, _, err := this.getContext().GetAll()
		if err != nil {
			responseError(ctx, err)
		} else {
			responseSuccess(ctx, results)
		}
	} else {
		result, err := this.getContext().Get(this.ObjectID)
		if err != nil {
			responseError(ctx, err)
		} else {
			responseSuccess(ctx, result)
		}
	}

}

func (this *capricious) Post(ctx *macaron.Context) {
	data := this.parseData(ctx)
	if id, err := this.getContext().Insert(&data); err != nil {
		responseError(ctx, err)
	} else {
		responseSuccess(ctx, id)
	}

}

func (this *capricious) Put(ctx *macaron.Context) {
	data := this.parseData(ctx)
	err := this.getContext().Update(this.ObjectID, &data)
	if err != nil {
		responseError(ctx, err)
	} else {
		responseSuccess(ctx, this.ObjectID)
	}
}

func (this *capricious) Delete(ctx *macaron.Context) {
	err := this.getContext().Delete(this.ObjectID)
	if err != nil {
		responseError(ctx, err)
	} else {
		responseSuccess(ctx, this.ObjectID)
	}
}
func (this *capricious) Prepare(ctx *macaron.Context) bool {
	// 如果是post或put  需要先进行处理
	if ctx.Req.Method == "POST" || ctx.Req.Method == "PUT" {
		if ctx.Req.MultipartForm == nil {
			// Workaround for multipart forms returning nil instead of an error
			// when content is not multipart; see https://code.google.com/p/go/issues/detail?id=6334
			if multipartReader, err := ctx.Req.MultipartReader(); err != nil {
				// errors.Add([]string{}, ERR_DESERIALIZATION, err.Error())
			} else {
				form, parseErr := multipartReader.ReadForm(MaxMemory)
				if parseErr != nil {
					// errors.Add([]string{}, ERR_DESERIALIZATION, parseErr.Error())
				}
				ctx.Req.MultipartForm = form
			}
		}
	}

	if this.APPName = ctx.Req.Header.Get("app_name"); this.APPName == "" {
		responseError(ctx, fmt.Errorf("未能从http head中读取app_name"))
		return false
	}
	if this.AccessToken = ctx.Req.Header.Get("access_token"); this.AccessToken == "" {
		responseError(ctx, fmt.Errorf("未能从http head中读取access_token"))
		return false
	}

	this.ContentType = ctx.Req.Header.Get("Content-Type")
	this.ObjectName = strings.TrimSpace(ctx.Params(":" + this.ObjectNameField))
	this.ObjectID = strings.TrimSpace(ctx.Params(":" + this.ObjectIDField))
	return true
}

func (this *capricious) getContext() *db.MongoDBContext {
	return this.dbConn.GetContext(this.APPName, this.ObjectName)
}

func (this *capricious) parseData(ctx *macaron.Context) db.DataStruct {
	data := db.NewDataStruct(nil)
	for key, value := range ctx.Req.MultipartForm.Value {
		if len(value) > 0 {
			data[key] = value[len(value)-1]
		}
	}
	return data
}

func responseSuccess(ctx *macaron.Context, result interface{}) {
	t := reflect.TypeOf(result)

	switch t.Kind() {
	case reflect.String:
		ctx.Resp.Header().Set("Content-Type", "text/plain")
		ctx.Resp.Write([]byte(result.(string)))
	default:
		ctx.Resp.Header().Set("Content-Type", "application/json")
		bytes, _ := json.Marshal(&result)
		ctx.Resp.Write(bytes)
	}
	ctx.Resp.WriteHeader(200)
}

func responseError(ctx *macaron.Context, err error) {
	ctx.Resp.WriteHeader(403)
	ctx.Resp.Write([]byte(err.Error()))
}

func Capricious(options ...Options) macaron.Handler {
	cap := newCapricious(prepareOptions(options))
	return func(ctx *macaron.Context, logger *log.Logger) {
		//过滤其他URL
		if strings.HasPrefix(ctx.Req.RequestURI, cap.URLPrefix) {
			if !cap.Prepare(ctx) {
				ctx.Invoke(responseError)
				return
			}
			switch ctx.Req.Method {
			case "GET":
				cap.Get(ctx)
			case "POST":
				cap.Post(ctx)
			case "PUT":
				cap.Put(ctx)
			case "DELETE":
				cap.Delete(ctx)
			}
		}
	}
}
