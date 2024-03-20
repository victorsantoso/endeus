package domain

import (
	"context"

	"github.com/victorsantoso/endeus/entity"
)

type RecipeRepository interface {
	// Recipe Categories
	CreateRecipeCategory(ctx context.Context, categoryTag string) error
	GetRecipeCategoryById(ctx context.Context, categoryId int64) (*entity.RecipeCategory, error)
	GetRecipeCategories(ctx context.Context) ([]entity.RecipeCategory, error)
	// Recipes
	CreateRecipe(ctx context.Context, recipe *entity.Recipe) error
	GetRecipeById(ctx context.Context, recipeId int64) (*entity.Recipe, error)
	GetRecipes(ctx context.Context, getRecipesQueryFilter *GetRecipesQueryFilter) ([]entity.Recipe, error)
	UpdateRecipeById(ctx context.Context, recipeId int64, updateRecipeByIdQueryFilter *UpdateRecipeByIdQueryFilter) error
	DeleteRecipeById(ctx context.Context, recipeId int64) error
	// Recipe Ratings
	CreateRecipeRating(ctx context.Context, recipeId, userId int64, rating int) error
	GetRecipeRatingSummary(ctx context.Context, recipeId int64) (float64, int, error)
}

type RecipeUsecase interface {
	// Recipe Categories
	CreateRecipeCategory(ctx context.Context, createRecipeCategoryDTO *CreateRecipeCategoryDTO) error
	GetRecipeCategoryById(ctx context.Context, categoryId int64) (*entity.RecipeCategory, error)
	GetRecipeCategories(ctx context.Context) ([]entity.RecipeCategory, error)
	// Recipes
	CreateRecipe(ctx context.Context, createRecipeDTO *CreateRecipeDTO) error
	GetRecipeById(ctx context.Context, recipeId int64) (*entity.Recipe, error)
	GetRecipes(ctx context.Context, getRecipesQueryFilter *GetRecipesQueryFilter) ([]entity.Recipe, error)
	UpdateRecipe(ctx context.Context, recipeId int64, updateRecipeDTO *UpdateRecipeDTO) error
	DeleteRecipeById(ctx context.Context, recipeId int64) error
	// Recipe Ratings
}

// Recipe Categories
type CreateRecipeCategoryDTO struct {
	CategoryTag string `json:"category_tag" binding:"required,min=3,max=60"`
}

type CreateRecipeCategoryResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type GetRecipeCategoryByIdResponse struct {
	RecipeCategory *entity.RecipeCategory `json:"recipe_category,omitempty"`
	Message        string                 `json:"messagee"`
	Code           int                    `json:"code"`
}

type GetRecipeCategoriesResponse struct {
	RecipeCategories []entity.RecipeCategory `json:"recipe_categories,omitempty"`
	Message          string                  `json:"message"`
	Code             int                     `json:"code"`
}

// Recipes
type CreateRecipeDTO struct {
	Title                string      `json:"title" binding:"required,min=6,max=60"`
	Header               string      `json:"header" binding:"required"`
	ImagePreview         string      `json:"image_preview" binding:"required"`
	Description          string      `json:"description,omitempty"`
	RecipeIngredients    interface{} `json:"recipe_ingredients" binding:"required"`
	CategoryId           int64       `json:"category_id" binding:"required,min=1"`
	EstimatedTimeMinutes int         `json:"estimated_time_minutes" binding:"required,min=3"`
}
type CreateRecipeResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
type GetRecipeByIdResponse struct {
	Recipe  *entity.Recipe `json:"recipe,omitempty"`
	Message string         `json:"message"`
	Code    int            `json:"code"`
}

type GetRecipesQueryFilter struct {
	Name       string
	CategoryId int64
	Limit      int
	Offset     int
}
type GetRecipesResponse struct {
	Recipes []entity.Recipe `json:"recipes,omitempty"`
	Message string          `json:"message"`
	Code    int             `json:"code"`
}

type UpdateRecipeDTO struct {
	Title                string      `json:"title,omitempty"`
	Header               string      `json:"header,omitempty"`
	ImagePreview         string      `json:"image_preview,omitempty"`
	Description          string      `json:"description,omitempty"`
	RecipeIngredients    interface{} `json:"recipe_ingredients,omitempty"`
	CategoryId           int64       `json:"category_id,omitempty" binding:"min=1"`
	EstimatedTimeMinutes int         `json:"estimated_time_minutes,omitempty" binding:"min=3"`
}
type UpdateRecipeByIdQueryFilter struct {
	Title                string
	Header               string
	ImagePreview         string
	Description          string
	RecipeIngredients    interface{}
	EstimatedTimeMinutes int
}
type UpdateRecipeResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type DeleteRecipeResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// // Recipe Ratings
// type CreateRecipeRatingDTO struct {
// 	Rating int `json:"rating" binding:"required,min=1,max=5"`
// }
