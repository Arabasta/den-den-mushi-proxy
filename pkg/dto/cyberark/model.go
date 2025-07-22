package cyberark

type Model struct {
	ID uint `gorm:"primaryKey;column:ID;type:bigint"` // todo: tmp for dev

	Object   string `gorm:"column:Objectname"`
	Hostname string `gorm:"column:Hostname"`
	Ip       string `gorm:"column:IP"`
}

func (Model) TableName() string {
	return "cyberarkobject_map"
}
