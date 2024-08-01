-- name: GetAllCategories :many
SELECT *
FROM transaction_categories;


-- name: CreateCategory :exec
INSERT INTO transaction_categories(description)
VALUES($1);

-- name: AnyCategoryByName :one
SELECT EXISTS (
    SELECT id
    FROM transaction_categories
    WHERE description = $1
);
