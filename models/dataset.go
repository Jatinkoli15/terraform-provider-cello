package models

type PVCDetails struct {
	DiskSize int    `json:"disk_size"`
	PvcType  string `json:"pvc_type"`
}

type Dataset struct {
	Name              string      `json:"name"`
	Encryption_Enable bool        `json:"encryption_enable"`
	Encryption_Type   *string     `json:"encryption_type,omitempty"`
	StorageType       string      `json:"storage_type"`
	Pvc               *PVCDetails `json:"pvc,omitempty"`
	BucketName        *string     `json:"bucket_name,omitempty"`
}
