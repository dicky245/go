package model

type FCMMessage struct {
	Message MessageContent `json:"message"`
}

// Message content
type MessageContent struct {
	Token        string       `json:"token"`
	Tokens       []string     `json:"tokens"`  // harus slice string agar bisa di-loop
	Notification Notification `json:"notification"`
	Data         DataPayload  `json:"data"`
}


// Notification section
type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Data section
type DataPayload struct {
	Screen     string `json:"screen"`
	JadwalID   string `json:"jadwal_id"`
	WaktuMulai string `json:"waktu_mulai"`
}

// NotificationRequest represents the request structure for sending notifications
type NotificationRequest struct {
	Title  string            `json:"title" binding:"required"`
	Body   string            `json:"body" binding:"required"`
	Screen string            `json:"screen,omitempty"`
	Data   map[string]string `json:"data,omitempty"`
}