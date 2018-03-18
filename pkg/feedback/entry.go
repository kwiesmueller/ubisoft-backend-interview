package feedback

// Entry definition
type Entry struct {
	ID        string `json:"id"`
	SessionID string `json:"sessionID"`
	UserID    string `json:"userID"`
	Rating    int8   `json:"rating"`
	Comment   string `json:"comment"`
}
