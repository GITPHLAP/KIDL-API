CREATE TABLE public.transactions (
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
	user_category int NULL,
	reserved bool NULL,
	user_accounting date NULL,
	CONSTRAINT transactions_amount_sender_recipient_name_description_balan_key UNIQUE (amount, sender, recipient, name, description, balance, currency)
);
