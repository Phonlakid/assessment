package controller

import (
	"database/sql"
	"net/http"

	"github.com/Phonlakid/assessment/db"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type User struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

func CreateexpensesHandler(c echo.Context) error {

	u := User{}
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.Conn.QueryRow("INSERT INTO expenses (title,amount,note,tags) values ($1, $2, $3, $4)  RETURNING id", u.Title, u.Amount, u.Note, pq.Array(u.Tags))
	err = row.Scan(&u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, u)
}

func GetUserHandler(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.Conn.Prepare("SELECT * FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query user statment:" + err.Error()})
	}

	row := stmt.QueryRow(id)
	u := User{}
	err = row.Scan(&u.ID, &u.Title, &u.Amount, &u.Note, (*pq.StringArray)(&u.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "user not found"})
	case nil:
		return c.JSON(http.StatusOK, u)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
	}
}
