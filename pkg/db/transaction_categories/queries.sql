-- name: GetAllCategories :many
SELECT *
FROM transaction_categories
ORDER BY order_no;


-- name: CreateCategory :exec
INSERT INTO transaction_categories(description, comment)
VALUES($1, $2);

-- name: AnyCategoryByName :one
SELECT EXISTS (
    SELECT id
    FROM transaction_categories
    WHERE description = $1
);

-- name: UpdateCategory :exec
UPDATE transaction_categories
SET description = $1, comment = $3
WHERE id =$2;

-- name: DeleteCategoryById :exec
DELETE FROM transaction_categories
WHERE id = $1;

-- name: AnyCategoryByNameAndId :one
SELECT EXISTS (
    SELECT id
    FROM transaction_categories
    WHERE description = $1 AND id != $2
);