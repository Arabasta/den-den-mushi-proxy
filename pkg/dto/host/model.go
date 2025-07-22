package host

type Model struct {
	ID          uint   `gorm:"column:Inventory_ID;primaryKey"`
	IpAddress   string `gorm:"column:IpAddress"`
	HostName    string `gorm:"column:HostName"`
	OSType      string `gorm:"column:OS_TYPE"`
	Status      string `gorm:"column:Status"`
	Environment string `gorm:"column:Environment"`
	Country     string `gorm:"column:Country"`
}

func (Model) TableName() string {
	return "slm_os"
}
