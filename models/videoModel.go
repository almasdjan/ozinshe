package models

type Video struct {
	Material_id uint
	Sezon       uint
	Series      uint
	Image_src   string `gorm:"unique"`
	Video_src   string `gorm:"unique"`
}

type Series struct {
	Series    uint
	Image_src string
	Video_src string
}

type SezonsAndSeries struct {
	SezonCount  uint
	SeriesCount uint
}
