package capricious

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/macaron"
	"log"
	"strings"
)

type ContentType int

const (
	URLENCODED ContentType = iota + 1
	MULTIPART
	JSON
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

	ContentType ContentType
}

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
	return &capricious{
		URLPrefix:       options.URLPrefix,
		ObjectNameField: options.ObjectNameField,
		ObjectIDField:   options.ObjectIDField,
	}
}

func (this *capricious) Get(ctx *macaron.Context) {

}

func (this *capricious) Post(ctx *macaron.Context) {

}

func (this *capricious) Put(ctx *macaron.Context) {

}

func (this *capricious) Delete(ctx *macaron.Context) {

}
func (this *capricious) Prepare(ctx *macaron.Context) bool {
	contentType := ctx.Req.Header.Get("Content-Type")
	if ctx.Req.Method == "POST" || ctx.Req.Method == "PUT" || len(contentType) > 0 {
		switch {
		case strings.Contains(contentType, "form-urlencoded"):
		case strings.Contains(contentType, "multipart/form-data"):
		case strings.Contains(contentType, "json"):
		default:
			var err error
			if contentType == "" {
				err = fmt.Errorf("Empty Content-Type")
			} else {
				err = fmt.Errorf("Unsupported Content-Type: %s", contentType)
			}
			ctx.Map(err)
			return false
		}
	}

	return true
}

func responseError(ctx *macaron.Context, err error) {
	ctx.Resp.WriteHeader(500)
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
