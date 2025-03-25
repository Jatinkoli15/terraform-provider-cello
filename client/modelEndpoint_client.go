package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"github.com/e2eterraformprovider/terraform-provider-tir/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (c *Client) NewEndoint(item *models.ModelEndpoint, projectID string, teamID string, activeIAM string) (map[string]interface{}, error) {

	repoJSON, _ := json.Marshal(item)
	buf := bytes.NewBuffer(repoJSON)

	url := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/serving/inference/"
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	log.Println("jatin111111")
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	params.Add("prefix", "models%2F")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println("Error while creating")
		return nil, err
	}
	log.Println("response", response)
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

func (c *Client) GetEndpoint(endpointID string, projectID string, teamID string, activeIAM string) (map[string]interface{}, error) {

	url := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/serving/inference/" + endpointID + "/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	log.Println("jatin111111GET")
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println("Error while Reading")
		return nil, err
	}
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBody, &jsonRes)
	if err != nil {
		return nil, err
	}
	log.Println("json is ", jsonRes)
	return jsonRes, nil
}

func (c *Client) DeleteEndpoint(endpointID string, projectID string, teamID string, activeIAM string) (map[string]interface{}, error) {
	url := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/serving/inference/" + endpointID + "/"
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	log.Println("jatin111111DELETE")
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println("Error while Reading")
		return nil, err
	}
	defer response.Body.Close()
	return nil, nil
}

func (c *Client) UpdateStartStopInference(endpointID string, projectID string, teamID string, activeIAM string, start_stop_flag string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"action": start_stop_flag,
	}
	jsonData, _ := json.Marshal(payload)
	buf := bytes.NewBuffer(jsonData)
	urlNode := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/serving/inference/" + endpointID + "/"
	req, err := http.NewRequest("PUT", urlNode, buf)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	req.Header.Set("User-Agent", "terraform/e2e")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println(response)
	return nil, nil
}

func (c *Client) UpdateEndpoint(item *models.ModelEndpoint, projectID string, teamID string, activeIAM string, endpointID string) (map[string]interface{}, error) {
	repoJSON, _ := json.Marshal(item)
	buf := bytes.NewBuffer(repoJSON)
	log.Println("request json", buf)
	url := c.Api_endpoint + "/teams/" + teamID + "/projects/" + projectID + "/serving/inference/" + endpointID + "/"
	req, err := http.NewRequest("PUT", url, buf)
	if err != nil {
		return nil, err
	}
	log.Println("jatin Update", req)
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("active_iam", activeIAM)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform/e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println("Error while Reading", err)
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("some problem occured")
	}
	log.Println("response update client", response)
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBody, &jsonRes)
	if err != nil {
		log.Println("Unmarshal", err)
		return nil, err
	}
	log.Println("json is ", jsonRes)

	return jsonRes, nil
}

func SetSchemaFromResponse(d *schema.ResourceData, response map[string]interface{}) error {
	// Extract the "data" field from the response
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response format: 'data' field is missing or not a map")
	}

	// Set basic fields
	if err := d.Set("name", data["name"]); err != nil {
		return fmt.Errorf("failed to set 'name': %v", err)
	}
	if err := d.Set("status", data["status"]); err != nil {
		return fmt.Errorf("failed to set 'status': %v", err)
	}
	if err := d.Set("created_at", data["created_at"]); err != nil {
		return fmt.Errorf("failed to set 'created_at': %v", err)
	}
	// Set SKU-related fields
	skuDetails, ok := data["sku_details"].(map[string]interface{})
	if ok {
		if err := d.Set("sku_name", skuDetails["specs"].(map[string]interface{})["name"]); err != nil {
			return fmt.Errorf("failed to set 'sku_name': %v", err)
		}
		if err := d.Set("sku_type", skuDetails["plan"].(map[string]interface{})["sku_type"]); err != nil {
			return fmt.Errorf("failed to set 'sku_type': %v", err)
		}
		if err := d.Set("committed_days", skuDetails["plan"].(map[string]interface{})["committed_days"]); err != nil {
			return fmt.Errorf("failed to set 'sku_type': %v", err)
		}
		if err := d.Set("currency", skuDetails["plan"].(map[string]interface{})["currency"]); err != nil {
			return fmt.Errorf("failed to set 'sku_type': %v", err)
		}
	}

	// Set storage-related fields
	if err := d.Set("storage_type", data["storage_type"]); err != nil {
		return fmt.Errorf("failed to set 'storage_type': %v", err)
	}
	if err := d.Set("disk_path", data["disk_path"]); err != nil {
		return fmt.Errorf("failed to set 'disk_path': %v", err)
	}
	if err := d.Set("sfs_path", data["sfs_path"]); err != nil {
		return fmt.Errorf("failed to set 'sfs_path': %v", err)
	}

	// Set replica-related fields
	if err := d.Set("replica", data["replica"]); err != nil {
		return fmt.Errorf("failed to set 'replica': %v", err)
	}
	if err := d.Set("committed_replicas", data["committed_replicas"]); err != nil {
		return fmt.Errorf("failed to set 'committed_replicas': %v", err)
	}

	// Set auto-scaling policy
	autoScalePolicy, ok := data["auto_scale_policy"].(map[string]interface{})
	if ok {
		stabilityPeriod, _ := strconv.Atoi(autoScalePolicy["stability_period"].(string))
		autoScalePolicy["stability_period"] = stabilityPeriod
		autoScalePolicyList := []map[string]interface{}{autoScalePolicy}
		if err := d.Set("auto_scale_policy", autoScalePolicyList); err != nil {
			return fmt.Errorf("failed to set 'auto_scale_policy': %v", err)
		}
	}

	// Set detailed info
	detailedInfo, ok := data["detailed_info"].(map[string]interface{})
	engineArgs, _ := detailedInfo["engine_args"].(map[string]interface{})
	stringEngineArgs := make(map[string]string)
	for key, value := range engineArgs {
		// Convert the value to a string using fmt.Sprintf
		stringEngineArgs[key] = fmt.Sprintf("%v", value)
	}
	detailedInfo["engine_args"] = stringEngineArgs
	detailedInfo["args"] = ""
	detailedInfo["commands"] = ""
	if ok {
		detailedInfoList := []map[string]interface{}{detailedInfo}
		if err := d.Set("detailed_info", detailedInfoList); err != nil {
			return fmt.Errorf("failed to set 'detailed_info': %v", err)
		}
	}

	// Set container and probe configurations
	customEndpointDetails, ok := data["custom_endpoint_details"].(map[string]interface{})
	if ok {
		container, ok := customEndpointDetails["container"].(map[string]interface{})
		if ok {
			advanceConfig, ok := container["advance_config"].(map[string]interface{})
			if ok {
				if err := d.Set("is_readiness_probe_enabled", advanceConfig["is_readiness_probe_enabled"]); err != nil {
					return fmt.Errorf("failed to set 'is_readiness_probe_enabled': %v", err)
				}
				if err := d.Set("is_liveness_probe_enabled", advanceConfig["is_liveness_probe_enabled"]); err != nil {
					return fmt.Errorf("failed to set 'is_liveness_probe_enabled': %v", err)
				}

				readinessProbe, ok := advanceConfig["readiness_probe"].(map[string]interface{})
				if ok {
					port, _ := strconv.Atoi(readinessProbe["port"].(string))
					readinessProbe["port"] = port
					readinessProbeList := []map[string]interface{}{readinessProbe}
					if err := d.Set("readiness_probe", readinessProbeList); err != nil {
						return fmt.Errorf("failed to set 'readiness_probe': %v", err)
					}
				}

				// livenessProbe, ok := advanceConfig["liveness_probe"].(map[string]interface{})
				// if ok {
				// 	port, _ := strconv.Atoi(livenessProbe["port"].(string))
				// 	livenessProbe["port"] = port
				//     livenessProbeList := []map[string]interface{}{livenessProbe}
				//     if err := d.Set("liveness_probe", livenessProbeList); err != nil {
				//         return fmt.Errorf("failed to set 'liveness_probe': %v", err)
				//     }
				// }
			}
		}
	}

	// Set resource details
	resourceDetails, ok := customEndpointDetails["resource_details"].(map[string]interface{})
	if ok {
		resourceDetailsList := []map[string]interface{}{resourceDetails}
		if err := d.Set("resource_details", resourceDetailsList); err != nil {
			return fmt.Errorf("failed to set 'resource_details': %v", err)
		}
	}

	// Set public IP
	if err := d.Set("public_ip", customEndpointDetails["public_ip"]); err != nil {
		return fmt.Errorf("failed to set 'public_ip': %v", err)
	}

	if d.Get("status") == "stopped" {
		d.Set("stop_inference", "stop")
	} else {
		d.Set("stop_inference", "start")
	}

	return nil
}
