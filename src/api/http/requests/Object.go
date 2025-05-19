package requests

import "time"

type Owner struct {
	Owner struct {
		Space string `json:"Space"`
		Local string `json:"Local"`
	} `json:"owner"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Info struct {
	ETag              string    `json:"etag"`
	Name              string    `json:"name"`
	LastModified      time.Time `json:"lastModified"`
	Size              int64     `json:"size"`
	ContentType       string    `json:"contentType"`
	Expires           time.Time `json:"expires"`
	UserTagCount      int       `json:"UserTagCount"`
	Owner             Owner     `json:"Owner"`
	StorageClass      string    `json:"storageClass"`
	IsLatest          bool      `json:"IsLatest"`
	IsDeleteMarker    bool      `json:"IsDeleteMarker"`
	VersionID         string    `json:"VersionID"`
	ReplicationStatus string    `json:"ReplicationStatus"`
	ReplicationReady  bool      `json:"ReplicationReady"`
	Expiration        time.Time `json:"Expiration"`
	ExpirationRuleID  string    `json:"ExpirationRuleID"`
	ChecksumCRC32     string    `json:"ChecksumCRC32"`
	ChecksumCRC32C    string    `json:"ChecksumCRC32C"`
	ChecksumSHA1      string    `json:"ChecksumSHA1"`
	ChecksumSHA256    string    `json:"ChecksumSHA256"`
}
type putObjectResponse struct {
	FileName         string `json:"file_name"`
	Folder           string `json:"folder"`
	OriginalFileName string `json:"original_file_name"`
	Size             string `json:"size"`
	Url              string `json:"url"`
}

// failureGetObjectListRequest struct for showing swagger format of failure response for get objects list
type failureGetObjectListRequest struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
}

// successGetObjectListRequest struct for showing swagger format of success response for get objects list
type successGetObjectListRequest struct {
	IsSuccessful bool   `json:"is_successful"`
	RequestUuid  string `json:"request_uuid"`
	RequestIp    string `json:"request_ip"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         struct {
		Object []Info `json:"objects"`
	} `json:"data,omitempty"`
}

// successPutObjectRequest struct for showing swagger format of success response for
type successPutObjectRequest struct {
	IsSuccessful bool   `json:"is_successful"`
	RequestUuid  string `json:"request_uuid"`
	RequestIp    string `json:"request_ip"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         struct {
		Objects putObjectResponse `json:"objects"`
	} `json:"data,omitempty"`
}

// failedPutObjectRequest struct for showing swagger format of failure response for remove bucket
type failurePutObjectRequest struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
}

// failureRemoveObjectsRequest struct for showing swagger format of failure response for remove objects of bucket
type failureRemoveObjectsRequest struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
}

// successGetObjectsRequest struct for showing swagger format of success response for remove objects of bucket
type successRemoveObjectsRequest struct {
	IsSuccessful bool   `json:"is_successful"`
	RequestUuid  string `json:"request_uuid"`
	RequestIp    string `json:"request_ip"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         struct {
		Objects []string `json:"object's name"`
	} `json:"data,omitempty"`
}

// failureRemoveObjectRequest struct for showing swagger format of failure response for remove object
type failureRemoveObjectRequest struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
}

// successGetObjectsRequest struct for showing swagger format of success response for remove object
type successRemoveObjectRequest struct {
	IsSuccessful bool   `json:"is_successful"`
	RequestUuid  string `json:"request_uuid"`
	RequestIp    string `json:"request_ip"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         struct {
		Objects []string `json:"object's name"`
	} `json:"data,omitempty"`
}

// failureGetObjectRequest struct for showing swagger format of failure response for get object data
type failureGetObjectRequest struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
}

// successGetObjectRequest struct for showing swagger format of success response for get object data
type successGetObjectRequest struct {
	IsSuccessful bool   `json:"is_successful"`
	RequestUuid  string `json:"request_uuid"`
	RequestIp    string `json:"request_ip"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         struct {
		Object putObjectResponse `json:"objects"`
	} `json:"data,omitempty"`
}
