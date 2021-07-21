package structs

// snapshot struct
type Snapshot struct {
  Session_id        string        `json:"session_id"`
  Body              []byte        `json:"data"`
}
