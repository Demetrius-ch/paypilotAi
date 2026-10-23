package models

type Summary struct {
	TotalRevenue     float64 `json:"total_revenue"`
	OrderCount       int     `json:"order_count"`
	PaidCount        int     `json:"paid_count"`
	PendingCount     int     `json:"pending_count"`
	PendingAmount    float64 `json:"pending_amount"`
	PreviousPeriod   float64 `json:"previous_period"`
	PercentageChange float64 `json:"percentage_change"`
}

type RevenueDataPoint struct {
	Period  string  `json:"period"`
	Revenue float64 `json:"revenue"`
}

type AnalyticsFilter struct {
	Period string // today, this_week, this_month, last_month, this_year
}
