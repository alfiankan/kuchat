package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func createUser(db *sql.DB, email, password string) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("registering", email, password)
	row := tx.QueryRow("INSERT INTO users (email, password) VALUES ($1, crypt($2, gen_salt('bf')) ) returning id, password", email, password)

	var id, passwd string
	if err := row.Scan(&id, &passwd); err != nil {
		log.Println(err)
		tx.Rollback()
	}

	if _, err = tx.Exec(`
	INSERT INTO vmq_auth_acl 
	(mountpoint, client_id, username, password, publish_acl, subscribe_acl) 
	VALUES ('', $1, $2, $3, $4, $5)
	`, id, email, passwd, `[{"pattern": "a/b/c"}, {"pattern": "c/b/#"}]`, `[{"pattern": "a/b/c"}, {"pattern": "c/b/#"}]`); err != nil {
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
	}
	return
}

func main() {
	e := echo.New()

	db, err := NewPgConnection("host=127.0.0.1 user=postgres password=postgres dbname=kuchat port=5432 sslmode=disable TimeZone=Asia/Jakarta")
	if err != nil {
		log.Fatal(err)
	}

	e.POST("/register", func(c echo.Context) error {

		var reqBody UserRequest

		if err := c.Bind(&reqBody); err != nil {
			log.Println(err)
			return c.String(http.StatusBadRequest, "Gagal")
		}
		log.Println(reqBody)

		if err := createUser(db, reqBody.Email, reqBody.Password); err != nil {
			log.Println(err)
			return c.String(http.StatusBadRequest, "Gagal")
		}

		return c.String(http.StatusOK, "Registered")
	})

	e.Logger.Fatal(e.Start(":1323"))

}
