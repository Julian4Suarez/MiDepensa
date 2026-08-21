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

// Category groups products by supermarket aisle, used for quick filtering and
// for grouping the generated shopping list.
type Category string

const (
	// CategoryFruitVeg covers fresh produce.
	CategoryFruitVeg Category = "FRUIT_VEG"
	// CategoryMeatFish covers meat, poultry and fish.
	CategoryMeatFish Category = "MEAT_FISH"
	// CategoryDairyEggs covers milk, cheese, butter and eggs.
	CategoryDairyEggs Category = "DAIRY_EGGS"
	// CategoryDryCanned covers non-perishable food: dry, canned and packaged.
	CategoryDryCanned Category = "DRY_CANNED"
	// CategoryDrinks covers everything drinkable.
	CategoryDrinks Category = "DRINKS"
	// CategoryHomeCare covers cleaning and personal-care supplies.
	CategoryHomeCare Category = "HOME_CARE"
)

// Categories lists every category in display order.
var Categories = []Category{
	CategoryFruitVeg,
	CategoryMeatFish,
	CategoryDairyEggs,
	CategoryDryCanned,
	CategoryDrinks,
	CategoryHomeCare,
}

// IsValid reports whether the category is one of the known values.
func (c Category) IsValid() bool {
	switch c {
	case CategoryFruitVeg, CategoryMeatFish, CategoryDairyEggs,
		CategoryDryCanned, CategoryDrinks, CategoryHomeCare:
		return true
	default:
		return false
	}
}
