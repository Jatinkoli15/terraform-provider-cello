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

type Client struct {
	Api_key      string
	Auth_token   string
	Api_endpoint string
	HttpClient   *http.Client
}

func NewClient(api_key string, auth_token string, api_endpoint string) *Client {
	return &Client{
		Api_key:      api_key,
		Auth_token:   auth_token,
		Api_endpoint: api_endpoint,
		HttpClient:   &http.Client{},
	}
}

func (c *Client) NewNode(item *models.NodeCreate, teamID string, projectID string, activeIAM string) (map[string]interface{}, error) {

	jsonPayload, _ := json.Marshal(item)
	buf := bytes.NewBuffer(jsonPayload)
	UrlNode := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/notebooks/"
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

func (c *Client) GetNode(nodeId string, projectID, teamID string, activeIAM string) (map[string]interface{}, error) {

	urlNode := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/notebooks/" + nodeId + "/"
	req, err := http.NewRequest("GET", urlNode, nil)
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
		return nil, fmt.Errorf("404 Not Found: Node with ID %s does not exist", nodeId)
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("got a non-200 status code: %v - %s", response.StatusCode, string(body))
	}
	resBody, _ := io.ReadAll(response.Body)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBody, &jsonRes)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return jsonRes, nil
	}
	return jsonRes, nil
}

func (c *Client) DeleteNode(nodeId string, projectID string, teamID string, activeIAM string) error {

	urlNode := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/notebooks/" + nodeId + "/"
	req, err := http.NewRequest("DELETE", urlNode, nil)
	if err != nil {
		return err
	}
	log.Println("[INFO] CLIENT | NODE READ")
	log.Println(req)
	log.Println(req.Body)
	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	return nil
}

func (c *Client) UpdateStartStopNode(nodeId string, projectID string, teamID string, activeIAM string, start_stop_flag bool) (map[string]interface{}, error) {

	urlNode := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/notebooks/" + nodeId + "/" + "actions/"
	req, err := http.NewRequest("PUT", urlNode, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	req.Header.Set("User-Agent", "terraform/e2e")
	if start_stop_flag {
		params.Add("action", "stop")
	} else {
		params.Add("action", "start")
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, _ := io.ReadAll(response.Body)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBody, &jsonRes)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return jsonRes, nil
	}
	return jsonRes, nil
}

func (c *Client) UpdatePlanNode(item *models.NodeAction, projectID string, teamID string, activeIAM string, nodeId string) (map[string]interface{}, error) {
	jsonPayload, _ := json.Marshal(item)
	buf := bytes.NewBuffer(jsonPayload)
	urlNode := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/notebooks/" + nodeId + "/"
	req, err := http.NewRequest("PUT", urlNode, buf)
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
	log.Println("this is the reqset", req)
	log.Println("this is the plan change req", req.Body)
	response, err := c.HttpClient.Do(req)
	log.Println("after cycle", response.Body)
	log.Println("this is error", err)
	if err != nil {
		log.Println("Error")
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
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

func (c *Client) UpdateImage(item *models.ImageDetail, projectID string, teamID string, activeIAM string, nodeId string) (map[string]interface{}, error) {

	jsonPayload, _ := json.Marshal(item)
	buf := bytes.NewBuffer(jsonPayload)
	urlNode := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/notebooks/" + nodeId + "/image_update/"
	req, err := http.NewRequest("PUT", urlNode, buf)
	if err != nil {
		return nil, err
	}
	log.Println("this is the req", req)
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
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
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

func (c *Client) GetImages(activeIAM string) (map[string]interface{}, error) {
	urlNode := c.Api_endpoint + "/gpu_service/" + "image/"
	req, err := http.NewRequest("GET", urlNode, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT | NODE READ")
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	params.Add("category", "notebook")
	params.Add("is_jupyterlab_enabled", "true")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	log.Println(req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	log.Printf("%s", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		log.Printf("[ERROR] CLIENT GET NDE | error when unmarshalling")
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) GetPlans(item *models.ImageDetail, activeIAM string) (map[string]interface{}, error) {
	jsonPayload, _ := json.Marshal(item)
	buf := bytes.NewBuffer(jsonPayload)
	urlNode := c.Api_endpoint + "/gpu_service/" + "sku/"
	req, err := http.NewRequest("GET", urlNode, buf)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT | NODE READ")
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	params.Add("service", "notebook")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	log.Println(req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println(response)
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	log.Printf("%s", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		log.Printf("[ERROR] CLIENT GET NDE | error when unmarshalling")
		return nil, err
	}
	return jsonRes, nil
}

func CheckResponseStatus(response *http.Response) error {
	if response.StatusCode != http.StatusCreated {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	return nil
}
