package models



type PrivateCluster struct { 
	Name string `json:"name"`
	NodesCount int `json:"nodes_count"`
	SKUName string `json:"sku_name"`
	SKUType string `json:"sku_type"`
	CommittedDays int `json:"committed_days"`
	CommittedInstancePolicy string `json:"committed_instance_policy"`
	Location string `json:"location"`
	Currency string `json:"currency"`
	Category string `json:"category"`
}