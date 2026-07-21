package tanswer

type SuccessEnvelope struct {
	Success  bool   `json:"success"`
	Task     string `json:"task"`
	Command  string `json:"command"`
	Query    any    `json:"query,omitempty"`
	Data     any    `json:"data"`
	Warnings any    `json:"warnings,omitempty"`
}

type ErrorEnvelope struct {
	Success bool         `json:"success"`
	Task    string       `json:"task"`
	Command string       `json:"command"`
	Error   CommandError `json:"error"`
}

type CommandError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
}

func NewSuccessEnvelope(task, command string, query, data, warnings any) SuccessEnvelope {
	return SuccessEnvelope{
		Success:  true,
		Task:     task,
		Command:  command,
		Query:    query,
		Data:     data,
		Warnings: warnings,
	}
}

func NewErrorEnvelope(task, command, code, message string, retryable bool) ErrorEnvelope {
	return ErrorEnvelope{
		Success: false,
		Task:    task,
		Command: command,
		Error: CommandError{
			Code:      code,
			Message:   message,
			Retryable: retryable,
		},
	}
}
