package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ApiCep struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func findApiCep(cep string, ch chan<- ApiCep) error {
	req, err := http.Get(fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}
	var apiCep ApiCep
	err = json.Unmarshal(res, &apiCep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}
	ch <- apiCep
	return nil
}

func findViaCep(cep string, ch chan<- ViaCep) error {
	req, err := http.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}
	var viaCep ViaCep
	err = json.Unmarshal(res, &viaCep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}
	ch <- viaCep
	return nil
}

func main() {
	apiCepChanel := make(chan ApiCep)
	viaCepChanel := make(chan ViaCep)
	for _, cep := range os.Args[1:] {
		go findApiCep(cep, apiCepChanel)
		go findViaCep(cep, viaCepChanel)
	}

	select {
	case apiCep := <-apiCepChanel:
		fmt.Println("ApiCep: ", apiCep)
	case viaCep := <-viaCepChanel:
		fmt.Println("ViaCep: ", viaCep)
	case <-time.After(time.Second):
		println("timeout")
	}

}
