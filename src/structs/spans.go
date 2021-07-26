package structs

// zipkin data format
type SpanStructInput struct {
  Trace_id          string              `json:"traceId"`
  Span_id           string              `json:"id"`
  Time_sent         int                 `json:"timestamp"`
  Duration          int                 `json:"duration"`

  Tags              map[string]string   `json:"tags"`
}

// span insertion struct
type CassandraSpan struct {
  // top-level span data
  Trace_id          string        `json:"trace_id"`
  Span_id           string        `json:"span_id"`
  Time_sent         int           `json:"time_sent"`
  Duration          string        `json:"time_duration"`
  Data              []byte        `json:"data"`

  // derived from tags
  Trigger_route     string        `json:"trigger_route"`
  User_id           string        `json:"user_id"`
  Session_id        string        `json:"session_id"`
  Chapter_id        string        `json:"chapter_id"`
  Status_code       int16         `json:"status_code"`
  Request_data      []byte        `json:"request_data"`

  Is_db             bool
}
