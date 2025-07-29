package puppet_trusted

type Model struct {
	Certname string `gorm:"CERTNAME"`
}

var tableName = "puppet_trusted"

func (Model) TableName() string {
	return tableName
}

func SetTableName(name string) {
	tableName = name
}
