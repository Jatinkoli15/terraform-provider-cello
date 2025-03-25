package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/e2eterraformprovider/terraform-provider-tir/models"
)

func (c *Client) NewDataset(item *models.Dataset, projectID string, teamID string, activeIAM string) (map[string]interface{}, error) {
	jsonPayload, _ := json.Marshal(item)
	buf := bytes.NewBuffer(jsonPayload)

	UrlNode := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/datasets/"
	req, err := http.NewRequest("POST", UrlNode, buf)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println("Error")
		return nil, err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBody, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) GetDataset(datasetID string, projectID string, teamID string, activeIAM string) (map[string]interface{}, error) {
	urlNode := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/datasets/" + datasetID + "/"
	req, err := http.NewRequest("GET", urlNode, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] CLIENT NODE READ | request failed: %v", err)
		return nil, err
	}
	defer response.Body.Close()
	var jsonRes map[string]interface{}
	resBody, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(resBody, &jsonRes)
	if err != nil {
		log.Printf("[ERROR] CLIENT GET NODE | error when unmarshalling: %v", err)
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DeleteDataset(datasetID string, projectID string, teamID string, activeIAM string) (map[string]interface{}, error) {
	urlNode := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/datasets/" + datasetID + "/"
	req, err := http.NewRequest("DELETE", urlNode, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] CLIENT NODE READ | request failed: %v", err)
		return nil, err
	}
	defer response.Body.Close()
	log.Printf("[INFO] CLIENT NODE READ | response code: %d", response.StatusCode)
	if response.StatusCode == http.StatusNotFound {
		log.Println("[INFO] Node not found, returning explicit 404 error")
		return nil, fmt.Errorf("404 Not Found: Node with ID %s does not exist", datasetID)
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("got a non-200 status code: %v - %s", response.StatusCode, string(body))
	}
	var jsonRes map[string]interface{}
	resBody, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(resBody, &jsonRes)
	if err != nil {
		log.Printf("[ERROR] CLIENT GET NODE | error when unmarshalling: %v", err)
		return nil, err
	}
	return jsonRes, nil

}
