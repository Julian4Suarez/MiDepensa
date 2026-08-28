package entities

// StockStatus is the current stock state, or archived when removed from active views.
type StockStatus string

const (
	// StatusOut means the product must be bought (rendered red).
	StatusOut StockStatus = "OUT"
	// StatusLow means the product is running low (rendered amber).
	StatusLow StockStatus = "LOW"
	// StatusOK means there is enough of the product (rendered green).
	StatusOK StockStatus = "OK"
	// StatusArchived removes the product from the active pantry views.
	StatusArchived StockStatus = "ARCHIVED"
)

// StockStatuses lists every stock status in display order.
var StockStatuses = []StockStatus{StatusOut, StatusLow, StatusOK, StatusArchived}

// IsValid reports whether the status is one of the known values.
func (s StockStatus) IsValid() bool {
	switch s {
	case StatusOut, StatusLow, StatusOK, StatusArchived:
		return true
	default:
		return false
	}
}

// ProductType is how often a product is bought. It is the coarse filter of the
// pantry screen, alongside Category.
type ProductType string

const (
	// TypeEssential is bought on every shopping trip.
	TypeEssential ProductType = "ESSENTIAL"
	// TypeSecondary is bought from time to time.
	TypeSecondary ProductType = "SECONDARY"
)

// ProductTypes lists every type in display order.
var ProductTypes = []ProductType{TypeEssential, TypeSecondary}

// IsValid reports whether the type is one of the known values.
func (t ProductType) IsValid() bool {
	switch t {
	case TypeEssential, TypeSecondary:
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
