CREATE TABLE transaction_categories_mapping(
	id SERIAL4 PRIMARY KEY,
	category_id INT NOT NULL,
	mapping_value text NOT NULL,
	CONSTRAINT unique_mapping_values UNIQUE (mapping_value),
	CONSTRAINT fk_transaction_categories_mapping_category_id FOREIGN KEY (category_id) REFERENCES transaction_categories(id)
);