package models

type Genre struct {
	ID        uint
	GenreName string
}

type Material_genre struct {
	Material_id uint
	Genre_id    uint
}
