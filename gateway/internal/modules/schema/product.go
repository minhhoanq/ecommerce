package schema

import "time"

type Product struct {
	Name        string
	Description string
	CategoryID  int32
	BrandID     int32
	Skus        []SKU
}

type SKU struct {
	Name       string
	Slug       string
	Price      Price
	Inventory  Inventory
	Attributes []Attributes
}

type Price struct {
	OriginalPrice int32
	EffectiveDate time.Time
}

type Inventory struct {
	Stock int32
}

type Attributes struct {
	AttributeID int32
	Value       string
}

func NewProduct(product Product) {

}
