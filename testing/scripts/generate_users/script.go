package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-faker/faker/v4"
)

type CreateUserReq struct {
	Username string `faker:"len=8" json:"username"`
	Password string `faker:"len=20" json:"password"`
}

const (
	HOST = "localhost"
	PORT = "8080"
)

var (
	url = fmt.Sprintf("http://%s:%s/users/new", HOST, PORT)
)

func main() {

	numUser := 1000

	f, err := os.OpenFile("users.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)

	for range numUser {
		sendUserRequest(w)
	}

}

func sendUserRequest(w *csv.Writer) {
	req := CreateUserReq{}

	if err := faker.FakeData(&req); err != nil {
		panic(err)
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	reqBody := bytes.NewBuffer(jsonBody)

	resp, err := http.Post(url, "application/json", reqBody)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	if resp.StatusCode != 201 {
		fmt.Println("user not created")
	}
	w.Write([]string{req.Username, req.Password})
	w.Flush()
}
