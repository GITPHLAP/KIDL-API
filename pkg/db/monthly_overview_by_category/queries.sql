-- name: GetOverviewByYearNMonth :one 
SELECT jsonb_agg(y) 
FROM (
	SELECT "year", json_agg(ym.mcaf) AS months
	FROM (
			SELECT "year", json_build_object( "month", json_agg(ymt.caf)) AS mcaf
			FROM (
				SELECT mobc."year", mobc."month", json_build_object(mobc."type", json_agg(
					json_build_object(
						'category', category, 'actual_sum', actual_sum, 'forecast_sum', forecast_sum
					)
				), mobc."type" || '_total', 					json_build_object(
						'total_actual', SUM(actual_sum) , 'total_forecast', SUM(forecast_sum)
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