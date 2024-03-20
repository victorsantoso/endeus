package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/victorsantoso/endeus/domain"
	"github.com/victorsantoso/endeus/entity"
)

type recipeRepository struct {
	dbConn *sql.DB
}

func NewRecipeRepository(dbConn *sql.DB) domain.RecipeRepository {
	return &recipeRepository{
		dbConn: dbConn,
	}
}

const (
	// query with variables handling to escape sql injections

	// Recipe Categories
	CreateRecipeCategoryQuery = `
		INSERT INTO recipe_categories(category_tag)
		VALUES($1);
	`
	GetRecipeCategoryByIdQuery = `
		SELECT category_id, category_tag
		FROM recipe_categories
		WHERE category_id = $1
		LIMIT 1
	`
	GetRecipeCategoriesQuery = `
		SELECT category_id, category_tag
		FROM recipe_categories
	`
	// Recipes
	CreateRecipeQuery = `
		INSERT INTO recipes(category_id, title, header, image_preview, description, estimated_time_minutes, recipe_ingredients, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, now()::timestamptz, now()::timestamptz);
	`
	GetRecipeByIdQuery = `
		SELECT recipe_id, category_id, title, header, image_preview, description, estimated_time_minutes, recipe_ingredients, created_at, updated_at
		FROM recipes
		WHERE recipe_id = $1
	`
	GetRecipesQuery = `
		SELECT recipe_id, category_id, title, header, image_preview, description, estimated_time_minutes, recipe_ingredients, created_at, updated_at
		FROM recipes
	`
	UpdateRecipeByIdQuery = `
		UPDATE recipes
		SET
			updated_at = now()::timestamptz
	`
	DeleteRecipeByIdQuery = `
		DELETE FROM recipes
		WHERE recipe_id = $1
	`
	// Recipe Ratings
	CreateRecipeRatingQuery = `
		INSERT INTO recipe_ratings(recipe_id, user_id, recipe_rating, created_at, updated_at)
		VALUES($1, $2, $3, now()::timestamptz, now()::timestamptz);
	`
	GetRecipeRatingSummaryQuery = `
		SELECT AVG(recipe_rating) AS average_rating, COUNT(recipe_rating) AS rating_count
		FROM recipe_ratings
		WHERE recipe_id = $1
		GROUP BY recipe_id;
	`
)

// Recipe Categories
func (rr *recipeRepository) CreateRecipeCategory(ctx context.Context, categoryTag string) error {
	tx, err := rr.dbConn.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, CreateRecipeCategoryQuery, categoryTag)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (rr *recipeRepository) GetRecipeCategoryById(ctx context.Context, categoryId int64) (*entity.RecipeCategory, error) {
	var recipeCategory entity.RecipeCategory
	row := rr.dbConn.QueryRowContext(ctx, GetRecipeCategoryByIdQuery, categoryId)
	if err := row.Scan(&recipeCategory.CategoryId, &recipeCategory.CategoryTag); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &recipeCategory, nil
}
func (rr *recipeRepository) GetRecipeCategories(ctx context.Context) ([]entity.RecipeCategory, error) {
	var recipeCategories []entity.RecipeCategory
	rows, err := rr.dbConn.QueryContext(ctx, GetRecipeCategoriesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recipeCategory entity.RecipeCategory
	for rows.Next() {
		if err := rows.Scan(&recipeCategory.CategoryId, &recipeCategory.CategoryTag); err != nil {
			return nil, err
		}
		recipeCategories = append(recipeCategories, recipeCategory)
	}
	return recipeCategories, nil
}

// Recipes
func (rr *recipeRepository) CreateRecipe(ctx context.Context, recipe *entity.Recipe) error {
	tx, err := rr.dbConn.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, CreateRecipeQuery, recipe.CategoryId, recipe.Title, recipe.Header, recipe.ImagePreview, recipe.Description, recipe.EstimatedTimeMinutes, recipe.RecipeIngredients, recipe.CreatedAt, recipe.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (rr *recipeRepository) GetRecipeById(ctx context.Context, recipeId int64) (*entity.Recipe, error) {
	var recipe entity.Recipe
	row := rr.dbConn.QueryRowContext(ctx, GetRecipeByIdQuery, recipeId)
	if err := row.Scan(&recipe.RecipeId, &recipe.CategoryId, &recipe.Title, &recipe.Header, &recipe.ImagePreview, &recipe.Description, &recipe.EstimatedTimeMinutes, &recipe.RecipeIngredients, &recipe.CreatedAt, &recipe.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &recipe, nil
}
func (rr *recipeRepository) GetRecipes(ctx context.Context, getRecipesQueryFilter *domain.GetRecipesQueryFilter) ([]entity.Recipe, error) {
	// var recipes []entity.Recipe
	// rows, err := rr.dbConn.QueryContext(ctx, GetRecipesQuery, 1)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// var recipe entity.Recipe
	// for rows.Next() {
	// 	if err := rows.Scan(&recipe.RecipeId, &recipe.CategoryId, &recipe.Title, &recipe.Header, &recipe.ImagePreview, &recipe.Description, &recipe.EstimatedTimeMinutes, &recipe.RecipeIngredients, &recipe.CreatedAt, &recipe.UpdatedAt); err != nil {
	// 		return nil, err
	// 	}
	// 	recipes = append(recipes, recipe)
	// }
	// return recipes, nil

	var recipes []entity.Recipe
	limit := getRecipesQueryFilter.Limit
	offset := getRecipesQueryFilter.Offset
	query := GetRecipesQuery
	var args []interface{}
	var conditions []string
	queryParamCount := 1
	// Handle filtering by name (title)
	if getRecipesQueryFilter.Name != "" {
		conditions = append(conditions, "LOWER(title) = $"+fmt.Sprintf("%d", queryParamCount))
		args = append(args, strings.ToLower(getRecipesQueryFilter.Name))
		queryParamCount++
	}
	if getRecipesQueryFilter.CategoryId != 0 {
		conditions = append(conditions, "category_id = $"+fmt.Sprintf("%d", queryParamCount))
		args = append(args, getRecipesQueryFilter.CategoryId)
		queryParamCount++
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", queryParamCount, queryParamCount+1)
	args = append(args, limit, offset)
	rows, err := rr.dbConn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe entity.Recipe
		if err := rows.Scan(&recipe.RecipeId, &recipe.CategoryId, &recipe.Title, &recipe.Header, &recipe.ImagePreview, &recipe.Description, &recipe.EstimatedTimeMinutes, &recipe.RecipeIngredients, &recipe.CreatedAt, &recipe.UpdatedAt); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (rr *recipeRepository) UpdateRecipeById(ctx context.Context, recipeId int64, updateRecipeByIdQueryFilter *domain.UpdateRecipeByIdQueryFilter) error {
	updateRecipeQuery := UpdateRecipeByIdQuery
	var updateValues []interface{}
	var args []interface{}
	queryFilterCount := 1
	// Check each field in the filter and update the query accordingly
	if updateRecipeByIdQueryFilter.Title != "" {
		updateRecipeQuery += ", title = $" + fmt.Sprintf("%d", queryFilterCount)
		updateValues = append(updateValues, updateRecipeByIdQueryFilter.Title)
		queryFilterCount++
	}
	if updateRecipeByIdQueryFilter.Header != "" {
		updateRecipeQuery += ", header = $" + fmt.Sprintf("%d", queryFilterCount)
		updateValues = append(updateValues, updateRecipeByIdQueryFilter.Header)
		queryFilterCount++
	}
	if updateRecipeByIdQueryFilter.ImagePreview != "" {
		updateRecipeQuery += ", image_preview = $" + fmt.Sprintf("%d", queryFilterCount)
		updateValues = append(updateValues, updateRecipeByIdQueryFilter.ImagePreview)
		queryFilterCount++
	}
	if updateRecipeByIdQueryFilter.Description != "" {
		updateRecipeQuery += ", description = $" + fmt.Sprintf("%d", queryFilterCount)
		updateValues = append(updateValues, updateRecipeByIdQueryFilter.Description)
		queryFilterCount++
	}
	updateRecipeQuery += " WHERE recipe_id = $1"
	args = append(args, recipeId)
	_, err := rr.dbConn.ExecContext(ctx, updateRecipeQuery, append(updateValues, args...)...)
	if err != nil {
		return err
	}
	return nil
}

func (rr *recipeRepository) DeleteRecipeById(ctx context.Context, recipeId int64) error {
	_, err := rr.dbConn.ExecContext(ctx, DeleteRecipeByIdQuery, recipeId)
	if err != nil {
		return err
	}
	return nil
}

// Recipe Ratings
func (rr *recipeRepository) CreateRecipeRating(ctx context.Context, recipeId, userId int64, rating int) error {
	tx, err := rr.dbConn.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, CreateRecipeRatingQuery, recipeId, userId, rating)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (rr *recipeRepository) GetRecipeRatingSummary(ctx context.Context, recipeId int64) (float64, int, error) {
	var ratingAvg float64
	var ratingCount int
	row := rr.dbConn.QueryRowContext(ctx, GetRecipeRatingSummaryQuery, recipeId)
	if err := row.Scan(&ratingAvg, &ratingCount); err != nil {
		return 0, 0, err
	}
	return ratingAvg, ratingCount, nil
}
