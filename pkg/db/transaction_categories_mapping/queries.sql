-- name: CreateCategoryMapping :exec
INSERT INTO transaction_categories_mapping(category_id, mapping_value)
VALUES($1, $2);

-- name: GetCategoryMapping :many
SELECT *
FROM transaction_categories_mapping;

-- name: UpdateCategoryMapping :exec
UPDATE transaction_categories_mapping
SET category_id = $1
    , mapping_value = $2
WHERE id = $3;

-- name: DeleteCategoryMapping :exec
DELETE FROM transaction_categories_mapping
WHERE id = $1;