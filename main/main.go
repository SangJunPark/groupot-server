package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	db, err := sql.Open("mysql", "root:qwerty1!A@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + "groupot")
	if err != nil {
		print(err)
	}

	db.

	_, err = db.Exec("USE " + "groupot")
	if err != nil {
		print(err)
	}

	_, err = db.Exec("INSERT INTO example VALUES(1,'t')")
	if err != nil {
		print(err)
	}

	e := echo.New()
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	e.GET("/", hello)
	//e.GET("/logout", bye)
	e.POST("/logout", bye)
	e.Start(":8080")
	// e.DELETE("/logout", hello)

}

type (
	Coree interface {
		Hiroo() string
		Byroo() int
	}
)

// Handler
func hello(c echo.Context) error {
	fmt.Println("asdf")
	r := make(map[string]string)
	fmt.Println(c.QueryParam("id"))
	r["Hellow"] = "World" + c.QueryParam("id")
	data, _ := json.Marshal(r)
	t := c.Param("id")
	fmt.Println(t)
	return c.JSON(http.StatusOK, string(data))
}

// Handler
func bye(c echo.Context) error {
	fmt.Println("bye")
	a := c.QueryParam("id")

	return c.String(http.StatusOK, "bye, World! "+a)
}
