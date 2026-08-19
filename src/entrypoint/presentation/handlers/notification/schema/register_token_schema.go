package schema

type RegisterTokenRequest struct {
	Token   string `json:"token"`
	Segment int    `json:"segment"`
}
