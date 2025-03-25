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

func (c *Client) NewRepo(item *models.ModelRepo, projectID string, teamID string, activeIAM string) (map[string]interface{}, error) {

	repoJSON, _ := json.Marshal(item)
	buf := bytes.NewBuffer(repoJSON)

	url := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/serving/model/"
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	log.Println("jatin111111")
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println("Error while creating")
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

func (c *Client) GetRepo(repoID string, projectID string, teamID string, activeIAM string) (map[string]interface{}, error) {

	url := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/serving/model/" + repoID + "/"
	req, err := http.NewRequest("GET", url, nil)
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
	log.Println("[INFO] Before Calling API")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println("Error while Reading")
		return nil, err
	}
	log.Println("[INFO] After Calling API")
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBody, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DeleteRepo(repoID string, projectID string, teamID string, activeIAM string) (map[string]interface{}, error) {
	url := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/serving/model/" + repoID + "/"
	req, err := http.NewRequest("DELETE", url, nil)
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
	log.Println("[INFO] Before API Call for DELETION")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println("Error while Reading")
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("some problem occured")
	}
	log.Println("[INFO] After API Call for DELETION")
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBody, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}
