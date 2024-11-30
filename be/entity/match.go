package entity

type Match struct {
	ID        int    `gorm:"primaryKey" json:"id,omitempty"`
	ProfileID int    `json:"profile_id,omitempty"`
	PartnerID int    `json:"partner_id,omitempty"`
	Status    string `json:"status,omitempty"`

	Profile Profile `gorm:"foreignKey:ProfileID;references:ID" json:"profile,omitempty"`
	Partner Profile `gorm:"foreignKey:PartnerID;references:ID" json:"-"`
}

const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
)
