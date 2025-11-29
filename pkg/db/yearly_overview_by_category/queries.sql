-- name: GetYearOverviewByYearAndType :many
SELECT category
	, month_1_plan_sum, month_1_actual_sum
	, month_2_plan_sum, month_2_actual_sum
	, month_3_plan_sum, month_3_actual_sum
	, month_4_plan_sum, month_4_actual_sum
	, month_5_plan_sum, month_5_actual_sum
	, month_6_plan_sum, month_6_actual_sum
	, month_7_plan_sum, month_7_actual_sum
	, month_8_plan_sum, month_8_actual_sum
	, month_9_plan_sum, month_9_actual_sum
	, month_10_plan_sum, month_10_actual_sum
	, month_11_plan_sum, month_11_actual_sum
	, month_12_plan_sum, month_12_actual_sum
FROM yearly_overview_by_category
WHERE year = $1 AND type = $2
ORDER BY category_order_no;
