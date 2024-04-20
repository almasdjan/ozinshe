package models

type Category struct {
	ID           uint
	CategoryName string
}

type Material_category struct {
	Material_id uint
	Category_id uint
}
