package hangouts

import "time"

// ChatMessage is message type from Pub/Sub events
type ChatMessage struct {
	Type      string    `json:"type"`
	EventTime time.Time `json:"eventTime"`
	Token     string    `json:"token"`
	Message   struct {
		Name   string `json:"name"`
		Sender struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			AvatarURL   string `json:"avatarUrl"`
			Email       string `json:"email"`
			Type        string `json:"type"`
		} `json:"sender"`
		CreateTime time.Time `json:"createTime"`
		Text       string    `json:"text"`
		Thread     struct {
			Name              string `json:"name"`
			RetentionSettings struct {
				State string `json:"state"`
			} `json:"retentionSettings"`
		} `json:"thread"`
		Space struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"space"`
		ArgumentText string `json:"argumentText"`
	} `json:"message"`
	User struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		AvatarURL   string `json:"avatarUrl"`
		Email       string `json:"email"`
		Type        string `json:"type"`
	} `json:"user"`
	Space struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		DisplayName string `json:"displayName"`
	} `json:"space"`
	ConfigCompleteRedirectURL string `json:"configCompleteRedirectUrl"`
}
