package jti

type Jti struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Subject   string `gorm:"column:subject" json:"subject"`
	ExpiresAt int64  `gorm:"column:expires_at" json:"expires_at"`
}
