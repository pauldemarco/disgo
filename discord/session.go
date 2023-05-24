package discord

type Session struct {
	SessionID        string `json:"sessionId"`
	ResumeGatewayURL string `json:"resumeGatewayUrl"`
	SequenceID       int    `json:"sequenceId"`
}
