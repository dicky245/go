package model

type FCMRequest struct {
	Message FCMMessage `json:"message"`
}

type FCMMessageData struct {
	TokenRequest        string             `json:"token"`
	NotificationRequest *FCMNotificationData `json:"notification,omitempty"`
	DataRequest         map[string]string    `json:"data,omitempty"`
}

type FCMNotificationData struct {
	TitleRequest string `json:"title"`
	BodyRequest  string `json:"body"`
}