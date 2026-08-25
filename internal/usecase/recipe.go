package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrRecipeNotFound    = errors.New("recipe not found")
	ErrMealPlanNotFound  = errors.New("meal plan not found")
	ErrShoppingListNotFound = errors.New("shopping list not found")
)

// RecipeUsecase handles recipe business logic
type RecipeUsecase struct {
	recipeRepo domain.RecipeRepository
}

// NewRecipeUsecase creates a new recipe use case
func NewRecipeUsecase(recipeRepo domain.RecipeRepository) *RecipeUsecase {
	return &RecipeUsecase{
		recipeRepo: recipeRepo,
	}
}

// CreateRecipe creates a new recipe
func (uc *RecipeUsecase) CreateRecipe(createdBy int, name, description, instructions string, prepTime, cookTime, servings int, caloriesPerServing, proteinPerServing, carbsPerServing, fatPerServing, fiberPerServing float64, cuisineType, dietType, difficulty, imageURL string, isPublic bool) (*domain.Recipe, error) {
	recipe := &domain.Recipe{
		Name:                 name,
		Description:          description,
		Instructions:         instructions,
		PrepTime:             prepTime,
		CookTime:             cookTime,
		Servings:             servings,
		CaloriesPerServing:   caloriesPerServing,
		ProteinPerServing:    proteinPerServing,
		CarbsPerServing:     carbsPerServing,
		FatPerServing:       fatPerServing,
		FiberPerServing:     fiberPerServing,
		CuisineType:          cuisineType,
		DietType:             dietType,
		Difficulty:           difficulty,
		ImageURL:             imageURL,
		IsPublic:             isPublic,
		CreatedBy:            createdBy,
	}

	if err := uc.recipeRepo.Create(recipe); err != nil {
		return nil, fmt.Errorf("failed to create recipe: %w", err)
	}

	return recipe, nil
}

// GetRecipe retrieves a recipe by ID
func (uc *RecipeUsecase) GetRecipe(id int) (*domain.Recipe, error) {
	recipe, err := uc.recipeRepo.GetByID(id)
	if err != nil {
		return nil, ErrRecipeNotFound
	}
	return recipe, nil
}

// GetUserRecipes retrieves all recipes for a user
func (uc *RecipeUsecase) GetUserRecipes(userID int) ([]domain.Recipe, error) {
	recipes, err := uc.recipeRepo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user recipes: %w", err)
	}
	return recipes, nil
}

// GetPublicRecipes retrieves all public recipes
func (uc *RecipeUsecase) GetPublicRecipes() ([]domain.Recipe, error) {
	recipes, err := uc.recipeRepo.ListPublic()
	if err != nil {
		return nil, fmt.Errorf("failed to get public recipes: %w", err)
	}
	return recipes, nil
}

// SearchRecipes searches for recipes based on criteria
func (uc *RecipeUsecase) SearchRecipes(query, dietType, cuisineType string) ([]domain.Recipe, error) {
	recipes, err := uc.recipeRepo.Search(query, dietType, cuisineType)
	if err != nil {
		return nil, fmt.Errorf("failed to search recipes: %w", err)
	}
	return recipes, nil
}

// UpdateRecipe updates an existing recipe
func (uc *RecipeUsecase) UpdateRecipe(id int, name, description, instructions *string, prepTime, cookTime, servings *int, caloriesPerServing, proteinPerServing, carbsPerServing, fatPerServing, fiberPerServing *float64, cuisineType, dietType, difficulty, imageURL *string, isPublic *bool) (*domain.Recipe, error) {
	recipe, err := uc.recipeRepo.GetByID(id)
	if err != nil {
		return nil, ErrRecipeNotFound
	}

	if name != nil {
		recipe.Name = *name
	}
	if description != nil {
		recipe.Description = *description
	}
	if instructions != nil {
		recipe.Instructions = *instructions
	}
	if prepTime != nil {
		recipe.PrepTime = *prepTime
	}
	if cookTime != nil {
		recipe.CookTime = *cookTime
	}
	if servings != nil {
		recipe.Servings = *servings
	}
	if caloriesPerServing != nil {
		recipe.CaloriesPerServing = *caloriesPerServing
	}
	if proteinPerServing != nil {
		recipe.ProteinPerServing = *proteinPerServing
	}
	if carbsPerServing != nil {
		recipe.CarbsPerServing = *carbsPerServing
	}
	if fatPerServing != nil {
		recipe.FatPerServing = *fatPerServing
	}
	if fiberPerServing != nil {
		recipe.FiberPerServing = *fiberPerServing
	}
	if cuisineType != nil {
		recipe.CuisineType = *cuisineType
	}
	if dietType != nil {
		recipe.DietType = *dietType
	}
	if difficulty != nil {
		recipe.Difficulty = *difficulty
	}
	if imageURL != nil {
		recipe.ImageURL = *imageURL
	}
	if isPublic != nil {
		recipe.IsPublic = *isPublic
	}

	if err := uc.recipeRepo.Update(recipe); err != nil {
		return nil, fmt.Errorf("failed to update recipe: %w", err)
	}

	return recipe, nil
}

// DeleteRecipe deletes a recipe
func (uc *RecipeUsecase) DeleteRecipe(id int) error {
	if err := uc.recipeRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}
	return nil
}

// AddIngredient adds an ingredient to a recipe
func (uc *RecipeUsecase) AddIngredient(recipeID int, name string, quantity float64, unit string) (*domain.RecipeIngredient, error) {
	ingredient := &domain.RecipeIngredient{
		RecipeID: recipeID,
		Name:     name,
		Quantity: quantity,
		Unit:     unit,
	}

	if err := uc.recipeRepo.AddIngredient(ingredient); err != nil {
		return nil, fmt.Errorf("failed to add ingredient: %w", err)
	}

	return ingredient, nil
}

// GetRecipeIngredients retrieves ingredients for a recipe
func (uc *RecipeUsecase) GetRecipeIngredients(recipeID int) ([]domain.RecipeIngredient, error) {
	ingredients, err := uc.recipeRepo.GetIngredients(recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe ingredients: %w", err)
	}
	return ingredients, nil
}

// MealPlanUsecase handles meal plan business logic
type MealPlanUsecase struct {
	mealPlanRepo   domain.MealPlanRepository
	recipeRepo     domain.RecipeRepository
	shoppingListRepo domain.ShoppingListRepository
}

// NewMealPlanUsecase creates a new meal plan use case
func NewMealPlanUsecase(mealPlanRepo domain.MealPlanRepository, recipeRepo domain.RecipeRepository, shoppingListRepo domain.ShoppingListRepository) *MealPlanUsecase {
	return &MealPlanUsecase{
		mealPlanRepo:   mealPlanRepo,
		recipeRepo:     recipeRepo,
		shoppingListRepo: shoppingListRepo,
	}
}

// CreateMealPlan creates a new meal plan
func (uc *MealPlanUsecase) CreateMealPlan(userID int, name, description string, startDate, endDate time.Time, targetCaloriesPerDay *int) (*domain.MealPlan, error) {
	if endDate.Before(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	mealPlan := &domain.MealPlan{
		UserID:              userID,
		Name:                name,
		Description:         description,
		StartDate:           startDate,
		EndDate:             endDate,
		TargetCaloriesPerDay: targetCaloriesPerDay,
		IsActive:            true,
	}

	if err := uc.mealPlanRepo.Create(mealPlan); err != nil {
		return nil, fmt.Errorf("failed to create meal plan: %w", err)
	}

	return mealPlan, nil
}

// GetMealPlan retrieves a meal plan by ID
func (uc *MealPlanUsecase) GetMealPlan(id int) (*domain.MealPlan, error) {
	mealPlan, err := uc.mealPlanRepo.GetByID(id)
	if err != nil {
		return nil, ErrMealPlanNotFound
	}
	return mealPlan, nil
}

// GetUserMealPlans retrieves all meal plans for a user
func (uc *MealPlanUsecase) GetUserMealPlans(userID int) ([]domain.MealPlan, error) {
	mealPlans, err := uc.mealPlanRepo.GetByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user meal plans: %w", err)
	}
	return mealPlans, nil
}

// GetActiveMealPlan retrieves the active meal plan for a user
func (uc *MealPlanUsecase) GetActiveMealPlan(userID int) (*domain.MealPlan, error) {
	mealPlan, err := uc.mealPlanRepo.GetActive(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active meal plan: %w", err)
	}
	return mealPlan, nil
}

// AddMealEntry adds a meal entry to a meal plan
func (uc *MealPlanUsecase) AddMealEntry(mealPlanID int, recipeID *int, mealType string, mealDate time.Time, calories *int, notes string) (*domain.MealPlanEntry, error) {
	entry := &domain.MealPlanEntry{
		MealPlanID: mealPlanID,
		RecipeID:   recipeID,
		MealType:   mealType,
		MealDate:   mealDate,
		Calories:   calories,
		Notes:      notes,
		IsCompleted: false,
	}

	if err := uc.mealPlanRepo.AddEntry(entry); err != nil {
		return nil, fmt.Errorf("failed to add meal entry: %w", err)
	}

	return entry, nil
}

// GetMealEntries retrieves entries for a meal plan
func (uc *MealPlanUsecase) GetMealEntries(mealPlanID int) ([]domain.MealPlanEntry, error) {
	entries, err := uc.mealPlanRepo.GetEntries(mealPlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal entries: %w", err)
	}
	return entries, nil
}

// GetMealsForDate retrieves meals for a specific date
func (uc *MealPlanUsecase) GetMealsForDate(userID int, date time.Time) ([]domain.MealPlanEntry, error) {
	entries, err := uc.mealPlanRepo.GetEntriesByDate(userID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals for date: %w", err)
	}
	return entries, nil
}

// MarkMealCompleted marks a meal entry as completed
func (uc *MealPlanUsecase) MarkMealCompleted(entryID int) error {
	if err := uc.mealPlanRepo.MarkEntryCompleted(entryID); err != nil {
		return fmt.Errorf("failed to mark meal as completed: %w", err)
	}
	return nil
}

// GenerateShoppingList generates a shopping list from a meal plan
func (uc *MealPlanUsecase) GenerateShoppingList(userID int, mealPlanID int, name string) (*domain.ShoppingList, error) {
	// Create shopping list
	shoppingList := &domain.ShoppingList{
		UserID:     userID,
		Name:       name,
		MealPlanID: &mealPlanID,
		IsActive:   true,
	}

	if err := uc.shoppingListRepo.Create(shoppingList); err != nil {
		return nil, fmt.Errorf("failed to create shopping list: %w", err)
	}

	// Get all meal entries for the meal plan
	entries, err := uc.mealPlanRepo.GetEntries(mealPlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal entries: %w", err)
	}

	// Aggregate ingredients from all recipes
	ingredientMap := make(map[string]*domain.ShoppingListItem)
	
	for _, entry := range entries {
		if entry.RecipeID != nil {
			ingredients, err := uc.recipeRepo.GetIngredients(*entry.RecipeID)
			if err != nil {
				continue
			}

			for _, ingredient := range ingredients {
				if existing, ok := ingredientMap[ingredient.Name]; ok {
					existing.Quantity += ingredient.Quantity
				} else {
					ingredientMap[ingredient.Name] = &domain.ShoppingListItem{
						ShoppingListID: shoppingList.ID,
						IngredientName: ingredient.Name,
						Quantity:       ingredient.Quantity,
						Unit:           ingredient.Unit,
						IsPurchased:    false,
					}
				}
			}
		}
	}

	// Add all aggregated ingredients to shopping list
	for _, item := range ingredientMap {
		if err := uc.shoppingListRepo.AddItem(item); err != nil {
			fmt.Printf("Failed to add ingredient to shopping list: %v\n", err)
		}
	}

	return shoppingList, nil
}