package models

type HealthCheckParams struct {
}

type GetOverviewsParams struct {
	Year   int     `path:"year"`
	Months []int32 `query:"months"`
}
