CREATE TABLE public.transactions_view (
	id serial4 NOT NULL,
	accounting date NULL,
	amount float8 NULL,
	sender int8 NULL,
	recipient int8 NULL,
	"name" varchar(255) NULL,
	description varchar(255) NULL,
	balance float8 NULL,
	currency varchar(5) NULL,
	aligned varchar(255) NULL,
	reserved bool NULL,
	category_id int NULL,
	type public.trans_type
);
