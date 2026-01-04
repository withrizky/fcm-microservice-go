package model

type FcmPayload struct {
	Title   string            `json:"title" binding:"required"`
	Body    string            `json:"body" binding:"required"`
	Target  string            `json:"target" binding:"required"` // Bisa Token Device atau Nama Topic
	IsTopic bool              `json:"is_topic"`                  // Set true jika Target adalah Topic
	Data    map[string]string `json:"data"`                      // Data tambahan (opsional)
}
