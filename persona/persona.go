package persona

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const APIKEY = "persona_sandbox_22f0d4e8-31ab-4a10-9316-84bb9a6180e3"
const URL = "https://withpersona.com/api/v1/"

func CreateAccount(email string) (string, error) {

	reqBody := map[string]interface{}{
		"email-address": email,
	}

	postBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	b, err := _makePersonaRequest(
		"accounts", 
		"POST", 
		"application/json", 
		postBody,
	)

	resBody := map[string]interface{}{}
	if err := json.Unmarshal(b, &resBody); err != nil {
		return "", err
	}
	data := resBody["data"].(map[string]interface{})

	actId := data["id"].(string)
	return actId, nil
}

func GetAccount(accountId string) (error) {
	url := "https://withpersona.com/api/v1/accounts/" + accountId

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Persona-Version", "2023-01-05")
	req.Header.Add("Authorization", "Bearer " + APIKEY)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	resBody := map[string]interface{}{}
	if err := json.Unmarshal(body, &resBody); err != nil {
		return err
	}
	data := resBody["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})
	fmt.Println(attributes)
	return nil
}

func _makePersonaRequest(path, method, contentType string, body []byte) ([]byte, error) {

	req, err := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("accept", "application/json")
	req.Header.Add("Persona-Version", "2023-01-05")
	req.Header.Add("content-type", contentType)
	req.Header.Add("Authorization", "Bearer " + APIKEY)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}