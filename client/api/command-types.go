package api

type CallIssueCommandResponse struct {
	StatusCode int
	Body       string
}

type ResponseBody struct {
	Message string
	Success bool
}
