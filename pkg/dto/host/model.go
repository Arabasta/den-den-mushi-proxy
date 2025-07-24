package host

type Model struct {
	ID          uint   `gorm:"column:INVENTORY_ID;primaryKey"`
	IpAddress   string `gorm:"column:IPADDRESS"`
	HostName    string `gorm:"column:HOSTNAME"`
	Appcode     string `gorm:"column:APPLICATION_CODE"`
	OSType      string `gorm:"column:OS_TYPE"`
	Status      string `gorm:"column:STATUS"`
	Environment string `gorm:"column:ENVIRONMENT"`
	Country     string `gorm:"column:COUNTRY"`
}

func (Model) TableName() string {
	return "slm_os"
}
