package models

import "github.com/jackc/pgx/v5/pgtype"

// region General
type MessageCodeData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseData struct {
	Meta MessageCodeData `json:"meta"`
	Data interface{}     `json:"data"`
}
type ResponseError struct {
	Meta MessageCodeError `json:"meta"`
}

type MessageCodeError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

// end region

type HealthCheckReturn struct {
	Version string `json:"version"`
}

// region Overview

type TypeData struct {
	ActualSum pgtype.Float8 `json:"actual_sum"`
	Category  pgtype.Text   `json:"category"`
	PlanSum   pgtype.Float8 `json:"plan_sum"`
}
type TypeTotalData struct {
	TotalActual pgtype.Float8 `json:"total_actual,omitempty"`
	TotalPlan   pgtype.Float8 `json:"total_plan,omitempty"`
}

type MonthData struct {
	Expenses      []TypeData     `json:"expenses,omitempty"`
	ExpensesTotal *TypeTotalData `json:"expenses_total,omitempty"`
	Income        []TypeData     `json:"income,omitempty"`
	IncomeTotal   *TypeTotalData `json:"income_total,omitempty"`
}
type MonthWrapper map[string][]MonthData
type OverviewReturn []struct {
	Year   int            `json:"year"`
	Months []MonthWrapper `json:"months"`
}

// end region
