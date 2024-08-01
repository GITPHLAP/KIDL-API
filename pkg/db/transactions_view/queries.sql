-- name: GetTransactionsFromViewByYear :many
SELECT id
    , accounting
    , amount
    , description
    , category
    , CASE 
        WHEN accounting > 0 THEN 'Indtægt' 
        ELSE 'Udgift'
    END AS type
    , balance 
FROM transactions_view
WHERE date_part('year', accounting) = @year::int;