package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var jsonError map[string]string
	if err = json.Unmarshal(body, &jsonError); err != nil {
		panic(err)
	}
	for k, v := range jsonError {
		switch k {
		case "Error":
			panic(v)
		}
	}
	var cotacao Cotacao
	if err = json.Unmarshal(body, &cotacao); err != nil {
		panic(err)
	}
	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %v", cotacao.Bid))
	if err != nil {
		panic(err)
	}
	fmt.Println("Cotação atual salva com sucesso")
}
