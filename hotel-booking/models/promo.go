package models

type Promo struct {
	Code        string
	DiscountPct int
	MaxDiscount int
	MinSpend    int
	Description string
	Active      bool
}
