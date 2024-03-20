-- Role Enum
CREATE TYPE role AS ENUM('ADMIN', 'READER');

-- Recipe Categories Table
CREATE TABLE public.recipe_categories (
    category_id SERIAL PRIMARY KEY NOT NULL,
    category_tag VARCHAR(60) NOT NULL
);

-- Users Table
CREATE TABLE public.users (
    user_id SERIAL PRIMARY KEY NOT NULL,
    role role DEFAULT 'READER' NOT NULL,
    email VARCHAR(60) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(60) NOT NULL,
    profile_image TEXT DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- Recipes Table
CREATE TABLE public.recipes (
    recipe_id SERIAL PRIMARY KEY NOT NULL, 
    category_id INTEGER NOT NULL,
    title VARCHAR(60) NOT NULL,
    header TEXT NOT NULL,
    image_preview TEXT NOT NULL,
    description TEXT DEFAULT NULL,
    estimated_time_minutes INTEGER NOT NULL,
    recipe_ingredients JSON NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_recipes_category_id FOREIGN KEY(category_id) REFERENCES recipe_categories(category_id)
);

-- Recipe Ratings Table
CREATE TABLE public.recipe_ratings (
    recipe_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    recipe_rating INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_recipe_ratings_recipe_id FOREIGN KEY(recipe_id) REFERENCES recipes(recipe_id),
    CONSTRAINT fk_recipe_ratings_user_id FOREIGN KEY(user_id) REFERENCES users(user_id)
);

-- Not indexed yet for searching etc