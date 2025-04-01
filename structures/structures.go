package structures

type JSONMessage struct {
	Status       int    `json:"status"`
	Message      string `json:"message,omitempty"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:",omitempty"`
}

type User struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

type OptsFile struct {
	IsAvatar bool `json:"IsAvatar,omitempty"`
	IsImage  bool `json:"IsImage,omitempty"`
}
