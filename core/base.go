package core

import (
	"goserver/utils"
	"github.com/gin-gonic/gin"
	"goserver/gpa"
	"reflect"
	"goserver/webs"
	"strings"
	"os"
)

func putFunRun(fun func()) {
	if len(initAfterFun) > 0 {
		go fun()
	} else {
		initAfterFun = append(initAfterFun, fun)
	}
}

func put(ele *utils.Element, v interface{}) {
	id, idb := ele.AttrValue("Id")
	if !idb {
		id = ele.Name()
	}
	_, de := data[id]
	if de {
		panic("Id重复" + id)
	} else {
		data[id] = v
	}
}

func getVerify(ele *utils.Element) webs.BaseFun {
	VerifyId, vb := ele.AttrValue("VerifyRef")
	if !vb {
		VerifyId = "Verify"
	}
	return data[VerifyId].(webs.BaseFun)
}

func getGpa(ele *utils.Element) *gpa.Gpa {
	ref := ele.Attr("GpaRef", "Gpa")
	web := data[ref].(*gpa.Gpa)
	return web
}

func getGin(ele *utils.Element) *gin.Engine {
	ref := ele.Attr("WebRef", "Web")
	web := data[ref].(*gin.Engine)
	return web
}

func doSubElement(ele *utils.Element, obj interface{}) {
	ns := ele.AllNodes()
	if len(ns) > 0 {
		spv := reflect.ValueOf(obj)
		for _, e := range ns {
			inputs := make([]reflect.Value, 2)
			inputs[0] = reflect.ValueOf(e)
			inputs[1] = reflect.ValueOf(data)
			m := spv.MethodByName(e.Name())
			m.Call(inputs)
		}
	}
}

func post(ele *utils.Element, fun func(param *webs.Param)) {
	getGin(ele).POST(ele.MustAttr("Url"), func(c *gin.Context) {
		wb := webs.NewParam(c)
		fun(wb)
		c.JSON(200, wb.Out)
	})
}

func getFile(file string) string {
	dir := os.Args[1]
	fg := []string{"/", "\\"}
	for _, flg := range fg {
		lst := strings.LastIndex(dir, flg)
		if lst > 0 {
			return dir[0:lst+1] + file
		}
	}
	return file
}
