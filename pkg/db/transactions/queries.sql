-- name: InsertNewTransactions :copyfrom
INSERT INTO public.transactions(accounting, amount, sender, recipient, "name", description, balance, currency, aligned, reserved) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);