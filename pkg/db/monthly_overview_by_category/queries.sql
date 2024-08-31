-- name: GetOverviewByYearNMonth :one 
SELECT jsonb_agg(y) 
FROM (
	SELECT "year", json_agg(ym.mcaf) AS months
	FROM (
			SELECT "year", json_build_object( "month", json_agg(ymt.caf)) AS mcaf
			FROM (
				SELECT mobc."year", mobc."month", json_build_object(mobc."type", json_agg(
					json_build_object(
						'category', category, 'actual_sum', actual_sum, 'plan_sum', plan_sum
					)
				), mobc."type" || '_total', 					json_build_object(
						'total_actual', SUM(actual_sum) , 'total_plan', SUM(plan_sum)
					)
				) AS caf
				FROM monthly_overview_by_category mobc
                WHERE year = $1 AND month = ANY(@months::int[])
				GROUP BY mobc.YEAR
					, mobc.MONTH
					, mobc.TYPE
			) ymt
			GROUP BY YEAR
				, MONTH
		) ym
	GROUP BY YEAR
) y;

-- name: Get13thMonthByYear :one 
WITH year_overview AS (
	SELECT YEAR, TYPE, category, SUM(actual_sum) AS actual_sum , sum(plan_sum) AS plan_sum 
	FROM monthly_overview_by_category mobc 
	WHERE YEAR = $1
	GROUP BY YEAR, TYPE, category
)
SELECT "year", json_agg(ym.caf) AS calculation
FROM (
	SELECT yo."year", json_build_object(yo."type", json_agg(
		json_build_object(
			'category', category, 'actual_sum', actual_sum, 'plan_sum', plan_sum
		)
	), yo."type" || '_total', 					json_build_object(
			'total_actual', SUM(actual_sum) , 'total_plan', SUM(plan_sum)
		)
	) AS caf
	FROM year_overview yo
	WHERE year = $1 
	GROUP BY yo.YEAR
		, yo.TYPE
) ym
GROUP BY YEAR;
	
