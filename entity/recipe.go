package entity

import "time"

type RecipeCategory struct {
	CategoryTag string `json:"category_tag"`
	CategoryId  int64  `json:"category_id"`
}

// Recipe will have adjusted memory padding to optimize memory
type Recipe struct {
	Title                string      `json:"title"`
	Header               string      `json:"header"`
	ImagePreview         string      `json:"image_preview"`
	Description          string      `json:"description"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
	RecipeIngredients    interface{} `json:"recipe_ingredients"`
	RecipeId             int64       `json:"recipe_id"`
	CategoryId           int64       `json:"category_id"`
	EstimatedTimeMinutes int         `json:"estimated_time_minutes"`
}

type RecipeRating struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RecipeId     int64     `json:"recipe_id"`
	UserId       int64     `json:"user_id"`
	RecipeRating int       `json:"recipe_rating"`
}

type DiscussionTestimonial struct {
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	Discussion              string    `json:"discussion"`
	DiscussionImage         string    `json:"discussion_image"`
	DiscussionTestimonialId int64     `json:"discussion_testimonial_id"`
	RecipeId                int64     `json:"recipe_id"`
	UserId                  int64     `json:"user_id"`
}
