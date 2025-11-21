package models

type UploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Key      string `json:"key"`
		Hash     string `json:"hash"`
		URL      string `json:"url"`
		FileSize int64  `json:"file_size"`
		MimeType string `json:"mime_type"`
	} `json:"data,omitempty"`
}

type ImageInfo struct {
	ID       string `json:"id"`
	Key      string `json:"key"`
	URL      string `json:"url"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
	Uploaded string `json:"uploaded"`
}

type ImageListResponse struct {
	Success bool        `json:"success"`
	Data    []ImageInfo `json:"data"`
	Total   int         `json:"total"`
}