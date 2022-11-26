package main

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	_ "gorm.io/driver/mysql"
	_ "gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Text   string
	Status string
}

func Env_load() {
	// .envを読み込む
	err := godotenv.Load(".env")
	if err != nil {
		panic("could not read .env file")
	}
}

func gormConnect() *gorm.DB {
	Env_load()

	DBMS := "mysql"
	USER := os.Getenv("MYSQL_USER")
	PASS := os.Getenv("MYSQL_PASSWORD")
	PROTOCOL := "tcp(127.0.0.1:3306)"
	DBNAME := "todo"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true"
	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}
	return db
}

// DBの初期化
func dbInit() {
	db := gormConnect()
	defer db.Close()
	//マイグレートを実行
	db.AutoMigrate(&Todo{})
}

// DBの作成
func dbCreate(text string, status string) {
	db := gormConnect()
	defer db.Close()

	db.Create(&Todo{Text: text, Status: status})
}

func dbUpdate(id int, text string, status string) {
	db := gormConnect()
	defer db.Close()

	var todo Todo
	db.First(&todo, id)
	//構造体Todoを呼んでいる
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
}

func dbDelete(id int) {
	db := gormConnect()
	defer db.Close()

	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
}

// DBの全件取得
func dbGetAll() []Todo {
	db := gormConnect()
	defer db.Close()

	// 要素を全て取得し、新しい順に並び替え
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	return todos
}

// DBの一件取得
func dbGetOne(id int) Todo {
	db := gormConnect()
	defer db.Close()

	var todo Todo
	//特定のレコードを取得
	db.First(&todo, id)
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
