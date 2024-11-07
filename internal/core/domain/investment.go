package domain

// Investment represents an investment made by an investor
type Investment struct {
	InvestorID int
	FarmID     int
	Amount     float64
}
