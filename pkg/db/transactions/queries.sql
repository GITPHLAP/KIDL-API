-- name: InsertNewTransactions :copyfrom
INSERT INTO public.transactions_test(accounting, amount, sender, recipient, "name", description, balance, currency, aligned, reserved) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetLatestTransaction :one
SELECT * 
FROM transactions_test
ORDER BY accounting DESC, id DESC 
LIMIT 1;

-- name: UpdateTransactionCategoryAndAccounting :exec
UPDATE transactions_test
SET user_category = $1, user_accounting=$2
WHERE id = $3;