package model

type StartEvent struct {
	TaskID           string `json:"task_id"`
	FileS3Path       string `json:"file_s3_path"`
	ProjectID        string `json:"project_id"`
	CacheProfileHash string `json:"cache_profile_hash"`
}

type CompletedEvent struct {
	TaskID         string `json:"task_id"`
	Status         string `json:"status"`
	ArtifactS3Path string `json:"artifact_s3_path,omitempty"`
	Error          string `json:"error,omitempty"`
}
