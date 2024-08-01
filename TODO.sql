CREATE TABLE transaction_categories(
	id SERIAL4 PRIMARY KEY,
	description text NOT NULL
);

CREATE TABLE forecasts (
	id SERIAL8 PRIMARY KEY,
	YEAR int NOT NULL,
	category_id INT NOT NULL,
	january_amount FLOAT8,
	february_amount FLOAT8,
	march_amount FLOAT8,
	april_amount FLOAT8,
	may_amount FLOAT8,
	june_amount FLOAT8,
	july_amount FLOAT8,
	august_amount FLOAT8,
	september_amount FLOAT8,
	october_amount FLOAT8,
	november_amount FLOAT8,
	december_amount FLOAT8,
	CONSTRAINT fk_forecasts_category_id FOREIGN KEY (category_id) REFERENCES transaction_categories(id)
);

CREATE TABLE transaction_categories_mapping(
	category_id INT NOT NULL,
	mapping_value text NOT NULL,
	CONSTRAINT unique_mapping_values UNIQUE (mapping_value),
	CONSTRAINT fk_transaction_categories_mapping_category_id FOREIGN KEY (category_id) REFERENCES transaction_categories(id)
);