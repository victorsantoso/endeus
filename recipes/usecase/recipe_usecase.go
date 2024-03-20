package usecase

import (
	"context"
	"database/sql"

	"github.com/apex/log"
	"github.com/victorsantoso/endeus/domain"
	"github.com/victorsantoso/endeus/entity"
)

type recipeUsecase struct {
	recipeRepository domain.RecipeRepository
}

func NewRecipeUsecase(recipeRepository domain.RecipeRepository) domain.RecipeUsecase {
	return &recipeUsecase{
		recipeRepository: recipeRepository,
	}
}

// Recipe Categories
func (ru *recipeUsecase) CreateRecipeCategory(ctx context.Context, createRecipeCategoryDTO *domain.CreateRecipeCategoryDTO) error {
	if err := ru.recipeRepository.CreateRecipeCategory(ctx, createRecipeCategoryDTO.CategoryTag); err != nil {
		log.Errorf("[recipe_usecase.CreateRecipeCategory] error creating recipe category, err: %v", err)
		return err
	}
	log.Debug("[recipe_usecase.CreateRecipeCategory] successfully created a new recipe category")
	return nil
}
func (ru *recipeUsecase) GetRecipeCategoryById(ctx context.Context, categoryId int64) (*entity.RecipeCategory, error) {
	recipeCategory, err := ru.recipeRepository.GetRecipeCategoryById(ctx, categoryId)
	if err != nil {
		if err == sql.ErrNoRows || recipeCategory == nil {
			log.Debugf("[recipe_usecase.GetRecipeCategoryById] no row found for category_id: %d, err: %v", categoryId, err)
			return nil, sql.ErrNoRows
		}
		log.Errorf("[recipe_usecase.GetRecipeCategoryById] error getting category by id, err: %v", err)
		return nil, err
	}
	log.Debugf("[recipe_usecase.GetRecipeCategoryById] recipe_category: %+v", recipeCategory)
	return recipeCategory, nil
}
func (ru *recipeUsecase) GetRecipeCategories(ctx context.Context) ([]entity.RecipeCategory, error) {
	recipeCategories, err := ru.recipeRepository.GetRecipeCategories(ctx)
	if err != nil {
		if err == sql.ErrNoRows || len(recipeCategories) == 0 {
			log.Debugf("[recipe_usecase.GetRecipeCategories] recipe categories are empty, err:%v", err)
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return recipeCategories, nil
}

// Recipes
func (ru *recipeUsecase) CreateRecipe(ctx context.Context, createRecipeDTO *domain.CreateRecipeDTO) error {
	if err := ru.recipeRepository.CreateRecipe(ctx, &entity.Recipe{
		Title:                createRecipeDTO.Title,
		Header:               createRecipeDTO.Header,
		ImagePreview:         createRecipeDTO.ImagePreview,
		Description:          createRecipeDTO.Description,
		RecipeIngredients:    createRecipeDTO.RecipeIngredients,
		CategoryId:           createRecipeDTO.CategoryId,
		EstimatedTimeMinutes: createRecipeDTO.EstimatedTimeMinutes,
	}); err != nil {
		log.Errorf("[recipe_usecase.CreateRecipe] error creating a new recipe, err: %v", err)
		return err
	}
	return nil
}
func (ru *recipeUsecase) GetRecipeById(ctx context.Context, recipeId int64) (*entity.Recipe, error) {
	recipe, err := ru.recipeRepository.GetRecipeById(ctx, recipeId)
	if err != nil {
		if err == sql.ErrNoRows || recipe == nil {
			log.Debugf("[recipe_usecase.GetRecipeById] no row found for recipe_id: %d, err: %v", recipeId, err)
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return recipe, nil
}
func (ru *recipeUsecase) GetRecipes(ctx context.Context, getRecipesQueryFilter *domain.GetRecipesQueryFilter) ([]entity.Recipe, error) {
	recipes, err := ru.recipeRepository.GetRecipes(ctx, getRecipesQueryFilter)
	if err != nil {
		if err == sql.ErrNoRows || recipes == nil {
			log.Debugf("[recipe_usecase.GetRecipes] no rows found, err: %v", err)
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return recipes, nil
}
func (ru *recipeUsecase) UpdateRecipe(ctx context.Context, recipeId int64, updateRecipeDTO *domain.UpdateRecipeDTO) error {
	if err := ru.recipeRepository.UpdateRecipeById(ctx, recipeId, &domain.UpdateRecipeByIdQueryFilter{
		Title:                updateRecipeDTO.Title,
		Header:               updateRecipeDTO.Header,
		ImagePreview:         updateRecipeDTO.Header,
		Description:          updateRecipeDTO.Description,
		RecipeIngredients:    updateRecipeDTO.RecipeIngredients,
		EstimatedTimeMinutes: updateRecipeDTO.EstimatedTimeMinutes,
	}); err != nil {
		log.Errorf("[recipe_usecase.UpdateRecipe] error updating recipe with recipe_id: %d, err: %v", recipeId, err)
		return err
	}
	return nil
}
func (ru *recipeUsecase) DeleteRecipeById(ctx context.Context, recipeId int64) error {
	if err := ru.recipeRepository.DeleteRecipeById(ctx, recipeId); err != nil {
		if err == sql.ErrNoRows {
			log.Errorf("[recipe_usecase.DeleteRecipeById] error deleting recipe with recipe_id: %d, err: %v", recipeId, err)
		}
		return err
	}
	return nil
}

// Recipe Ratings
func (ru *recipeUsecase) CreateRecipeRating(ctx context.Context, recipeId, userId int64, rating int) error {
	if err := ru.recipeRepository.CreateRecipeRating(ctx, recipeId, userId, rating); err != nil {
		if err == sql.ErrNoRows {
			log.Debugf("[recipe_usecase.CreateRecipeRating] error creating recipe rating with recipe_id: %d, err: %v", recipeId, err)
		}
		return err
	}
	return nil
}
func (ru *recipeUsecase) GetRecipeRatingSummary(ctx context.Context, recipeId int64) (float64, int, error) {
	avg, ratingCount, err := ru.recipeRepository.GetRecipeRatingSummary(ctx, recipeId)
	if err != nil {
		log.Debugf("[recipe_usecase.GetRecipeRatingSummary] error getting recipe rating summary with recipe_id: %d, err: %v", recipeId, err)
		return 0, 0, err
	}
	return avg, ratingCount, nil
}
