package models

type LivenessProbe struct {
	Protocol         string `json:"protocol"`
	InitialDelaySecs int    `json:"initial_delay_seconds"`
	SuccessThreshold int    `json:"success_threshold"`
	FailureThreshold int    `json:"failure_threshold"`
	Port             int    `json:"port"`
	PeriodSeconds    int    `json:"period_seconds"`
	TimeoutSeconds   int    `json:"timeout_seconds"`
	Path             string `json:"path,omitempty"`
	GRPCService      string `json:"grpc_service,omitempty"`
	Commands         string `json:"commands,omitempty"`
}

type ReadinessProbe struct {
	Protocol         string `json:"protocol"`
	InitialDelaySecs int    `json:"initial_delay_seconds"`
	SuccessThreshold int    `json:"success_threshold"`
	FailureThreshold int    `json:"failure_threshold"`
	Port             int    `json:"port"`
	PeriodSeconds    int    `json:"period_seconds"`
	TimeoutSeconds   int    `json:"timeout_seconds"`
	Path             string `json:"path,omitempty"`
	GRPCService      string `json:"grpc_service,omitempty"`
	Commands         string `json:"commands,omitempty"`
}

type AdvanceConfig struct {
	ImagePullPolicy         string         `json:"image_pull_policy"`
	IsReadinessProbeEnabled bool           `json:"is_readiness_probe_enabled"`
	IsLivenessProbeEnabled  bool           `json:"is_liveness_probe_enabled"`
	ReadinessProbe          ReadinessProbe `json:"readiness_probe"`
	LivenessProbe           LivenessProbe  `json:"liveness_probe"`
}

type PrivateImageDetails struct {
	RegistryNamespaceId int `json:"registry_namespace_id,omitempty"`
}

type EnvVariables struct {
	Key      string      `json:"key"`
	Value    string      `json:"value"`
	Disabled interface{} `json:"disabled"`
	Required bool        `json:"required"`
}

type Container struct {
	ContainerName       string              `json:"container_name"`
	ContainerType       string              `json:"container_type"`
	PrivateImageDetails PrivateImageDetails `json:"private_image_details"`
	AdvanceConfig       AdvanceConfig       `json:"advance_config"`
}

type ResourceDetails struct {
	DiskSize     int            `json:"disk_size"`
	MountPath    string         `json:"mount_path"`
	EnvVariables []EnvVariables `json:"env_variables"`
}

type Rules struct {
	Metric           string `json:"metric"`
	CustomMetricName string `json:"custom_metric_name,omitempty"`
	ConditionType    string `json:"condition_type"`
	Value            int    `json:"value"`
	WatchPeriod      int    `json:"watch_period"`
	Granularity      int    `json:"granularity"`
	Window           int    `json:"window"`
}

type CustomEndpointDetails struct {
	ServicePort     bool            `json:"service_port"`
	MetricPort      bool            `json:"metric_port"`
	Container       Container       `json:"container"`
	ResourceDetails ResourceDetails `json:"resource_details"`
	PublicIP        string          `json:"public_ip"`
}

type AutoScalePolicy struct {
	MinReplica      int         `json:"min_replica"`
	MaxReplica      int         `json:"max_replica"`
	StabilityPeriod int         `json:"stability_period"`
	Rules           interface{} `json:"rules"`
}

type ModelEndpoint struct {
	Action                  string                `json:"action,omitempty"`
	Name                    string                `json:"name"`
	Path                    string                `json:"path"`
	CustomEndpointDetails   CustomEndpointDetails `json:"custom_endpoint_details"`
	ModelID                 *int                  `json:"model_id"`
	SKUName                 string                `json:"sku_name,omitempty"`
	SkuType                 string                `json:"sku_type,omitempty"`
	CommittedDays           int                   `json:"committed_days"`
	CommittedInstancePolicy string                `json:"committed_instance_policy"`
	Replica                 int                   `json:"replica"`
	CommittedReplicas       int                   `json:"committed_replicas"`
	Framework               string                `json:"framework"`
	IsAutoScaleEnabled      bool                  `json:"is_auto_scale_enabled"`
	AutoScalePolicy         AutoScalePolicy       `json:"auto_scale_policy"`
	DetailedInfo            interface{}           `json:"detailed_info"`
	ModelLoadIntegrationID  *int                  `json:"model_load_integration_id"`
	DatasetID               *int                  `json:"dataset_id"`
	DatasetPath             string                `json:"dataset_path"`
	PrivateCloudID          *int                  `json:"private_cloud_id,omitempty"`
	CustomSKU               interface{}           `json:"custom_sku,omitempty"`
	ClusterType             string                `json:"cluster_type"`
	StorageType             string                `json:"storage_type"`
	SFSId                   int                   `json:"sfs_id,omitempty"`
	SFSPath                 string                `json:"sfs_path,omitempty"`
	DiskPath                string                `json:"disk_path,omitempty"`
	Category                string                `json:"category"`
	Location                string                `json:"location"`
	Currency                string                `json:"currency"`
}
