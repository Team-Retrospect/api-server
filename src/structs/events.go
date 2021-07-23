package structs

// event insertion struct
type CassandraEvent struct {
  User_id           string        `json:"user_id"`
  Session_id        string        `json:"session_id"`
  Chapter_id        string        `json:"chapter_id"`
  Body              []byte        `json:"data"`
}
