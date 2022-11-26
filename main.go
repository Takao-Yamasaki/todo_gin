package main

import (
	_ "text/tabwriter"

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

func main() {
	// Ginのルーターを作成
	router := gin.Default()
	// HTMLを読み込むディレクトリを指定
	router.LoadHTMLGlob("templates/*.html")

	data := "Hello Go/Gin!!"

	// index.htmlにGetで繋いでいる
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"data": data})
	})

	// ポートを指定していないので、デフォルトのポート(8080)で待受
	router.Run()
}
