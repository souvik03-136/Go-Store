package utils

type BaseResponse struct {
	Success    bool        `json:"success,omitempty"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	MetaData   interface{} `json:"metadata,omitempty"`
	Error      *Error      `json:"error,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
