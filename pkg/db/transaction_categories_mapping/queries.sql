-- name: CreateCategoryMapping :exec
INSERT INTO transaction_categories_mapping(category_id, mapping_value)
VALUES($1, $2);
