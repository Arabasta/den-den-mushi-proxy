package puppet_trusted

type Model struct {
	Certname string `gorm:"CERTNAME"`
}

func (Model) TableName() string {
	return "PUPPET_TRUSTED"
}
