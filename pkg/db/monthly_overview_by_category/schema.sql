CREATE TABLE monthly_overview_by_category (
	YEAR int NOT NULL,
	month int NOT NULL,
	type text,
	category_id INT NOT NULL,
	actual_sum FLOAT8,
	forecast_sum FLOAT8
);