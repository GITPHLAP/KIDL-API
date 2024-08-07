CREATE TABLE monthly_overview_by_category (
	YEAR int NOT NULL,
	month int NOT NULL,
	type public.trans_type,
	category text NOT NULL,
	actual_sum FLOAT8,
	forecast_sum FLOAT8
);