package modelEndpoint

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"github.com/e2eterraformprovider/terraform-provider-tir/constants"
	"github.com/e2eterraformprovider/terraform-provider-tir/models"
	"github.com/e2eterraformprovider/terraform-provider-tir/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceModel() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource. This is a required field and must be unique within the project.",
			},
			"server_options": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Specifies the server options for the resource. This is typically used for server types like TRITON, PYTORCH, NEMO, and TENSOR RT.",
			},
			"sku_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SKU (Stock Keeping Unit) name for the resource. This defines the type of resource being deployed.",
			},
			"sku_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SKU type for the resource. This defines the category or classification of the SKU.",
			},
			"committed_instance_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The policy for committed instances. This defines how committed instances are managed and billed.",
			},
			"committed_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The number of days the instance is committed for. This is used for billing and resource allocation.",
			},
			"model_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The path to the model file or directory. This is used to specify the location of the model to be deployed.",
			},
			"framework": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The framework used for the model. This could be TensorFlow, PyTorch, etc.",
			},
			"model_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The unique identifier for the model. This is used to reference the model in the system.",
			},
			"model_load_integration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The integration ID used for loading the model. This is typically used for custom model loading workflows.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of cluster the resource is deployed on. ",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of storage used for the resource.",
			},
			"disk_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/mnt/models",
				Description: "The path where the disk is mounted. This is used to specify the location for model storage.",
			},
			"sfs_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/shared/.cache",
				Description: "The path for shared file storage. This is used for caching and shared resources.",
			},
			"sfs_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The ID of the shared file storage. This is used to reference the shared storage resource.",
			},
			"image_pull_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Always",
				ValidateFunc: validation.StringInSlice([]string{
					"Always",
					"IfNotPresent",
				}, false),
				Description: "The policy for pulling container images. Options are 'Always' or 'IfNotPresent'.",
			},
			"is_auto_scale_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether auto-scaling is enabled for the resource.",
			},
			"auto_scale_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Default:     nil,
				Description: "The policy for auto-scaling the resource. This includes min/max replicas and scaling rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_replicas": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The minimum number of replicas to maintain during auto-scaling.",
						},
						"max_replicas": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The maximum number of replicas to scale up to during auto-scaling.",
						},
						"rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Default:     nil,
							Description: "The rules for auto-scaling based on metrics and conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "The metric to monitor for auto-scaling",
									},
									"value": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     10,
										Description: "The threshold value for the metric to trigger scaling.",
									},
									"condition_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "limit",
										Description: "The type of condition to apply for scaling.",
									},
									"watch_period": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     60,
										Description: "The period (in seconds) to watch the metric before scaling.",
									},
									"granularity": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     1,
										Description: "The granularity of the metric data collection.",
									},
									"window": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     1,
										Description: "The time window (in seconds) for evaluating the metric.",
									},
									"custom_metric_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "The name of a custom metric to use for scaling.",
									},
								},
							},
						},
						"stability_period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     300,
							Description: "The period (in seconds) to wait after scaling before scaling again.",
						},
					},
				},
			},
			"replica": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The number of replicas to deploy for the resource.",
			},
			"committed_replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The number of replicas that are committed for the resource.",
			},
			"detailed_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Default:     nil,
				Description: "Detailed information about the resource, including commands, args, and logging settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"commands": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Commands to execute when the resource is deployed.",
						},
						"args": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Arguments to pass to the commands when the resource is deployed.",
						},
						"hugging_face_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The Hugging Face model ID associated with the resource.",
						},
						"tokenizer": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The tokenizer to use for the model.",
						},
						"server_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The version of the server being used.",
						},
						"world_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The world size for distributed training or inference.",
						},
						"error_log": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Enable or disable error logging.",
						},
						"info_log": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Enable or disable info logging.",
						},
						"warning_log": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Enable or disable warning logging.",
						},
						"log_verbose_level": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The verbosity level for logging.",
						},
						"model_serve_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The type of model serving (e.g., real-time, batch).",
						},
						"engine_args": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Additional engine-specific arguments for the model.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"is_readiness_probe_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable or disable the readiness probe for the resource.",
			},
			"is_liveness_probe_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable or disable the liveness probe for the resource.",
			},
			"readiness_probe": {
				Type:        schema.TypeList,
				Optional:    true,
				Default:     nil,
				Description: "Configuration for the readiness probe.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "http",
							Description: "The protocol to use for the readiness probe (e.g., http, tcp).",
						},
						"initial_delay_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "The initial delay (in seconds) before the readiness probe starts.",
						},
						"success_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The number of successful probes required to mark the resource as ready.",
						},
						"failure_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     3,
							Description: "The number of failed probes before the resource is marked as not ready.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     8080,
							Description: "The port to use for the readiness probe.",
						},
						"period_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "The period (in seconds) between readiness probe checks.",
						},
						"timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "The timeout (in seconds) for the readiness probe.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "/health",
							Description: "The path to check for the readiness probe.",
						},
						"grpc_service": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The gRPC service to check for the readiness probe.",
						},
						"commands": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Commands to execute for the readiness probe.",
						},
					},
				},
			},
			"liveness_probe": {
				Type:        schema.TypeList,
				Optional:    true,
				Default:     nil,
				Description: "Configuration for the liveness probe.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "http",
							Description: "The protocol to use for the liveness probe (e.g., http, tcp).",
						},
						"initial_delay_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "The initial delay (in seconds) before the liveness probe starts.",
						},
						"success_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The number of successful probes required to mark the resource as live.",
						},
						"failure_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     3,
							Description: "The number of failed probes before the resource is marked as not live.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     8080,
							Description: "The port to use for the liveness probe.",
						},
						"period_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "The period (in seconds) between liveness probe checks.",
						},
						"timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "The timeout (in seconds) for the liveness probe.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "/metrics",
							Description: "The path to check for the liveness probe.",
						},
						"grpc_service": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The gRPC service to check for the liveness probe.",
						},
						"commands": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Commands to execute for the liveness probe.",
						},
					},
				},
			},
			"resource_details": {
				Type:        schema.TypeList,
				Optional:    true,
				Default:     nil,
				Description: "Additional details about the resource, such as disk size, mount path, and environment variables.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     100,
							Description: "The size of the disk (in GB) allocated for the resource.",
						},
						"mount_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The path where the disk is mounted.",
						},
						"env_variables": {
							Type:        schema.TypeList,
							Optional:    true,
							Default:     nil,
							Description: "Environment variables to be set for the resource.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key for the environment variable.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The value for the environment variable.",
									},
									"required": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Indicates whether the environment variable is required.",
									},
									"disabled": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "A map of disabled environment variables.",
										Elem: &schema.Schema{
											Type: schema.TypeBool,
										},
									},
								},
							},
						},
					},
				},
			},
			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "Indicates whether a public IP address is assigned to the resource.",
			},
			"container_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the container associated with the resource. This is computed automatically.",
			},
			"container_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of container used for the resource (e.g., public, private).",
			},
			"private_cloud_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The ID of the private cloud where the resource is deployed.",
			},
			"custom_sku": {
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     nil,
				Description: "A map of custom SKU configurations for the private cloud .",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"service_port": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Indicates whether a service port is exposed for the resource.",
			},
			"metric_port": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether a metric port is exposed for the resource.",
			},
			"dataset_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The ID of the dataset associated with the resource.",
			},
			"dataset_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The path to the dataset used by the resource.",
			},
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the team that owns the resource.",
			},
			"active_iam": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IAM (Identity and Access Management) role associated with the resource.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project where the resource is deployed.",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location or region where the resource is deployed.",
			},
			"currency": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The currency used for billing the resource.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the resource. This is computed automatically.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created. This is computed automatically.",
			},
			"stop_inference": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "start",
				Description: "Indicates whether to stop or start inference for the resource. Default is 'start'.",
			},
		},
		CreateContext: resourceCreateModelEndpoint,
		ReadContext:   resourceReadModelEndpoint,
		UpdateContext: resourceUpdateModelEndpoint,
		DeleteContext: resourceDeleteModelEndpoint,
	}
}

func resourceCreateModelEndpoint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	log.Println("Starting resourceCreateModelEndpoint function")
	if d.Get("stop_inference") != "start" {
		return diag.Errorf("Field stop_inference must be [start] at the time of creation")
	}
	di := d.Get("detailed_info").([]interface{})
	detailed_info := di[0].(map[string]interface{})
	err, containerName := constants.GetContainerName(detailed_info["server_version"].(string), d.Get("model_id").(string), d.Get("framework").(string))
	log.Println("Container name:", containerName)

	if err != nil {
		return diag.Errorf("Error finding the framework, please enter the correct framework")
	}
	detailedInfoList := d.Get("detailed_info").([]interface{})
	detailedInfo := detailedInfoList[0].(map[string]interface{})
	originalEngineArgs := detailedInfo["engine_args"].(map[string]interface{})
	originalCommands := detailedInfo["commands"].(string)
	originalArgs := detailedInfo["args"].(string)

	_, endpointNode := createPayloadForInference(d)
	// log.Println("Repository JSON:", buf)

	response, error := apiClient.NewEndoint(&endpointNode, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if error != nil {
		return diag.Errorf("Some error occurred while creating the model repository. Please check the config you have provided!! %s", error)
	}
	data := response["data"].(map[string]interface{})
	if nodeID, ok := data["id"].(float64); ok {
		d.SetId(strconv.Itoa(int(math.Round(nodeID))))
		log.Println("ID SET")
	} else {
		return diag.Errorf("failed to extract node ID from response")
	}

	log.Println(d.Id())
	d.Set("container_name", containerName)
	detailedInfo["engine_args"] = originalEngineArgs
	detailedInfo["commands"] = originalCommands
	detailedInfo["args"] = originalArgs
	d.Set("status", data["status"].(string))
	d.Set("created_at", data["created_at"].(string))
	log.Println("resourceCreateModelEndpoint completed successfully")
	var diags diag.Diagnostics
	return diags
}

func resourceReadModelEndpoint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	endpointID := d.Id()

	response, err := apiClient.GetEndpoint(endpointID, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			log.Println("[INFO] Repo not found, setting ID to empty")
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding item with id: %s - %v", endpointID, err)
	}
	if err := client.SetSchemaFromResponse(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUpdateModelEndpoint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	endpointID := d.Id()
	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	activeIAM := d.Get("active_iam").(string)
	action := d.Get("stop_inference").(string)

	if d.HasChange("framework") {
		oldFramework, _ := d.GetChange("framework")
		d.Set("framework", oldFramework)
		return diag.Errorf("You cannot change framework once created inference!!")
	}
	if d.HasChange("stop_inference") {
		_, err := apiClient.UpdateStartStopInference(endpointID, projectID, teamID, activeIAM, action)
		if err != nil {
			d.Set("stop_inference", "start")
			return diag.Errorf("Not able to stop/start node")
		}
	} else {
		detailedInfoList := d.Get("detailed_info").([]interface{})
		detailedInfo := detailedInfoList[0].(map[string]interface{})
		originalEngineArgs := detailedInfo["engine_args"].(map[string]interface{})
		originalCommands := detailedInfo["commands"].(string)
		originalArgs := detailedInfo["args"].(string)
		err, endpointNode := createPayloadForInference(d)
		if err != nil {
			log.Println(err)
			return diag.Errorf("Something went wrong creating payload,,,please check the config")
		}
		if d.Get("status") == "stopped" {
			endpointNode.Action = "update"
		} else {
			endpointNode.Action = "patch"
		}
		response, error := apiClient.UpdateEndpoint(&endpointNode, projectID, teamID, activeIAM, endpointID)
		if error != nil {
			return diag.Errorf("Something went wrong please check the config file %s", error)
		}
		detailedInfo["engine_args"] = originalEngineArgs
		detailedInfo["commands"] = originalCommands
		detailedInfo["args"] = originalArgs
		log.Println("updated response", response)
	}
	return diags
}

func resourceDeleteModelEndpoint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	endpointID := d.Id()

	_, err := apiClient.DeleteEndpoint(endpointID, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			log.Println("[INFO] Repo not found, setting ID to empty")
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding item with id: %s - %v", endpointID, err)
	}
	d.SetId("")
	return diags
}

func buildEnvVariablesFromSchema(resource_details map[string]interface{}) ([]models.EnvVariables, error) {

	var envVars []models.EnvVariables
	if rawEnvVars, ok := resource_details["env_variables"]; ok {
		envVarsList := rawEnvVars.([]interface{})
		for _, rawEnvVar := range envVarsList {
			envVarMap := rawEnvVar.(map[string]interface{})
			envVar := models.EnvVariables{
				Key:      envVarMap["key"].(string),
				Value:    envVarMap["value"].(string),
				Required: envVarMap["required"].(bool),
			}
			if disabled, ok := envVarMap["disabled"]; ok {
				disabledMap := disabled.(map[string]interface{})
				envVar.Disabled = disabledMap
			}
			envVars = append(envVars, envVar)
		}
	}

	return envVars, nil
}

func convertEngineArgs(engineArgsMap map[string]interface{}) (map[string]interface{}, error) {
	convertedMap := make(map[string]interface{})

	for key, value := range engineArgsMap {
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("engine_args value for key %s is not a string", key)
		}
		if intValue, err := strconv.Atoi(strValue); err == nil {
			convertedMap[key] = intValue
			continue
		}
		if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
			convertedMap[key] = floatValue
			continue
		}
		if boolValue, err := strconv.ParseBool(strValue); err == nil {
			convertedMap[key] = boolValue
			continue
		}

		convertedMap[key] = strValue
	}

	return convertedMap, nil
}

func createPayloadForInference(d *schema.ResourceData) (diag.Diagnostics, models.ModelEndpoint) {

	rp := d.Get("readiness_probe").([]interface{})
	lp := d.Get("liveness_probe").([]interface{})
	readiness_probe := rp[0].(map[string]interface{})
	liveness_probe := lp[0].(map[string]interface{})
	log.Println("Readiness probe:", readiness_probe)
	log.Println("Liveness probe:", liveness_probe)

	readinessProbeNode := models.ReadinessProbe{
		Protocol:         readiness_probe["protocol"].(string),
		InitialDelaySecs: readiness_probe["initial_delay_seconds"].(int),
		SuccessThreshold: readiness_probe["success_threshold"].(int),
		FailureThreshold: readiness_probe["failure_threshold"].(int),
		Port:             readiness_probe["port"].(int),
		PeriodSeconds:    readiness_probe["period_seconds"].(int),
		TimeoutSeconds:   readiness_probe["timeout_seconds"].(int),
		Path:             readiness_probe["path"].(string),
		GRPCService:      readiness_probe["grpc_service"].(string),
		Commands:         readiness_probe["commands"].(string),
	}
	log.Println("Readiness probe node created:", readinessProbeNode)

	livenessProbeNode := models.LivenessProbe{
		Protocol:         liveness_probe["protocol"].(string),
		InitialDelaySecs: liveness_probe["initial_delay_seconds"].(int),
		SuccessThreshold: liveness_probe["success_threshold"].(int),
		FailureThreshold: liveness_probe["failure_threshold"].(int),
		Port:             liveness_probe["port"].(int),
		PeriodSeconds:    liveness_probe["period_seconds"].(int),
		TimeoutSeconds:   liveness_probe["timeout_seconds"].(int),
		Path:             liveness_probe["path"].(string),
		GRPCService:      liveness_probe["grpc_service"].(string),
		Commands:         liveness_probe["commands"].(string),
	}
	log.Println("Liveness probe node created:", livenessProbeNode)

	advancedConfigNode := models.AdvanceConfig{
		ImagePullPolicy:         d.Get("image_pull_policy").(string),
		IsReadinessProbeEnabled: d.Get("is_readiness_probe_enabled").(bool),
		IsLivenessProbeEnabled:  d.Get("is_liveness_probe_enabled").(bool),
		ReadinessProbe:          readinessProbeNode,
		LivenessProbe:           livenessProbeNode,
	}
	log.Println("Advanced config node created:", advancedConfigNode)

	rd := d.Get("resource_details").([]interface{})
	resource_details := rd[0].(map[string]interface{})
	log.Println("resource_details", resource_details)
	envVarNode, _ := buildEnvVariablesFromSchema(resource_details)
	resourceDetailsNode := models.ResourceDetails{
		DiskSize:     resource_details["disk_size"].(int),
		MountPath:    resource_details["mount_path"].(string),
		EnvVariables: envVarNode,
	}
	log.Println("Resource details node created:", resourceDetailsNode)
	di := d.Get("detailed_info").([]interface{})
	detailed_info := di[0].(map[string]interface{})
	err, containerName := constants.GetContainerName(detailed_info["server_version"].(string), d.Get("model_id").(string), d.Get("framework").(string))
	log.Println("Container name:", containerName)

	if err != nil {
		log.Println("Error getting container name:", err)
		return diag.Errorf("Something went wrong please check the config file"), models.ModelEndpoint{}
	}

	containerNode := models.Container{
		ContainerName:       containerName,
		ContainerType:       d.Get("container_type").(string),
		PrivateImageDetails: models.PrivateImageDetails{}, // if want to make private image then there is a field registry_namespace_id
		AdvanceConfig:       advancedConfigNode,
	}
	log.Println("Container node created:", containerNode)

	customEndpointDetailsNode := models.CustomEndpointDetails{
		ServicePort:     d.Get("service_port").(bool),
		MetricPort:      d.Get("metric_port").(bool),
		Container:       containerNode,
		ResourceDetails: resourceDetailsNode,
		PublicIP:        "no",
	}
	log.Println("Custom endpoint details node created:", customEndpointDetailsNode)
	frameName, _ := constants.GetFrameworkName(d.Get("framework").(string))

	// originalEngineArgs := detailed_info["engine_args"].(map[string]interface{})
	engine_args, _ := convertEngineArgs(detailed_info["engine_args"].(map[string]interface{}))
	detailed_info["engine_args"] = engine_args
	log.Println("before")
	commands := detailed_info["commands"].(string)
	args := detailed_info["args"].(string)
	log.Println("after")
	detailed_info["commands"] = base64.StdEncoding.EncodeToString([]byte(commands))
	detailed_info["args"] = base64.StdEncoding.EncodeToString([]byte(args))
	log.Println("engine_args", engine_args)
	log.Println("detailed_info", detailed_info)
	endpointNode := models.ModelEndpoint{
		Name:                   d.Get("name").(string),
		Path:                   d.Get("model_path").(string),
		Category:               "inference_service",
		CustomEndpointDetails:  customEndpointDetailsNode,
		Replica:                d.Get("replica").(int),
		CommittedReplicas:      d.Get("committed_replicas").(int),
		Framework:              frameName,
		IsAutoScaleEnabled:     d.Get("is_auto_scale_enabled").(bool),
		DetailedInfo:           detailed_info,
		DatasetPath:            d.Get("dataset_path").(string),
		ClusterType:            d.Get("cluster_type").(string),
		StorageType:            d.Get("storage_type").(string),
		SFSPath:                d.Get("sfs_path").(string),
		DiskPath:               d.Get("disk_path").(string),
		ModelLoadIntegrationID: nil,
		ModelID:                nil,
		PrivateCloudID:         nil,
		Location:               d.Get("location").(string),
		Currency:               d.Get("currency").(string),
	}
	log.Println("Endpoint node created:", endpointNode)

	if d.Get("model_id") != "" {
		log.Println("line 748", d.Get("model_id").(string))
		ModelID, _ := strconv.Atoi(d.Get("model_id").(string))
		endpointNode.ModelID = &ModelID
	} else if d.Get("model_load_integration_id") != "" {
		log.Println("line 752", d.Get("model_load_integration_id").(string))
		ModelIntegrationID, _ := strconv.Atoi(d.Get("model_load_integration_id").(string))
		endpointNode.ModelLoadIntegrationID = &ModelIntegrationID
	}

	if d.Get("private_cloud_id").(string) != "" {
		log.Println("line 758", d.Get("private_cloud_id").(string))
		private_cloud_id, _ := strconv.Atoi(d.Get("private_cloud_id").(string))
		endpointNode.PrivateCloudID = &private_cloud_id
		if d.Get("custom_sku") == nil {
			log.Println("Custom SKU is required for private cloud")
			return diag.Errorf("Please provide the custom sku for private cloud"), models.ModelEndpoint{}
		}
		CustomSku := d.Get("custom_sku").(map[string]interface{})
		endpointNode.CustomSKU = CustomSku
	} else {
		endpointNode.SKUName = d.Get("sku_name").(string)
		endpointNode.SkuType = d.Get("sku_type").(string)
		log.Println(d.Get("committed_days").(int))
		endpointNode.CommittedDays = d.Get("committed_days").(int)
		log.Println(endpointNode.CommittedDays)
		endpointNode.CommittedInstancePolicy = d.Get("committed_instance_policy").(string)
	}

	if d.Get("sfs_id").(string) != "" {
		SFSId, _ := strconv.Atoi(d.Get("sfs_id").(string))
		endpointNode.SFSId = SFSId
	}

	if d.Get("dataset_id").(string) != "" {
		DatasetID, _ := strconv.Atoi(d.Get("dataset_id").(string))
		endpointNode.DatasetID = &DatasetID
	}

	autoScalePolicyList := d.Get("auto_scale_policy").([]interface{})
	log.Println("autolist", autoScalePolicyList)
	autoScalePolicyMap, ok := autoScalePolicyList[0].(map[string]interface{})
	if !ok {
		return diag.Errorf("Please check the config file"), models.ModelEndpoint{}
	}
	log.Println("automap", autoScalePolicyMap)
	autoScalePolicyRulesList := autoScalePolicyMap["rules"].([]interface{})
	// log.Println("autorules",autoScalePolicyRulesList)
	autoScalePolicyModel := models.AutoScalePolicy{
		MinReplica:      autoScalePolicyMap["min_replicas"].(int),
		MaxReplica:      autoScalePolicyMap["max_replicas"].(int),
		StabilityPeriod: autoScalePolicyMap["stability_period"].(int),
		Rules:           autoScalePolicyRulesList,
	}
	endpointNode.AutoScalePolicy = autoScalePolicyModel
	// repoJSON, _ := json.Marshal(autoScalePolicyModel)
	// buf := bytes.NewBuffer(repoJSON)
	// log.Println("auto_scale",buf)

	return nil, endpointNode

}
