package models

type Video struct {
	Id          uint
	Material_id uint
	Sezon       uint
	Series      uint
	Video_src   string `gorm:"unique"`
	Viewed      int
}

type Videos struct {
	Sezon     uint
	Series    uint
	Video_src string `gorm:"unique"`
}

type Series struct {
	Id        uint
	Series    uint
	Image_src string
	Video_src string
}

type SezonsAndSeries struct {
	SezonCount  uint
	SeriesCount uint
}
