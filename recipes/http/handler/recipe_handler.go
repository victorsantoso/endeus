package handler

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorsantoso/endeus/domain"
	"github.com/victorsantoso/endeus/entity"
)

type recipeHandler struct {
	recipeUsecase domain.RecipeUsecase
}

func NewRecipeHandler(g *gin.Engine, authMiddleware gin.HandlerFunc, recipeUsecase domain.RecipeUsecase) {
	recipeHandler := &recipeHandler{
		recipeUsecase: recipeUsecase,
	}

	// No Auth needed to access this group
	noAuthGroup := g.Group("/api/v1")
	// Recipe Category
	noAuthGroup.GET("/recipe_categories", recipeHandler.GetRecipeCategories)
	noAuthGroup.GET("/recipe_category/:recipeCategoryId", recipeHandler.GetRecipeCategoryById)
	// Recipe
	noAuthGroup.GET("/recipe/:recipeId", recipeHandler.GetRecipeById)
	noAuthGroup.GET("/recipes", recipeHandler.GetRecipes)

	// Auth group with ADMIN role only
	authGroup := g.Group("/api/v1", authMiddleware)
	// Recipe Category
	authGroup.POST("/recipe_category", recipeHandler.CreateRecipeCategory)
	// Recipe
	authGroup.POST("/recipe", recipeHandler.CreateRecipe)
	authGroup.PUT("/recipe/:recipeId", recipeHandler.UpdateRecipe)
	authGroup.DELETE("/recipe/:recipeId", recipeHandler.DeleteRecipe)
}

// Recipe Category
func (rh *recipeHandler) CreateRecipeCategory(c *gin.Context) {
	key, _ := c.Get("user")
	user, ok := key.(*entity.User)
	if !ok || user.Role != domain.ADMIN {
		c.JSON(http.StatusForbidden, &domain.CreateRecipeCategoryResponse{
			Message: domain.ErrForbidenAccess.Error(),
			Code:    http.StatusForbidden,
		})
		return
	}
	createRecipeCategoryDTO := &domain.CreateRecipeCategoryDTO{}
	if err := c.ShouldBindJSON(createRecipeCategoryDTO); err != nil {
		c.JSON(http.StatusBadRequest, &domain.CreateRecipeCategoryResponse{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	if err := rh.recipeUsecase.CreateRecipeCategory(context.Background(), createRecipeCategoryDTO); err != nil {
		c.JSON(http.StatusInternalServerError, &domain.CreateRecipeCategoryResponse{
			Message: domain.ErrInternalServerError.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, &domain.CreateRecipeCategoryResponse{
		Message: "successfully created a new category",
		Code:    http.StatusOK,
	})
}
func (rh *recipeHandler) GetRecipeCategoryById(c *gin.Context) {
	recipeCategoryParam := c.Param("recipeCategoryId")
	categoryId, err := strconv.Atoi(recipeCategoryParam)
	if err != nil || categoryId <= 0 {
		c.JSON(http.StatusBadRequest, &domain.GetRecipeCategoryByIdResponse{
			Message: domain.ErrBadRequest.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	recipeCategory, err := rh.recipeUsecase.GetRecipeCategoryById(context.Background(), int64(categoryId))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, &domain.GetRecipeCategoryByIdResponse{
				Message: domain.ErrNotFound.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &domain.GetRecipeCategoryByIdResponse{
			Message: domain.ErrInternalServerError.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, &domain.GetRecipeCategoryByIdResponse{
		RecipeCategory: recipeCategory,
		Message:        "successfully retrieved recipe category by id",
		Code:           http.StatusOK,
	})
}
func (rh *recipeHandler) GetRecipeCategories(c *gin.Context) {
	recipeCategories, err := rh.recipeUsecase.GetRecipeCategories(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, &domain.GetRecipeCategoriesResponse{
				Message: domain.ErrNotFound.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &domain.GetRecipeCategoriesResponse{
			Message: domain.ErrInternalServerError.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, &domain.GetRecipeCategoriesResponse{
		RecipeCategories: recipeCategories,
		Message:          "successfully get all recipe categories",
		Code:             http.StatusOK,
	})
}

// Recipe
func (rh *recipeHandler) CreateRecipe(c *gin.Context) {
	key, _ := c.Get("user")
	user, ok := key.(*entity.User)
	if !ok || user.Role != domain.ADMIN {
		c.JSON(http.StatusForbidden, &domain.CreateRecipeResponse{
			Message: domain.ErrForbidenAccess.Error(),
			Code:    http.StatusForbidden,
		})
		return
	}
	createRecipeDTO := &domain.CreateRecipeDTO{}
	if err := c.ShouldBindJSON(createRecipeDTO); err != nil {
		c.JSON(http.StatusBadRequest, &domain.CreateRecipeResponse{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	if err := rh.recipeUsecase.CreateRecipe(context.Background(), createRecipeDTO); err != nil {
		c.JSON(http.StatusInternalServerError, &domain.CreateRecipeResponse{
			Message: domain.ErrInternalServerError.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, &domain.CreateRecipeResponse{
		Message: "successfully created a new recipe",
		Code:    http.StatusOK,
	})
}
func (rh *recipeHandler) GetRecipeById(c *gin.Context) {
	recipeParam := c.Param("recipeId")
	recipeId, err := strconv.Atoi(recipeParam)
	if err != nil || recipeId <= 0 {
		c.JSON(http.StatusBadRequest, &domain.GetRecipeByIdResponse{
			Message: domain.ErrInvalidId.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	recipe, err := rh.recipeUsecase.GetRecipeById(context.Background(), int64(recipeId))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, &domain.GetRecipeByIdResponse{
				Message: domain.ErrNotFound.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &domain.GetRecipeByIdResponse{
			Message: domain.ErrInternalServerError.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, &domain.GetRecipeByIdResponse{
		Recipe:  recipe,
		Message: "successfully retrieved recipe by id",
		Code:    http.StatusOK,
	})
}
func (rh *recipeHandler) GetRecipes(c *gin.Context) {
	// get queries
	nameQuery := c.Query("name")
	categoryIdQuery := c.Query("category_id")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	queryFilter := &domain.GetRecipesQueryFilter{}
	if len(nameQuery) > 0 {
		queryFilter.Name = nameQuery
	}
	categoryId, _ := strconv.Atoi(categoryIdQuery)
	if categoryId > 0 {
		queryFilter.CategoryId = int64(categoryId)
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}
	queryFilter.Limit = limitInt
	queryFilter.Offset = offsetInt

	recipes, err := rh.recipeUsecase.GetRecipes(context.Background(), queryFilter)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, &domain.GetRecipesResponse{
				Message: domain.ErrNotFound.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &domain.GetRecipesResponse{
			Message: domain.ErrInternalServerError.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, &domain.GetRecipesResponse{
		Recipes: recipes,
		Message: "successfully retrieved recipes",
		Code:    http.StatusOK,
	})
}
func (rh *recipeHandler) UpdateRecipe(c *gin.Context) {
	key, _ := c.Get("user")
	user, ok := key.(*entity.User)
	if !ok || user.Role != domain.ADMIN {
		c.JSON(http.StatusForbidden, &domain.UpdateRecipeResponse{
			Message: domain.ErrForbidenAccess.Error(),
			Code:    http.StatusForbidden,
		})
		return
	}
	recipeParam := c.Param("recipeId")
	recipeId, err := strconv.Atoi(recipeParam)
	if err != nil || recipeId <= 0 {
		c.JSON(http.StatusBadRequest, &domain.UpdateRecipeResponse{
			Message: domain.ErrInvalidId.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	updateRecipeDTO := &domain.UpdateRecipeDTO{}
	if err := c.ShouldBindJSON(updateRecipeDTO); err != nil {
		c.JSON(http.StatusBadRequest, &domain.UpdateRecipeResponse{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	err = rh.recipeUsecase.UpdateRecipe(context.Background(), int64(recipeId), updateRecipeDTO)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, &domain.UpdateRecipeResponse{
				Message: domain.ErrNotFound.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &domain.UpdateRecipeResponse{
			Message: domain.ErrInternalServerError.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, &domain.UpdateRecipeResponse{
		Message: "succesfully updated recipe",
		Code:    http.StatusOK,
	})
}

func (rh *recipeHandler) DeleteRecipe(c *gin.Context) {
	key, _ := c.Get("user")
	user, ok := key.(*entity.User)
	if !ok || user.Role != domain.ADMIN {
		c.JSON(http.StatusForbidden, &domain.DeleteRecipeResponse{
			Message: domain.ErrForbidenAccess.Error(),
			Code:    http.StatusForbidden,
		})
		return
	}
	recipeParam := c.Param("recipeId")
	recipeId, err := strconv.Atoi(recipeParam)
	if err != nil || recipeId <= 0 {
		c.JSON(http.StatusBadRequest, &domain.DeleteRecipeResponse{
			Message: domain.ErrInvalidId.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	err = rh.recipeUsecase.DeleteRecipeById(context.Background(), int64(recipeId))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, &domain.DeleteRecipeResponse{
				Message: domain.ErrNotFound.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &domain.DeleteRecipeResponse{
			Message: domain.ErrInternalServerError.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, &domain.DeleteRecipeResponse{
		Message: "succesfully updated recipe",
		Code:    http.StatusOK,
	})
}
