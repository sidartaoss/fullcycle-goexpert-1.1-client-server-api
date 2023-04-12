package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const cambioUrl = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type Cambio struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type Cotacao struct {
	Bid string `json:"bid"`
}

type jsonError struct {
	Error string
}

func main() {

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	stmts := `
		DROP TABLE IF EXISTS cotacoes;
		CREATE TABLE cotacoes (code TEXT, codein TEXT, name TEXT, high TEXT, low TEXT, var_bid TEXT, pct_change TEXT, bid TEXT, ask TEXT, timestamp TEXT, create_date TEXT);
	`
	_, err = db.Exec(stmts)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/cotacao", BuscaCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request iniciado")
	defer log.Println("Request finalizado")
	select {
	case <-time.After(1 * time.Nanosecond):

		cambio, err := BuscaCambio()
		if err != nil {
			log.Println("chegou aqui com erro...")
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(jsonError{Error: err.Error()})
			return
		}

		err = insertCambio(cambio)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(jsonError{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Accept", "application/json")
		json.NewEncoder(w).Encode(Cotacao{Bid: cambio.Usdbrl.Bid})

	case <-ctx.Done():
		log.Println("Request cancelado pelo cliente")
	}

}

func BuscaCambio() (*Cambio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), (200 * time.Millisecond))
	defer cancel()
	select {
	case <-time.After(1 * time.Nanosecond):
		req, err := http.NewRequestWithContext(ctx, "GET", cambioUrl, nil)
		if err != nil {
			return nil, err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		var data Cambio
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, err
		}
		return &data, nil
	case <-ctx.Done():
		return nil, errors.New("tempo excedido para a chamada do endpoint de cambio do dolar")
	}
}

func insertCambio(c *Cambio) error {

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	select {
	case <-time.After(1 * time.Nanosecond):
		stmt, err := db.PrepareContext(ctx, "insert into cotacoes(code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) values (?,?,?,?,?,?,?,?,?,?,?)")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(c.Usdbrl.Code, c.Usdbrl.Codein, c.Usdbrl.Name, c.Usdbrl.High, c.Usdbrl.Low, c.Usdbrl.VarBid, c.Usdbrl.PctChange, c.Usdbrl.Bid, c.Usdbrl.Ask, c.Usdbrl.Timestamp, c.Usdbrl.CreateDate)
		if err != nil {
			return err
		}
		return nil
	case <-ctx.Done():
		return errors.New("banco de dados excedeu o tempo para persistir dados")
	}
}
