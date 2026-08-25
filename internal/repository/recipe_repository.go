package repository

import (
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) domain.RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) Create(recipe *domain.Recipe) error {
	return r.db.Create(recipe).Error
}

func (r *recipeRepository) GetByID(id int) (*domain.Recipe, error) {
	var recipe domain.Recipe
	if err := r.db.First(&recipe, id).Error; err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *recipeRepository) ListByUser(userID int) ([]domain.Recipe, error) {
	var recipes []domain.Recipe
	if err := r.db.Where("created_by = ?", userID).Order("created_at DESC").Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func (r *recipeRepository) ListPublic() ([]domain.Recipe, error) {
	var recipes []domain.Recipe
	if err := r.db.Where("is_public = ?", true).Order("created_at DESC").Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func (r *recipeRepository) Search(query string, dietType, cuisineType string) ([]domain.Recipe, error) {
	var recipes []domain.Recipe
	db := r.db.Where("is_public = ?", true)

	if query != "" {
		db = db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if dietType != "" {
		db = db.Where("diet_type = ?", dietType)
	}
	if cuisineType != "" {
		db = db.Where("cuisine_type = ?", cuisineType)
	}

	if err := db.Order("created_at DESC").Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func (r *recipeRepository) Update(recipe *domain.Recipe) error {
	return r.db.Save(recipe).Error
}

func (r *recipeRepository) Delete(id int) error {
	return r.db.Delete(&domain.Recipe{}, id).Error
}

func (r *recipeRepository) AddIngredient(ingredient *domain.RecipeIngredient) error {
	return r.db.Create(ingredient).Error
}

func (r *recipeRepository) GetIngredients(recipeID int) ([]domain.RecipeIngredient, error) {
	var ingredients []domain.RecipeIngredient
	if err := r.db.Where("recipe_id = ?", recipeID).Find(&ingredients).Error; err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (r *recipeRepository) DeleteIngredient(id int) error {
	return r.db.Delete(&domain.RecipeIngredient{}, id).Error
}

type mealPlanRepository struct {
	db *gorm.DB
}

func NewMealPlanRepository(db *gorm.DB) domain.MealPlanRepository {
	return &mealPlanRepository{db: db}
}

func (r *mealPlanRepository) Create(mealPlan *domain.MealPlan) error {
	return r.db.Create(mealPlan).Error
}

func (r *mealPlanRepository) GetByID(id int) (*domain.MealPlan, error) {
	var mealPlan domain.MealPlan
	if err := r.db.First(&mealPlan, id).Error; err != nil {
		return nil, err
	}
	return &mealPlan, nil
}

func (r *mealPlanRepository) GetByUser(userID int) ([]domain.MealPlan, error) {
	var mealPlans []domain.MealPlan
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&mealPlans).Error; err != nil {
		return nil, err
	}
	return mealPlans, nil
}

func (r *mealPlanRepository) GetActive(userID int) (*domain.MealPlan, error) {
	var mealPlan domain.MealPlan
	if err := r.db.Where("user_id = ? AND is_active = ?", userID, true).First(&mealPlan).Error; err != nil {
		return nil, err
	}
	return &mealPlan, nil
}

func (r *mealPlanRepository) Update(mealPlan *domain.MealPlan) error {
	return r.db.Save(mealPlan).Error
}

func (r *mealPlanRepository) Delete(id int) error {
	return r.db.Delete(&domain.MealPlan{}, id).Error
}

func (r *mealPlanRepository) AddEntry(entry *domain.MealPlanEntry) error {
	return r.db.Create(entry).Error
}

func (r *mealPlanRepository) GetEntries(mealPlanID int) ([]domain.MealPlanEntry, error) {
	var entries []domain.MealPlanEntry
	if err := r.db.Where("meal_plan_id = ?", mealPlanID).Order("meal_date ASC, meal_type").Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *mealPlanRepository) GetEntriesByDate(userID int, date time.Time) ([]domain.MealPlanEntry, error) {
	var entries []domain.MealPlanEntry
	dateStart := date.Truncate(24 * time.Hour)
	dateEnd := dateStart.Add(24 * time.Hour)
	
	if err := r.db.Joins("JOIN meal_plans ON meal_plan_entries.meal_plan_id = meal_plans.id").
		Where("meal_plans.user_id = ? AND meal_plan_entries.meal_date >= ? AND meal_plan_entries.meal_date < ?", userID, dateStart, dateEnd).
		Order("meal_type").
		Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *mealPlanRepository) UpdateEntry(entry *domain.MealPlanEntry) error {
	return r.db.Save(entry).Error
}

func (r *mealPlanRepository) DeleteEntry(id int) error {
	return r.db.Delete(&domain.MealPlanEntry{}, id).Error
}

func (r *mealPlanRepository) MarkEntryCompleted(id int) error {
	return r.db.Model(&domain.MealPlanEntry{}).Where("id = ?", id).Update("is_completed", true).Error
}

type shoppingListRepository struct {
	db *gorm.DB
}

func NewShoppingListRepository(db *gorm.DB) domain.ShoppingListRepository {
	return &shoppingListRepository{db: db}
}

func (r *shoppingListRepository) Create(shoppingList *domain.ShoppingList) error {
	return r.db.Create(shoppingList).Error
}

func (r *shoppingListRepository) GetByID(id int) (*domain.ShoppingList, error) {
	var shoppingList domain.ShoppingList
	if err := r.db.First(&shoppingList, id).Error; err != nil {
		return nil, err
	}
	return &shoppingList, nil
}

func (r *shoppingListRepository) GetByUser(userID int) ([]domain.ShoppingList, error) {
	var shoppingLists []domain.ShoppingList
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&shoppingLists).Error; err != nil {
		return nil, err
	}
	return shoppingLists, nil
}

func (r *shoppingListRepository) GetActive(userID int) (*domain.ShoppingList, error) {
	var shoppingList domain.ShoppingList
	if err := r.db.Where("user_id = ? AND is_active = ?", userID, true).First(&shoppingList).Error; err != nil {
		return nil, err
	}
	return &shoppingList, nil
}

func (r *shoppingListRepository) Update(shoppingList *domain.ShoppingList) error {
	return r.db.Save(shoppingList).Error
}

func (r *shoppingListRepository) Delete(id int) error {
	return r.db.Delete(&domain.ShoppingList{}, id).Error
}

func (r *shoppingListRepository) AddItem(item *domain.ShoppingListItem) error {
	return r.db.Create(item).Error
}

func (r *shoppingListRepository) GetItems(shoppingListID int) ([]domain.ShoppingListItem, error) {
	var items []domain.ShoppingListItem
	if err := r.db.Where("shopping_list_id = ?", shoppingListID).Order("ingredient_name").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *shoppingListRepository) UpdateItem(item *domain.ShoppingListItem) error {
	return r.db.Save(item).Error
}

func (r *shoppingListRepository) DeleteItem(id int) error {
	return r.db.Delete(&domain.ShoppingListItem{}, id).Error
}

func (r *shoppingListRepository) MarkItemPurchased(id int) error {
	now := time.Now()
	return r.db.Model(&domain.ShoppingListItem{}).Where("id = ?", id).Updates(map[string]interface{}{"is_purchased": true, "purchased_at": now}).Error
}