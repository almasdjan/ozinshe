package models

type Genre struct {
	ID        uint
	Image     string
	GenreName string
}

type Material_genre struct {
	Material_id uint
	Genre_id    uint
}

type Genrejson struct {
	Genre string
}
