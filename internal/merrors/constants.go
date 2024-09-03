package merrors

var errorType = struct {
	validation         string
	server             string
	Unauthorized       string
	conflict           string
	ServiceUnavailable string
	Forbidden          string
	Downstream         string
	NotFound           string
}{
	validation:         "validation",
	server:             "server",
	Unauthorized:       "unauthorized",
	conflict:           "conflict",
	ServiceUnavailable: "service unavailable",
	Forbidden:          "forbidden",
	Downstream:         "downstream",
	NotFound:           "not found",
}
