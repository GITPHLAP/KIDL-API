-- name: GetOverviewByYearNType :many 
SELECT *
FROM monthly_overview_by_category
WHERE year = $1 AND type = $2;


-- name: GetTotalOverviewBy :many 
SELECT year, month, type, SUM(actual_sum), SUM(forecast_sum)
FROM monthly_overview_by_category
WHERE year = $1 AND type = $2
GROUP BY year, month, type;

