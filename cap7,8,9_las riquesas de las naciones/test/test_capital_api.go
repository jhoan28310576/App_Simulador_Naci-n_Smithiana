package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://localhost:8080/api/simular-retorno"
	payload := map[string]interface{}{
		"capital":       10000,
		"sector":        "comercio",
		"riesgo":        3,
		"competencia":   10,
		"salario_medio": 0,
		"pais":          "USA",
	}
	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error en la petición:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Código de estado:", resp.StatusCode)
	fmt.Println("Respuesta:")
	fmt.Println(string(body))
}
