package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	gorm.Model
	Text   string
	Status string
}

// DBの初期化
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("could not database open")
	}
	//マイグレートを実行
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

// DBの作成
func dbCreate(text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("could not database open")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("could not database open")
	}
	var todo Todo
	db.First(&todo, id)
	//構造体Todoを呼んでいる
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("could not database open")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

// DBの全件取得
func dbGetAll() []Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("could not database open")
	}
	// 要素を全て取得し、新しい順に並び替え
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

// DBの一件取得
func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("could not database open")
	}
	var todo Todo
	//特定のレコードを取得
	db.First(&todo, id)
	db.Close()
	return todo
}

func main() {
	// Ginのルーターを作成
	router := gin.Default()
	// HTMLを読み込むディレクトリを指定
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	//Index
	// index.htmlにGetで繋いでいる
	router.GET("/", func(c *gin.Context) {
		todos := dbGetAll()
		c.HTML(200, "index.html", gin.H{
			"todos": todos,
		})
	})

	//Create
	router.POST("/new", func(c *gin.Context) {
		text := c.PostForm("text")
		status := c.PostForm("status")
		dbCreate(text, status)
		c.Redirect(302, "/")
	})

	//Detail
	router.GET("/detail/:id", func(c *gin.Context) {
		n := c.Param("id")
		// idの値をint型に変換
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		c.HTML(200, "detail.html", gin.H{"todo": todo})
	})

	//Update
	router.POST("/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		text := c.PostForm("text")
		status := c.PostForm("status")
		dbUpdate(id, text, status)
		c.Redirect(302, "/")
	})

	//削除確認
	router.GET("/delete_check/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		todo := dbGetOne(id)
		c.HTML(200, "delete.html", gin.H{"todo": todo})
	})

	//Delete
	router.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		c.Redirect(302, "/")
	})

	// ポートを指定していないので、デフォルトのポート(8080)で待受
	router.Run()
}
