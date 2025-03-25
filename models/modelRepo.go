package models

type ModelRepo struct {
	ModelType   string `json:"model_type"`
	Name        string `json:"name"`
	BucketName  string `json:"bucket_name"`
	SecretKey   string `json:"secret_key"`
	AccessKey   string `json:"access_key"`
	StorageType string `json:"storage_type"`
}
