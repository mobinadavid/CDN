package requests

// successMakeBucketRequest struct for showing swagger format of success response for create bucket
type successMakeBucketRequest struct {
	IsSuccessful bool   `json:"is_successful"`
	RequestUuid  string `json:"request_uuid"`
	RequestIp    string `json:"request_ip"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         struct {
	} `json:"data,omitempty"`
}

// failedMakeBucketRequest struct for showing swagger format of failure response for create bucket
type failureMakeBucketRequest struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
}

// successRemoveBucketRequest struct for showing swagger format of success response for remove bucket
type successRemoveBucketRequest struct {
	IsSuccessful bool   `json:"is_successful"`
	RequestUuid  string `json:"request_uuid"`
	RequestIp    string `json:"request_ip"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         struct {
	} `json:"data,omitempty"`
}

// failedRemoveBucketRequest struct for showing swagger format of failure response for remove bucket
type failureRemoveBucketRequest struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
}

// failureGetBucketListRequest struct for showing swagger format of failure response for get buckets list
type failureGetBucketListRequest struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
}

// successGetBucketListRequest struct for showing swagger format of success response for get buckets list
type successGetBucketListRequest struct {
	IsSuccessful bool   `json:"is_successful"`
	RequestUuid  string `json:"request_uuid"`
	RequestIp    string `json:"request_ip"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         struct {
		Users struct {
			Limit       int `json:"limit"`
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
			TotalItems  int `json:"total_items"`
		}
	} `json:"data,omitempty"`
}

// failureGetTagRequest struct for showing swagger format of failure response for get tag data
type failureGetTagRequest struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
}

// successGetTagRequest struct for showing swagger format of success response for get tag data
type successGetTagRequest struct {
	IsSuccessful bool   `json:"is_successful"`
	RequestUuid  string `json:"request_uuid"`
	RequestIp    string `json:"request_ip"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         struct {
	} `json:"data,omitempty"`
}

// failureRemoveTagRequest struct for showing swagger format of failure response for remove tag
type failureRemoveTagRequest struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
}

// successGetTagRequest struct for showing swagger format of success response for remove tag
type successRemoveTagRequest struct {
	IsSuccessful bool   `json:"is_successful"`
	RequestUuid  string `json:"request_uuid"`
	RequestIp    string `json:"request_ip"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         struct {
	} `json:"data,omitempty"`
}
