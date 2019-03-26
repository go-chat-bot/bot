package api

type Channel struct {
	Id           string `json:"_id"`
	Name         string `json:"name"`
	MessageCount int `json:"msgs"`
	UserNames    []string `json:"usernames"`

	User         User `json:"u"`

	ReadOnly     bool `json:"ro"`
	Timestamp    string `json:"ts"`
	T            string `json:"t"`
	UpdatedAt    string `json:"_updatedAt"`
	SysMes       bool `json:"sysMes"`
}

type User struct {
	Id       string `json:"_id"`
	UserName string `json:"username"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"pass"`
}

type Message struct {
	Id        string `json:"_id"`
	ChannelId string `json:"rid"`
	Text      string `json:"msg"`
	Timestamp string `json:"ts"`
	User      User `json:"u"`
}

type Info struct {
	Version string `json:"version"`

	Build struct {
		Date string `json:"date"`
		NodeVersion string `json:"nodeVersion"`
		Arch string `json:"arch"`
		Platform string `json:"platform"`
		OsRelease string `json:"osRelease"`
		TotalMemory int64 `json:"totalMemory"`
		FreeMemory int64 `json:"freeMemory"`
		CpuCount int `json:"cpus"`
	} `json:"build"`

	Travis struct {
		BuildNumber string `json:"buildNumber"`
		Branch string `json:"branch"`
		Tag string `json:"tag"`
	} `json:"travis"`

	Commit  struct {
		Hash string `json:"hash"`
		Date string `json:"date"`
		Author string `json:"author"`
		Subject string `json:"subject"`
		Tag string `json:"tag"`
		Branch string `json:"branch"`
	} `json:"commit"`

	GraphicsMagick struct {
		Enabled bool `json:"enabled"`
	} `json:"GraphicsMagick"`

	ImageMagick struct {
		Enabled bool `json:"enabled"`
		Version string `json:"version"`
	} `json:"ImageMagick"`
}