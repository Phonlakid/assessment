//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	C "github.com/Phonlakid/assessment/controller"
	m "github.com/Phonlakid/assessment/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

var serverPort string = "2565"

func SetupServer() {
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		h := handler.NewApplication(db)

		e.POST("/expenses", C.CreateExpensesHandler)
		e.GET("/expenses/:id", C.GetExpensesHandler)
		e.PUT("/expenses/:id", C.UpdateExpenseHandler)
		e.GET("/expenses", C.GetExpenseHandler)

		e.Start(":" + serverPort)
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%s", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

}
func TestCreateexpenses(t *testing.T) {
	// Setup server
	SetupServer()

	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
	var e m.Expenses

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, float64(79), e.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)
}

func TestGetexpenses(t *testing.T) {
	// Setup server
	SetupServer()

	c := seedExpenses(t)

	var Exp m.Expenses
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&Exp)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, Exp.ID)
	assert.NotEmpty(t, Exp.Title)
	assert.NotEmpty(t, Exp.Amount)
	assert.NotEmpty(t, Exp.Note)
	assert.NotEmpty(t, Exp.Tags)
}

func TestUpdateexpenses(t *testing.T) {
	id := seedExpenses(t).ID
	e := m.Expenses{
		ID:     id,
		Title:  "Gundam",
		Amount: 1,
		Note:   "Gundam Freedom",
		Tags:   []string{"gadget", "shopping"},
	}
	payload, _ := json.Marshal(e)
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(id)), bytes.NewBuffer(payload))
	var info m.Expenses
	err := res.Decode(&info)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.Title, info.Title)
	assert.Equal(t, e.Amount, info.Amount)
	assert.Equal(t, e.Note, info.Note)
	assert.Equal(t, e.Tags, info.Tags)
}

func TestGetAllExpenses(t *testing.T) {
	SetupServer()

	seedExpenses(t)
	var exps []m.Expenses

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&exps)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(exps), 0)
}

func seed(t *testing.T) m.Expenses {
	var e m.Expenses
	body := bytes.NewBufferString(`{
		"title": "PS5",
		"amount": 13999,
		"note": "god of war only", 
		"tags": ["gadget", "shopping"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&e)
	if err != nil {
		t.Fatal("can't create expenses:", err)
	}
	return e
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	jsonStr := StreamToString(r.Body)
	return json.Unmarshal([]byte(jsonStr), v)
}
func request(method, url string, body io.Reader) *Response {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "November 10, 2009")
	res, err := http.DefaultClient.Do(req)
	return &Response{res, err}
}

func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}
