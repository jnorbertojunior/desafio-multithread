package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	canalBrasilApi := make(chan string)
	canalViacep := make(chan string)
	cep := "01153000"

	go fazerRequisicao(fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep), canalBrasilApi)
	go fazerRequisicao(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep), canalViacep)

	select {
	case res := <-canalBrasilApi:
		log.Println("BrasilAPI", res)

	case res := <-canalViacep:
		log.Println("ViaCEP", res)

	case <-time.After(time.Second):
		log.Println(errors.New("Timeout"))
	}
}

func fazerRequisicao(url string, canal chan string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	enderecoJson, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	canal <- string(enderecoJson)
}
