package frontend

type UpsertBudgetInput struct {
	Categories []BudgetCategory `json:"categories" binding:"required"`
}

type GetBudgetOutput struct {
	Categories []BudgetCategory `json:"categories" binding:"required"`
}

type BudgetCategory struct {
	PersonalFinanceCategory string  `json:"personal_finance_category" binding:"required"`
	Amount                  float64 `json:"amount" binding:"required"`
}

var VALID_PERSONAL_FINANCE_CATEGORIES = []string{
	"LOAN_PAYMENTS",
	"BANK_FEES",
	"FOOD_AND_DRINK",
	"GENERAL_MERCHANDISE",
	"GENERAL_SERVICES",
	"ENTERTAINMENT",
	"TRAVEL",
	"TRANSPORTATION",
	"RENT_AND_UTILITIES",
	"HOME_IMPROVEMENT",
	"MEDICAL",
	"PERSONAL_CARE",
	"GOVERNMENT_AND_NON_PROFIT",
	"TRANSFER_OUT",
}
