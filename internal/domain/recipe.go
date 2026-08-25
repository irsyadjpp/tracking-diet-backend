package domain

import "time"

// Recipe represents a food recipe
type Recipe struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Instructions string   `gorm:"type:text" json:"instructions"`
	PrepTime    int       `gorm:"comment:in minutes" json:"prep_time"`
	CookTime    int       `gorm:"comment:in minutes" json:"cook_time"`
	Servings    int       `gorm:"not null" json:"servings"`
	CaloriesPerServing float64 `json:"calories_per_serving"`
	ProteinPerServing float64 `json:"protein_per_serving"`
	CarbsPerServing  float64 `json:"carbs_per_serving"`
	FatPerServing    float64 `json:"fat_per_serving"`
	FiberPerServing  float64 `json:"fiber_per_serving"`
	CuisineType  string    `gorm:"size:50" json:"cuisine_type"`
	DietType     string    `gorm:"size:50" json:"diet_type"` // vegetarian, vegan, keto, paleo, etc.
	Difficulty   string    `gorm:"size:20;default:medium" json:"difficulty"` // easy, medium, hard
	ImageURL     string    `gorm:"size:500" json:"image_url"`
	IsPublic     bool      `gorm:"default:true" json:"is_public"`
	CreatedBy    int       `gorm:"index" json:"created_by"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Recipe model
func (Recipe) TableName() string {
	return "recipes"
}

// RecipeIngredient represents ingredients for a recipe
type RecipeIngredient struct {
	ID          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipeID    int     `gorm:"not null;index" json:"recipe_id"`
	Name        string  `gorm:"size:200;not null" json:"name"`
	Quantity    float64 `gorm:"not null" json:"quantity"`
	Unit        string  `gorm:"size:50" json:"unit"`
}

// TableName specifies the table name for RecipeIngredient model
func (RecipeIngredient) TableName() string {
	return "recipe_ingredients"
}

// MealPlan represents a user's meal plan
type MealPlan struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"not null;index" json:"user_id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	EndDate     time.Time `gorm:"not null" json:"end_date"`
	TargetCaloriesPerDay *int `json:"target_calories_per_day"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for MealPlan model
func (MealPlan) TableName() string {
	return "meal_plans"
}

// MealPlanEntry represents a single meal entry in a meal plan
type MealPlanEntry struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	MealPlanID  int       `gorm:"not null;index" json:"meal_plan_id"`
	RecipeID    *int      `gorm:"index" json:"recipe_id"`
	MealType    string    `gorm:"size:50;not null" json:"meal_type"` // breakfast, lunch, dinner, snack
	MealDate    time.Time `gorm:"not null;index" json:"meal_date"`
	Calories    *int      `json:"calories"`
	Notes       string    `gorm:"type:text" json:"notes"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
}

// TableName specifies the table name for MealPlanEntry model
func (MealPlanEntry) TableName() string {
	return "meal_plan_entries"
}

// ShoppingList represents a shopping list for meal planning
type ShoppingList struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"not null;index" json:"user_id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	MealPlanID  *int      `gorm:"index" json:"meal_plan_id"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for ShoppingList model
func (ShoppingList) TableName() string {
	return "shopping_lists"
}

// ShoppingListItem represents an item in a shopping list
type ShoppingListItem struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ShoppingListID int       `gorm:"not null;index" json:"shopping_list_id"`
	IngredientName string   `gorm:"size:200;not null" json:"ingredient_name"`
	Quantity       float64   `gorm:"not null" json:"quantity"`
	Unit           string    `gorm:"size:50" json:"unit"`
	IsPurchased    bool      `gorm:"default:false" json:"is_purchased"`
	PurchasedAt    *time.Time `json:"purchased_at"`
}

// TableName specifies the table name for ShoppingListItem model
func (ShoppingListItem) TableName() string {
	return "shopping_lists_items"
}

// RecipeRepository defines operations for managing recipes
type RecipeRepository interface {
	Create(recipe *Recipe) error
	GetByID(id int) (*Recipe, error)
	ListByUser(userID int) ([]Recipe, error)
	ListPublic() ([]Recipe, error)
	Search(query string, dietType, cuisineType string) ([]Recipe, error)
	Update(recipe *Recipe) error
	Delete(id int) error
	AddIngredient(ingredient *RecipeIngredient) error
	GetIngredients(recipeID int) ([]RecipeIngredient, error)
	DeleteIngredient(id int) error
}

// MealPlanRepository defines operations for managing meal plans
type MealPlanRepository interface {
	Create(mealPlan *MealPlan) error
	GetByID(id int) (*MealPlan, error)
	GetByUser(userID int) ([]MealPlan, error)
	GetActive(userID int) (*MealPlan, error)
	Update(mealPlan *MealPlan) error
	Delete(id int) error
	AddEntry(entry *MealPlanEntry) error
	GetEntries(mealPlanID int) ([]MealPlanEntry, error)
	GetEntriesByDate(userID int, date time.Time) ([]MealPlanEntry, error)
	UpdateEntry(entry *MealPlanEntry) error
	DeleteEntry(id int) error
	MarkEntryCompleted(id int) error
}

// ShoppingListRepository defines operations for managing shopping lists
type ShoppingListRepository interface {
	Create(shoppingList *ShoppingList) error
	GetByID(id int) (*ShoppingList, error)
	GetByUser(userID int) ([]ShoppingList, error)
	GetActive(userID int) (*ShoppingList, error)
	Update(shoppingList *ShoppingList) error
	Delete(id int) error
	AddItem(item *ShoppingListItem) error
	GetItems(shoppingListID int) ([]ShoppingListItem, error)
	UpdateItem(item *ShoppingListItem) error
	DeleteItem(id int) error
	MarkItemPurchased(id int) error
}