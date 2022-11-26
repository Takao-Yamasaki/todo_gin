package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// HTMLを読み込むディレクトリを指定
	router.LoadHTMLGlob("templates/*.html")

	data := "Hello Go/Gin!!"

	// index.htmlにGetで繋いでいる
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{"data": data})
	})

	router.Run()
}
