package host

type Model struct {
	ID uint `gorm:"column:INVENTORY_ID;primaryKey"`

	IpAddress  string `gorm:"column:IPADDRESS"`
	HostName   string `gorm:"column:HOSTNAME"`
	OSType     string `gorm:"column:PLATFORM"`    // e.g. "Linux", "Windows", "AIX", "Solaris"
	SystemType string `gorm:"column:SYSTEM_TYPE"` // e.g. "App", "DB" "Web". Lots of fields are empty

	Status      string `gorm:"column:STATUS"`      // e.g. "Active", "Inactive", "Decommissioned"
	Environment string `gorm:"column:ENVIRONMENT"` // e.g. "PROD", "UAT", "DEV", "SIT"
	Country     string `gorm:"column:COUNTRY"`

	Lob     string `gorm:"column:LOB"`
	Appcode string `gorm:"column:APPLICATION_CODE"`
}

func (Model) TableName() string {
	return "slm_os"
}
