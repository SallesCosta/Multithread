package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type CepParam struct {
	cep string `json:"cep"`
}

type ViaCEP struct {
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

type ApiCEP struct {
	Status   int    `json:"status"`
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

func main() {
	var cep string
	// cep := "87053259"
	c1 := make(chan ApiCEP)
	c2 := make(chan ViaCEP)

	fmt.Println("Por favor digite o cep...")
	fmt.Scan(&cep)

	var s = strings.SplitAfter(cep, string([]rune(cep)[4]))
	var cept = s[0] + "-" + s[1]

	// ApiCEP
	go func() {
		req, err := http.Get("https://cdn.apicep.com/file/apicep/" + cept + ".json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		}

		defer req.Body.Close()

		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler resposta %v\n", err)
		}

		var data ApiCEP
		err = json.Unmarshal(res, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer o parse da resposta: %v\n", err)
		}
		c1 <- data
	}()

	//ViaCEP
	go func() {
		req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		}

		defer req.Body.Close()

		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler resposta %v\n", err)
		}

		var data ViaCEP
		err = json.Unmarshal(res, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer o parse da resposta: %v\n", err)
		}
		c2 <- data
	}()

	select {
	case data := <-c1:
		fmt.Printf("API Source: ApiCEP\n")
		fmt.Printf("Data: %s\n", data)

	case data := <-c2:
		fmt.Printf("API Source: ViaCEP\n")
		fmt.Printf("Data: %s\n", data)

	case <-time.After(time.Second * 1):
		println("timeout")
	}
}
