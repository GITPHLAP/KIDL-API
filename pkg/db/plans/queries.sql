-- name: InsertPlan :exec 
INSERT INTO public.plans ("year", category_id, january_amount, february_amount
    , march_amount, april_amount, may_amount, june_amount, july_amount, august_amount
    , september_amount, october_amount, november_amount, december_amount, "type") 
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);

-- name: GetPlansByYear :many
SELECT *
FROM plans
WHERE year = $1;

-- name: UpdatePlanById :exec
UPDATE plans
SET "year"=$1
    , category_id=$2
    , january_amount=$3
    , february_amount=$4
    , march_amount=$5
    , april_amount=$6
    , may_amount=$7
    , june_amount=$8
    , july_amount=$9
    , august_amount=$10
    , september_amount=$11
    , october_amount=$12
    , november_amount=$13
    , december_amount=$14
    , "type"=$15
WHERE id=$16;