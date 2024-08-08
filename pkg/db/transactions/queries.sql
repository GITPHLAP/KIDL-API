-- name: InsertNewTransactions :copyfrom
INSERT INTO public.transactions(accounting, amount, sender, recipient, "name", description, balance, currency, aligned, reserved) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetLatestTransaction :one
SELECT * 
FROM transactions
ORDER BY accounting DESC, id DESC 
LIMIT 1;

-- name: UpdateTransactionCategoryAndAccounting :exec
UPDATE transactions
SET user_category_id = $1, user_accounting=$2
WHERE id = $3;