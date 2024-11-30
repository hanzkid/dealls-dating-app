package entity

type Profile struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
}

type ProfileViewLog struct {
	ID        uint `json:"id"`
	ViewerID  uint `json:"viewer_id"`
	ProfileID uint `json:"profile_id"`
}
