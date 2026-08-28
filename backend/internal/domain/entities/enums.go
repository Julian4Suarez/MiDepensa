package entities

// ItemStatus is the product's state in the shopping workflow.
type ItemStatus string

const (
	// StatusDiscarded means the product is not needed on this shopping trip.
	StatusDiscarded ItemStatus = "DISCARDED"
	// StatusPending means the product still needs a decision.
	StatusPending ItemStatus = "PENDING"
	// StatusInCart means the product belongs on the shopping list.
	StatusInCart ItemStatus = "IN_CART"
	// StatusArchived removes the product from the active pantry views.
	StatusArchived ItemStatus = "ARCHIVED"
)

// ItemStatuses lists every item status in display order.
var ItemStatuses = []ItemStatus{StatusDiscarded, StatusPending, StatusInCart, StatusArchived}

// IsValid reports whether the status is one of the known values.
func (s ItemStatus) IsValid() bool {
	switch s {
	case StatusDiscarded, StatusPending, StatusInCart, StatusArchived:
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
	// CategoryFruit contains fresh and preserved fruit sold as fruit.
	CategoryFruit Category = "FRUIT"
	// CategoryVegetables contains vegetables regardless of preservation method.
	CategoryVegetables Category = "VEGETABLES"
	// CategoryMeatFish covers meat, poultry and fish.
	CategoryMeatFish Category = "MEAT_FISH"
	// CategoryDairyEggs covers milk, cheese, butter and eggs.
	CategoryDairyEggs             Category = "DAIRY_EGGS"
	CategoryBakeryBreakfastSnacks Category = "BAKERY_BREAKFAST_SNACKS"
	CategoryPastaRiceLegumes      Category = "PASTA_RICE_LEGUMES"
	CategoryCookingCondiments     Category = "COOKING_CONDIMENTS"
	// CategoryDrinks covers everything drinkable.
	CategoryDrinks            Category = "DRINKS"
	CategoryHouseholdCleaning Category = "HOUSEHOLD_CLEANING"
	CategoryPersonalCare      Category = "PERSONAL_CARE"
)

// Categories lists every category in display order.
var Categories = []Category{
	CategoryFruit,
	CategoryVegetables,
	CategoryMeatFish,
	CategoryDairyEggs,
	CategoryBakeryBreakfastSnacks,
	CategoryPastaRiceLegumes,
	CategoryCookingCondiments,
	CategoryDrinks,
	CategoryHouseholdCleaning,
	CategoryPersonalCare,
}

// IsValid reports whether the category is one of the known values.
func (c Category) IsValid() bool {
	switch c {
	case CategoryFruit, CategoryVegetables, CategoryMeatFish, CategoryDairyEggs,
		CategoryBakeryBreakfastSnacks, CategoryPastaRiceLegumes,
		CategoryCookingCondiments, CategoryDrinks, CategoryHouseholdCleaning,
		CategoryPersonalCare:
		return true
	default:
		return false
	}
}
