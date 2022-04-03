package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type AddUrlRequest struct {
	Url      string `json:"url"`
	ExpireAt string `json:"expireAt"`
}

var db *sql.DB

func handleRedirect(context *gin.Context) {
	const layout = "2006-01-02T15:04:05Z"
	var url, expireAt string

	id := context.Param("ID")
	err := db.QueryRow("SELECT Url, ExpireAt FROM `shortURL` WHERE ID=?", id).Scan(&url, &expireAt)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": 404})
		return
	}

	expire, err := time.Parse(layout, expireAt)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": 500})
		fmt.Println(err)
		return
	}

	if time.Now().After(expire) {
		context.JSON(http.StatusNotFound, gin.H{"status": 404})
	} else {
		context.Redirect(302, url)
	}
}

func handleNewURL(context *gin.Context) {
	request := AddUrlRequest{}
	err := context.BindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": 400})
		return
	}
	result, err := db.Exec("INSERT INTO `shortURL` (Url, ExpireAt) VALUES (?, ?)", request.Url, request.ExpireAt)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": 500})
		fmt.Println(err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": 500})
		fmt.Println(err)
		return
	}
	shortURL := fmt.Sprintf("http://localhost/%d", id)
	context.JSON(http.StatusOK, gin.H{"id": id, "shortUrl": shortURL})
}

func main() {
	_db, err := sql.Open("sqlite3", "./url.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db = _db
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `shortURL` (`ID` INTEGER PRIMARY KEY AUTOINCREMENT, `Url` MEDIUMTEXT, `ExpireAt` CHAR(20));")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	err = router.SetTrustedProxies(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	router.GET("/:ID", handleRedirect)
	router.POST("/api/v1/urls", handleNewURL)
	router.Run(":80")
}
