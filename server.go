package main

import (
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"strings"
	"os"
)

func main() {
	app := iris.New(iris.Configuration{})
	app.Adapt(
		iris.DevLogger(),
		httprouter.New(),
	)

	githubSrc := "https://github.com/nousmotards/mail-template/raw/master/"

	app.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		if len(ctx.Path()) <= 1 {
			return
		}
		filename := ctx.Path()[1:]
		absPath, _ := filepath.Abs("./")
		filename = fmt.Sprintf("%s/%s", absPath, filename)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			ctx.WriteString(fmt.Sprintf("File %s not found", filename))
			return
		}

		ctx.SetStatusCode(iris.StatusOK)
		if len(filename) > 5 && filename[len(filename)-5:] == ".html" {
			file, err := ioutil.ReadFile(filename)
			if err != nil {
				return
			}
			ctx.SetContentType("text/html")
			ctx.WriteString(strings.Replace(string(file), githubSrc, "", -1))
		} else {
			ctx.ServeFile(filename, true)
		}
	})

	app.Listen(":4499")
}
