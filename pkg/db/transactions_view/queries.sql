-- name: GetTransactionsFromViewByYear :many
SELECT *
FROM transactions_view
WHERE date_part('year', accounting) = @year::int;