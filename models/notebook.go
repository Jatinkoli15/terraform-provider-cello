package models

type NodeCreate struct {
	Name                    string   `json:"name"`
	ClusterType             string   `json:"cluster_type"`
	DiskSizeInGB            int      `json:"disk_size_in_gb"`
	EnableSSH               bool     `json:"enable_ssh"`
	ImageType               string   `json:"image_type"`
	ImageName               string   `json:"image_name"`
	ImageVersion            string   `json:"image_version"`
	InstanceType            string   `json:"instance_type"`
	IsJupyterLabEnabled     bool     `json:"is_jupyterlab_enabled"`
	NotebookType            string   `json:"notebook_type"`
	NotebookURL             string   `json:"notebook_url"`
	SfsPath                 string   `json:"sfs_path"`
	SKUName                 string   `json:"sku_name"`
	Location                string   `json:"location"`
	SKUType                 string   `json:"sku_type"`
	AddOns                  []string `json:"add_ons"`
	DatasetIDList           []string `json:"dataset_id_list"`
	CommittedDays           int      `json:"committed_days"`
	Currency                string   `json:"currency"`
	Category                string   `json:"category"`
	CommittedInstancePolicy string   `json:"committed_instance_policy"`
	PublicSSHKeys           []string `json:"public_key"`
}

// type Notebook struct {
// 	AddOns                   []string `json:"add_ons"`
// 	ClusterType              string        `json:"cluster_type"`
// 	// CommittedInstancePolicy  string        `json:"committed_instance_policy"`
// 	DatasetIDList            []string `json:"dataset_id_list"`
// 	DiskSizeInGB             int           `json:"disk_size_in_gb"`
// 	EnableSSH                bool          `json:"enable_ssh"`
// 	ImageType                string        `json:"image_type"`
// 	ImageVersionID           int           `json:"image_version_id"`
// 	InstanceType             string        `json:"instance_type"`
// 	IsJupyterLabEnabled      bool          `json:"is_jupyterlab_enabled"`
// 	Name                     string        `json:"name"`
// 	// NextSkuItemPriceID       int           `json:"next_sku_item_price_id"`
// 	NotebookType             string        `json:"notebook_type"`
// 	NotebookURL              string        `json:"notebook_url"`
// 	SfsPath                  string        `json:"sfs_path"`
// 	SkuID                    int           `json:"sku_id"`
// 	SkuItemPriceID           int           `json:"sku_item_price_id"`
// }

type NodeAction struct {
	CommittedInstancePolicy string `json:"committed_instance_policy"`
	SKUType                 string `json:"sku_type"`
	CommittedDays           int    `json:"committed_days"`
	Category                string `json:"category"`
	Currency                string `json:"currency"`
	Location                string `json:"location"`
	SKUName                 string `json:"sku_name"`
}

type ImageDetail struct {
	ImageName           string `json:"image_name"`
	ImageVersion        string `json:"image_version"`
	IsJupyterLabEnabled bool   `json:"is_jupyterlab_enabled"`
	ImageType           string `json:"image_type"`
}
