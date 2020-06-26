package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Responses struct {
	Responses []Response `json:"response"`
}

type Response struct {
	Created string `json:"created.timestamp"`
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

type Model struct {
	CreatedAt string `json:"created.timestamp"`
	ModelID   string `json:"modelId"`
	ID        string `json:"id"`
	Schema    struct {
		Inputs []Input `json:"input"`
	} `json:"schema"`
}

type Input struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://url/model", nil)
	req.Header.Add("Authorization", "Bearer token")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	responseMap := Responses{}
	err = json.Unmarshal(body, &responseMap)

	for _, i := range responseMap.Responses {
		modelName := i.ID
		get := fmt.Sprintf("https://url/model/%s", modelName)
		newReq, err := http.NewRequest("GET", get, nil)
		newReq.Header.Add("Authorization", "Bearer token")
		newResp, err := client.Do(newReq)
		if err != nil {
			panic(err)
		}
		defer newResp.Body.Close()

		body, err := ioutil.ReadAll(newResp.Body)

		if err != nil {
			panic(err)
		}

		newResponseMap := Model{}
		err = json.Unmarshal(body, &newResponseMap)

		newArray := []string{}

		for _, z := range newResponseMap.Schema.Inputs {
			name := z.Name
			nameLC := strings.ToLower(name)
			// fmt.Println(nameLC)

			if strings.Contains(nameLC, "npanxx") || strings.Contains(nameLC, "serviceability") || strings.Contains(nameLC, "zipcode") || strings.Contains(nameLC, "broadband") {
				// fmt.Println(nameLC)
				if nameLC != "zipcode_ivr" && nameLC != "vertical_broadband_flag" && nameLC != "zipcode" && nameLC != "zipcodescore" {
					newArray = append(newArray, nameLC)
				}
			}
		}
		if len(newArray) > 0 {
			fmt.Printf("%s: %v\n", modelName, newArray)
		}

	}

}
