package vehiclecommand

type CallIssueCommandResponse struct {
	StatusCode int
	Body       string
}

type ResponseBody struct {
	Message string
	Success bool
}

type CommandRequest struct {
	AuthToken string `json:"authToken"`
	Vin       string `json:"vin"`
	Command   string `json:"command"`
}

type CommandResponse struct {
	Success bool   `json:"status"`
	Message string `json:"message"`
}
