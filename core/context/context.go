package context

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
)

const (
	LogIdKey    = "LOG_ID"
	QuestionKey = "QUESTION"
	QTypeKey    = "Q_TYPE"
	SrcKey      = "SRC"
	LocationKey = "LOCATION"
)

type Context struct {
	context.Context
	fields map[string]interface{}
}

func (ctx *Context) Fields() logrus.Fields {
	_, file, no, ok := runtime.Caller(1)
	delete(ctx.fields, LocationKey)
	if ok {
		file = filepath.Base(file)
		ctx.fields[LocationKey] = fmt.Sprintf("%s:%d", file, no)
	}
	return ctx.fields
}

func NewContext(writer dns.ResponseWriter, request *dns.Msg) *Context {
	question := request.Question[0]
	src := writer.RemoteAddr().String()
	ctx := &Context{
		Context: context.Background(),
		fields: map[string]interface{}{
			LogIdKey:    strconv.FormatInt(int64(request.Id), 10),
			QuestionKey: question.Name,
			QTypeKey:    dns.Type(question.Qtype).String(),
			SrcKey:      src,
			//SrcKey:    src[:strings.LastIndex(src, ":")],
		},
	}
	return ctx
}

func NewEmptyContext(logId uint16) *Context {
	ctx := &Context{
		Context: context.Background(),
		fields: map[string]interface{}{
			LogIdKey: strconv.FormatInt(int64(logId), 10),
		},
	}
	return ctx
}
