package entities

// StockStatus is how much of a product is left in a pantry.
type StockStatus string

const (
	// StatusOut means the product must be bought (rendered red).
	StatusOut StockStatus = "OUT"
	// StatusLow means the product is running low (rendered amber).
	StatusLow StockStatus = "LOW"
	// StatusOK means there is enough of the product (rendered green).
	StatusOK StockStatus = "OK"
)

// StockStatuses lists every stock status in display order.
var StockStatuses = []StockStatus{StatusOut, StatusLow, StatusOK}

// IsValid reports whether the status is one of the known values.
func (s StockStatus) IsValid() bool {
	switch s {
	case StatusOut, StatusLow, StatusOK:
		return true
	default:
		return false
	}
}

// PantryView is the section of the pantry a product is filed under.
type PantryView string

const (
	// ViewPrimary holds the products checked on every shopping trip.
	ViewPrimary PantryView = "PRIMARY"
	// ViewSecondary holds the products checked occasionally.
	ViewSecondary PantryView = "SECONDARY"
	// ViewOther holds everything else.
	ViewOther PantryView = "OTHER"
)

// PantryViews lists every view in display order.
var PantryViews = []PantryView{ViewPrimary, ViewSecondary, ViewOther}

// IsValid reports whether the view is one of the known values.
func (v PantryView) IsValid() bool {
	switch v {
	case ViewPrimary, ViewSecondary, ViewOther:
		return true
	default:
		return false
	}
}

// Category groups products by what they are, used for quick filtering.
type Category string

const (
	// CategoryFresh covers produce, dairy, meat, fish and bread.
	CategoryFresh Category = "FRESH"
	// CategoryPantry covers dry and canned goods.
	CategoryPantry Category = "PANTRY"
	// CategoryDrinks covers everything drinkable.
	CategoryDrinks Category = "DRINKS"
	// CategoryHomeCare covers cleaning and personal-care supplies.
	CategoryHomeCare Category = "HOME_CARE"
)

// Categories lists every category in display order.
var Categories = []Category{CategoryFresh, CategoryPantry, CategoryDrinks, CategoryHomeCare}

// IsValid reports whether the category is one of the known values.
func (c Category) IsValid() bool {
	switch c {
	case CategoryFresh, CategoryPantry, CategoryDrinks, CategoryHomeCare:
		return true
	default:
		return false
	}
}
