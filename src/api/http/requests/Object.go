package requests

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
		Users struct {
			Limit       int `json:"limit"`
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
			TotalItems  int `json:"total_items"`
		}
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
	} `json:"data,omitempty"`
}
