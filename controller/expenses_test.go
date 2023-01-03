//go:build unit
// +build unit

package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	m "github.com/Phonlakid/assessment/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateexpenses(t *testing.T) {
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
	seedExpenses(t)
	var exps []m.Expenses

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&exps)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(exps), 0)
}

func seedExpenses(t *testing.T) m.Expenses {
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
	return json.NewDecoder(r.Body).Decode(v)
}
