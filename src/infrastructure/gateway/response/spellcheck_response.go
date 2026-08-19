package response

type SpellCheckResponse struct {
	Success      bool   `json:"success"`
	FileContents string `json:"file_contents,omitempty"`
	Error        string `json:"error_message,omitempty"`
}
