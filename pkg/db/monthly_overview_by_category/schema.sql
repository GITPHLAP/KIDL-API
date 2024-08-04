CREATE TABLE monthly_overview_by_category (
	YEAR int NOT NULL,
	month int NOT NULL,
	type text,
	category text NOT NULL,
	actual_sum FLOAT8,
	forecast_sum FLOAT8
);