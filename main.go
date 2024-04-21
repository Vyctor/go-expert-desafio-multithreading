package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type BrasilApiCEPDataOutput struct {
	CEP          string `json:"cep"`
	STATE        string `json:"state"`
	CITY         string `json:"city"`
	NEIGHBORHOOD string `json:"neighborhood"`
	STREET       string `json:"street"`
	SERVICE      string `json:"service"`
}

type ViaApiCEPDataOutput struct {
	CEP          string `json:"cep"`
	STREET       string `json:"logradouro"`
	COMPLEMENT   string `json:"complemento"`
	NEIGHBORHOOD string `json:"bairro"`
	CITY         string `json:"localidade"`
	STATE        string `json:"uf"`
	IBGE         string `json:"ibge"`
	GIA          string `json:"gia"`
	DDD          string `json:"ddd"`
	SIAFI        string `json:"siafi"`
}

func main() {
	cep := os.Args[1]

	if len(cep) != 8 {
		fmt.Println("CEP inválido")
		return
	}

	fmt.Println("Buscando CEP: ", cep)

	brasilApiChannel := make(chan *BrasilApiCEPDataOutput)
	viaCepChannel := make(chan *ViaApiCEPDataOutput)

	go func() {
		brasilApiResponse, brasilApiErr := GetAddressByBrasilApiRequest(cep)
		if brasilApiErr != nil {
			fmt.Printf("Erro ao buscar CEP: %s\n", brasilApiErr)
			return
		}
		brasilApiChannel <- brasilApiResponse
	}()

	go func() {
		viaCepApiResponse, viaCepErr := GetAddressByViaCepApiRequest(cep)
		if viaCepErr != nil {
			fmt.Printf("Erro ao buscar CEP: %s\n", viaCepErr)
			return
		}
		viaCepChannel <- viaCepApiResponse
	}()

	select {
	case brasilApiResponse := <-brasilApiChannel:
		if brasilApiResponse != nil {
			fmt.Println("Endereço encontrado pela Brasil API")
			fmt.Println("CEP: ", brasilApiResponse.CEP)
			fmt.Println("Estado: ", brasilApiResponse.STATE)
			fmt.Println("Cidade: ", brasilApiResponse.CITY)
			fmt.Println("Bairro: ", brasilApiResponse.NEIGHBORHOOD)
			fmt.Println("Rua: ", brasilApiResponse.STREET)
		}
	case viaCepApiResponse := <-viaCepChannel:
		if viaCepApiResponse != nil {
			fmt.Println("Endereço encontrado pela Via CEP API")
			fmt.Println("CEP: ", viaCepApiResponse.CEP)
			fmt.Println("Estado: ", viaCepApiResponse.STATE)
			fmt.Println("Cidade: ", viaCepApiResponse.CITY)
			fmt.Println("Bairro: ", viaCepApiResponse.NEIGHBORHOOD)
			fmt.Println("Rua: ", viaCepApiResponse.STREET)
		}
	case <-time.After(1 * time.Second):
		fmt.Println("Erro! Nenhum endereço encontrado")
	}

}

func GetAddressByBrasilApiRequest(cep string) (*BrasilApiCEPDataOutput, error) {
	brasilApiReq, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)

	if err != nil {
		return nil, err
	}

	defer brasilApiReq.Body.Close()
	res, err := io.ReadAll(brasilApiReq.Body)

	if err != nil {
		return nil, err
	}

	var brasilApiResponse BrasilApiCEPDataOutput

	err = json.Unmarshal(res, &brasilApiResponse)

	if err != nil {
		return nil, err
	}

	return &brasilApiResponse, nil
}

func GetAddressByViaCepApiRequest(cep string) (*ViaApiCEPDataOutput, error) {
	viaCepApiReq, err := http.Get("http://viacep.com.br/ws/" + cep + "/json")

	if err != nil {
		return nil, err
	}

	defer viaCepApiReq.Body.Close()
	res, err := io.ReadAll(viaCepApiReq.Body)

	if err != nil {
		return nil, err
	}

	var viaCepApiResponse ViaApiCEPDataOutput

	err = json.Unmarshal(res, &viaCepApiResponse)

	if err != nil {
		return nil, err
	}

	return &viaCepApiResponse, nil
}
