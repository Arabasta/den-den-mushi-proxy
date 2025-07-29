package cyberark

type Model struct {
	ID uint `gorm:"primaryKey;column:ID;type:bigint"` // todo: tmp for dev

	Object   string `gorm:"column:OBJECTNAME"`
	Hostname string `gorm:"column:HOSTNAME"`
	Ip       string `gorm:"column:IP"`
}

var tableName = "cyberarkobject_map"

func (Model) TableName() string {
	return tableName
}

func SetTableName(name string) {
	tableName = name
}
