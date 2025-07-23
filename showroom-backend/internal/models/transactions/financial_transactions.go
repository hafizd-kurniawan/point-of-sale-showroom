package transactions

import (
	"time"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// FinancialTransactionType represents the type of financial transaction
type FinancialTransactionType string

const (
	FinancialTransactionVehiclePurchase FinancialTransactionType = "vehicle_purchase"
	FinancialTransactionVehicleSale     FinancialTransactionType = "vehicle_sale"
	FinancialTransactionPartsPurchase   FinancialTransactionType = "parts_purchase"
	FinancialTransactionServiceRevenue  FinancialTransactionType = "service_revenue"
	FinancialTransactionExpense         FinancialTransactionType = "expense"
	FinancialTransactionCashAdjustment  FinancialTransactionType = "cash_adjustment"
	FinancialTransactionCreditPayment   FinancialTransactionType = "credit_payment"
	FinancialTransactionSupplierPayment FinancialTransactionType = "supplier_payment"
	FinancialTransactionWarrantyClaim   FinancialTransactionType = "warranty_claim"
	FinancialTransactionInsuranceClaim  FinancialTransactionType = "insurance_claim"
)

// FinancialTransactionCategory represents the category of financial transaction
type FinancialTransactionCategory string

const (
	FinancialCategoryIncome     FinancialTransactionCategory = "income"
	FinancialCategoryExpense    FinancialTransactionCategory = "expense"
	FinancialCategoryAsset      FinancialTransactionCategory = "asset"
	FinancialCategoryLiability  FinancialTransactionCategory = "liability"
	FinancialCategoryEquity     FinancialTransactionCategory = "equity"
)

// AccountType represents the type of account for financial transactions
type AccountType string

const (
	AccountTypeCash              AccountType = "cash"
	AccountTypeBank              AccountType = "bank"
	AccountTypeAccountsReceivable AccountType = "accounts_receivable"
	AccountTypeAccountsPayable   AccountType = "accounts_payable"
	AccountTypeInventory         AccountType = "inventory"
	AccountTypeFixedAsset        AccountType = "fixed_asset"
	AccountTypeRevenue           AccountType = "revenue"
	AccountTypeCOGS              AccountType = "cogs"
	AccountTypeOperatingExpense  AccountType = "operating_expense"
)

// FinancialTransactionStatus represents the status of a financial transaction
type FinancialTransactionStatus string

const (
	FinancialStatusDraft     FinancialTransactionStatus = "draft"
	FinancialStatusPending   FinancialTransactionStatus = "pending"
	FinancialStatusApproved  FinancialTransactionStatus = "approved"
	FinancialStatusPosted    FinancialTransactionStatus = "posted"
	FinancialStatusCancelled FinancialTransactionStatus = "cancelled"
	FinancialStatusReversed  FinancialTransactionStatus = "reversed"
)

// FinancialTransaction represents a financial transaction in the system
type FinancialTransaction struct {
	TransactionID            int                         `json:"transaction_id" db:"transaction_id"`
	TransactionNumber        string                      `json:"transaction_number" db:"transaction_number"`
	TransactionType          FinancialTransactionType    `json:"transaction_type" db:"transaction_type"`
	TransactionCategory      FinancialTransactionCategory `json:"transaction_category" db:"transaction_category"`
	AccountType              AccountType                 `json:"account_type" db:"account_type"`
	AccountCode              string                      `json:"account_code" db:"account_code"`
	AccountName              string                      `json:"account_name" db:"account_name"`
	DebitAmount              float64                     `json:"debit_amount" db:"debit_amount"`
	CreditAmount             float64                     `json:"credit_amount" db:"credit_amount"`
	NetAmount                float64                     `json:"net_amount" db:"net_amount"`
	PaymentMethod            *PaymentMethod              `json:"payment_method,omitempty" db:"payment_method"`
	PaymentReference         *string                     `json:"payment_reference,omitempty" db:"payment_reference"`
	TransactionDate          time.Time                   `json:"transaction_date" db:"transaction_date"`
	PostingDate              *time.Time                  `json:"posting_date,omitempty" db:"posting_date"`
	RelatedVehicleID         *int                        `json:"related_vehicle_id,omitempty" db:"related_vehicle_id"`
	RelatedCustomerID        *int                        `json:"related_customer_id,omitempty" db:"related_customer_id"`
	RelatedSupplierID        *int                        `json:"related_supplier_id,omitempty" db:"related_supplier_id"`
	RelatedPurchaseOrderID   *int                        `json:"related_purchase_order_id,omitempty" db:"related_purchase_order_id"`
	RelatedSalesID           *int                        `json:"related_sales_id,omitempty" db:"related_sales_id"`
	RelatedServiceID         *int                        `json:"related_service_id,omitempty" db:"related_service_id"`
	Description              string                      `json:"description" db:"description"`
	TransactionMemo          *string                     `json:"transaction_memo,omitempty" db:"transaction_memo"`
	ProcessedBy              int                         `json:"processed_by" db:"processed_by"`
	ApprovedBy               *int                        `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt               *time.Time                  `json:"approved_at,omitempty" db:"approved_at"`
	Status                   FinancialTransactionStatus  `json:"status" db:"status"`
	SupportingDocumentsJSON  *string                     `json:"supporting_documents_json,omitempty" db:"supporting_documents_json"`
	CreatedAt                time.Time                   `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time                   `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time                  `json:"deleted_at,omitempty" db:"deleted_at"`

	// Related data
	ProcessedByName string `json:"processed_by_name,omitempty" db:"processed_by_name"`
	ApprovedByName  string `json:"approved_by_name,omitempty" db:"approved_by_name"`
	CustomerName    string `json:"customer_name,omitempty" db:"customer_name"`
	SupplierName    string `json:"supplier_name,omitempty" db:"supplier_name"`
}

// FinancialTransactionCreateRequest represents a request to create a financial transaction
type FinancialTransactionCreateRequest struct {
	TransactionType          FinancialTransactionType    `json:"transaction_type" binding:"required"`
	TransactionCategory      FinancialTransactionCategory `json:"transaction_category" binding:"required"`
	AccountType              AccountType                 `json:"account_type" binding:"required"`
	AccountCode              string                      `json:"account_code" binding:"required,max=20"`
	AccountName              string                      `json:"account_name" binding:"required,max=100"`
	DebitAmount              float64                     `json:"debit_amount" binding:"min=0"`
	CreditAmount             float64                     `json:"credit_amount" binding:"min=0"`
	PaymentMethod            *PaymentMethod              `json:"payment_method,omitempty"`
	PaymentReference         *string                     `json:"payment_reference,omitempty" binding:"omitempty,max=100"`
	TransactionDate          time.Time                   `json:"transaction_date" binding:"required"`
	RelatedVehicleID         *int                        `json:"related_vehicle_id,omitempty" binding:"omitempty,min=1"`
	RelatedCustomerID        *int                        `json:"related_customer_id,omitempty" binding:"omitempty,min=1"`
	RelatedSupplierID        *int                        `json:"related_supplier_id,omitempty" binding:"omitempty,min=1"`
	RelatedPurchaseOrderID   *int                        `json:"related_purchase_order_id,omitempty" binding:"omitempty,min=1"`
	RelatedSalesID           *int                        `json:"related_sales_id,omitempty" binding:"omitempty,min=1"`
	RelatedServiceID         *int                        `json:"related_service_id,omitempty" binding:"omitempty,min=1"`
	Description              string                      `json:"description" binding:"required,max=500"`
	TransactionMemo          *string                     `json:"transaction_memo,omitempty"`
}

// CashFlowDaily represents daily cash flow summary
type CashFlowDaily struct {
	CashflowID               int        `json:"cashflow_id" db:"cashflow_id"`
	TransactionDate          time.Time  `json:"transaction_date" db:"transaction_date"`
	OpeningBalanceCash       float64    `json:"opening_balance_cash" db:"opening_balance_cash"`
	OpeningBalanceBank       float64    `json:"opening_balance_bank" db:"opening_balance_bank"`
	TotalCashIn              float64    `json:"total_cash_in" db:"total_cash_in"`
	TotalCashOut             float64    `json:"total_cash_out" db:"total_cash_out"`
	TotalBankIn              float64    `json:"total_bank_in" db:"total_bank_in"`
	TotalBankOut             float64    `json:"total_bank_out" db:"total_bank_out"`
	ClosingBalanceCash       float64    `json:"closing_balance_cash" db:"closing_balance_cash"`
	ClosingBalanceBank       float64    `json:"closing_balance_bank" db:"closing_balance_bank"`
	TotalClosingBalance      float64    `json:"total_closing_balance" db:"total_closing_balance"`
	VarianceAmount           float64    `json:"variance_amount" db:"variance_amount"`
	ReconciliationStatus     string     `json:"reconciliation_status" db:"reconciliation_status"`
	ReconciledBy             *int       `json:"reconciled_by,omitempty" db:"reconciled_by"`
	ReconciledAt             *time.Time `json:"reconciled_at,omitempty" db:"reconciled_at"`
	ReconciliationNotes      *string    `json:"reconciliation_notes,omitempty" db:"reconciliation_notes"`
	VarianceExplanation      *string    `json:"variance_explanation,omitempty" db:"variance_explanation"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`

	// Related data
	ReconciledByName string `json:"reconciled_by_name,omitempty" db:"reconciled_by_name"`
}

// BusinessExpenseCategory represents the category of business expense
type BusinessExpenseCategory string

const (
	ExpenseCategorySalary           BusinessExpenseCategory = "salary"
	ExpenseCategoryUtilities        BusinessExpenseCategory = "utilities"
	ExpenseCategoryRent             BusinessExpenseCategory = "rent"
	ExpenseCategoryInsurance        BusinessExpenseCategory = "insurance"
	ExpenseCategoryMarketing        BusinessExpenseCategory = "marketing"
	ExpenseCategoryMaintenance      BusinessExpenseCategory = "maintenance"
	ExpenseCategoryFuel             BusinessExpenseCategory = "fuel"
	ExpenseCategoryOfficeSupplies   BusinessExpenseCategory = "office_supplies"
	ExpenseCategoryProfessionalFees BusinessExpenseCategory = "professional_fees"
	ExpenseCategoryTaxes            BusinessExpenseCategory = "taxes"
	ExpenseCategoryDepreciation     BusinessExpenseCategory = "depreciation"
	ExpenseCategoryOther            BusinessExpenseCategory = "other"
)

// BusinessExpenseStatus represents the status of business expense
type BusinessExpenseStatus string

const (
	ExpenseStatusRequested BusinessExpenseStatus = "requested"
	ExpenseStatusApproved  BusinessExpenseStatus = "approved"
	ExpenseStatusRejected  BusinessExpenseStatus = "rejected"
	ExpenseStatusPaid      BusinessExpenseStatus = "paid"
)

// BusinessExpense represents a business expense
type BusinessExpense struct {
	ExpenseID                  int                     `json:"expense_id" db:"expense_id"`
	ExpenseNumber              string                  `json:"expense_number" db:"expense_number"`
	ExpenseCategory            BusinessExpenseCategory `json:"expense_category" db:"expense_category"`
	ExpenseSubcategory         *string                 `json:"expense_subcategory,omitempty" db:"expense_subcategory"`
	Amount                     float64                 `json:"amount" db:"amount"`
	ExpenseDate                time.Time               `json:"expense_date" db:"expense_date"`
	PaymentDate                *time.Time              `json:"payment_date,omitempty" db:"payment_date"`
	Description                string                  `json:"description" db:"description"`
	ReceiptNumber              *string                 `json:"receipt_number,omitempty" db:"receipt_number"`
	VendorName                 *string                 `json:"vendor_name,omitempty" db:"vendor_name"`
	PaymentMethod              PaymentMethod           `json:"payment_method" db:"payment_method"`
	PaymentReference           *string                 `json:"payment_reference,omitempty" db:"payment_reference"`
	RequestedBy                int                     `json:"requested_by" db:"requested_by"`
	ApprovedBy                 *int                    `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt                 *time.Time              `json:"approved_at,omitempty" db:"approved_at"`
	TransactionID              *int                    `json:"transaction_id,omitempty" db:"transaction_id"`
	Status                     BusinessExpenseStatus   `json:"status" db:"status"`
	ApprovalNotes              *string                 `json:"approval_notes,omitempty" db:"approval_notes"`
	RejectionReason            *string                 `json:"rejection_reason,omitempty" db:"rejection_reason"`
	SupportingDocumentsJSON    *string                 `json:"supporting_documents_json,omitempty" db:"supporting_documents_json"`
	IsRecurring                bool                    `json:"is_recurring" db:"is_recurring"`
	RecurringFrequencyMonths   *int                    `json:"recurring_frequency_months,omitempty" db:"recurring_frequency_months"`
	NextRecurringDate          *time.Time              `json:"next_recurring_date,omitempty" db:"next_recurring_date"`
	CreatedAt                  time.Time               `json:"created_at" db:"created_at"`
	DeletedAt                  *time.Time              `json:"deleted_at,omitempty" db:"deleted_at"`

	// Related data
	RequestedByName string `json:"requested_by_name,omitempty" db:"requested_by_name"`
	ApprovedByName  string `json:"approved_by_name,omitempty" db:"approved_by_name"`
}

// BusinessExpenseCreateRequest represents a request to create a business expense
type BusinessExpenseCreateRequest struct {
	ExpenseCategory            BusinessExpenseCategory `json:"expense_category" binding:"required"`
	ExpenseSubcategory         *string                 `json:"expense_subcategory,omitempty" binding:"omitempty,max=100"`
	Amount                     float64                 `json:"amount" binding:"required,min=0"`
	ExpenseDate                time.Time               `json:"expense_date" binding:"required"`
	Description                string                  `json:"description" binding:"required,max=500"`
	ReceiptNumber              *string                 `json:"receipt_number,omitempty" binding:"omitempty,max=100"`
	VendorName                 *string                 `json:"vendor_name,omitempty" binding:"omitempty,max=255"`
	PaymentMethod              PaymentMethod           `json:"payment_method" binding:"required"`
	PaymentReference           *string                 `json:"payment_reference,omitempty" binding:"omitempty,max=100"`
	IsRecurring                bool                    `json:"is_recurring"`
	RecurringFrequencyMonths   *int                    `json:"recurring_frequency_months,omitempty" binding:"omitempty,min=1,max=12"`
}

// FinancialFilterParams represents filtering parameters for financial queries
type FinancialFilterParams struct {
	TransactionType     *FinancialTransactionType    `json:"transaction_type,omitempty" form:"transaction_type"`
	TransactionCategory *FinancialTransactionCategory `json:"transaction_category,omitempty" form:"transaction_category"`
	AccountType         *AccountType                 `json:"account_type,omitempty" form:"account_type"`
	Status              *FinancialTransactionStatus  `json:"status,omitempty" form:"status"`
	ExpenseCategory     *BusinessExpenseCategory     `json:"expense_category,omitempty" form:"expense_category"`
	ExpenseStatus       *BusinessExpenseStatus       `json:"expense_status,omitempty" form:"expense_status"`
	DateFrom            *time.Time                   `json:"date_from,omitempty" form:"date_from"`
	DateTo              *time.Time                   `json:"date_to,omitempty" form:"date_to"`
	MinAmount           *float64                     `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount           *float64                     `json:"max_amount,omitempty" form:"max_amount"`
	CustomerID          *int                         `json:"customer_id,omitempty" form:"customer_id"`
	SupplierID          *int                         `json:"supplier_id,omitempty" form:"supplier_id"`
	VehicleID           *int                         `json:"vehicle_id,omitempty" form:"vehicle_id"`
	Search              string                       `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// IsValid validates the financial transaction type
func (f FinancialTransactionType) IsValid() bool {
	switch f {
	case FinancialTransactionVehiclePurchase, FinancialTransactionVehicleSale,
		 FinancialTransactionPartsPurchase, FinancialTransactionServiceRevenue,
		 FinancialTransactionExpense, FinancialTransactionCashAdjustment,
		 FinancialTransactionCreditPayment, FinancialTransactionSupplierPayment,
		 FinancialTransactionWarrantyClaim, FinancialTransactionInsuranceClaim:
		return true
	default:
		return false
	}
}

// IsValid validates the financial transaction category
func (f FinancialTransactionCategory) IsValid() bool {
	switch f {
	case FinancialCategoryIncome, FinancialCategoryExpense, FinancialCategoryAsset,
		 FinancialCategoryLiability, FinancialCategoryEquity:
		return true
	default:
		return false
	}
}

// IsValid validates the account type
func (a AccountType) IsValid() bool {
	switch a {
	case AccountTypeCash, AccountTypeBank, AccountTypeAccountsReceivable,
		 AccountTypeAccountsPayable, AccountTypeInventory, AccountTypeFixedAsset,
		 AccountTypeRevenue, AccountTypeCOGS, AccountTypeOperatingExpense:
		return true
	default:
		return false
	}
}

// IsValid validates the financial transaction status
func (f FinancialTransactionStatus) IsValid() bool {
	switch f {
	case FinancialStatusDraft, FinancialStatusPending, FinancialStatusApproved,
		 FinancialStatusPosted, FinancialStatusCancelled, FinancialStatusReversed:
		return true
	default:
		return false
	}
}

// IsValid validates the business expense category
func (b BusinessExpenseCategory) IsValid() bool {
	switch b {
	case ExpenseCategorySalary, ExpenseCategoryUtilities, ExpenseCategoryRent,
		 ExpenseCategoryInsurance, ExpenseCategoryMarketing, ExpenseCategoryMaintenance,
		 ExpenseCategoryFuel, ExpenseCategoryOfficeSupplies, ExpenseCategoryProfessionalFees,
		 ExpenseCategoryTaxes, ExpenseCategoryDepreciation, ExpenseCategoryOther:
		return true
	default:
		return false
	}
}

// IsValid validates the business expense status
func (b BusinessExpenseStatus) IsValid() bool {
	switch b {
	case ExpenseStatusRequested, ExpenseStatusApproved, ExpenseStatusRejected, ExpenseStatusPaid:
		return true
	default:
		return false
	}
}

// CalculateNetAmount calculates the net amount for a financial transaction
func (f *FinancialTransaction) CalculateNetAmount() {
	f.NetAmount = f.DebitAmount - f.CreditAmount
}

// IsIncome checks if the transaction is income
func (f *FinancialTransaction) IsIncome() bool {
	return f.TransactionCategory == FinancialCategoryIncome
}

// IsExpense checks if the transaction is expense
func (f *FinancialTransaction) IsExpense() bool {
	return f.TransactionCategory == FinancialCategoryExpense
}

// CanBeReversed checks if the transaction can be reversed
func (f *FinancialTransaction) CanBeReversed() bool {
	return f.Status == FinancialStatusPosted
}

// IsApproved checks if the expense is approved
func (b *BusinessExpense) IsApproved() bool {
	return b.Status == ExpenseStatusApproved
}

// CanBePaid checks if the expense can be paid
func (b *BusinessExpense) CanBePaid() bool {
	return b.Status == ExpenseStatusApproved
}

// CalculateClosingBalance calculates the closing balance for daily cash flow
func (c *CashFlowDaily) CalculateClosingBalance() {
	c.ClosingBalanceCash = c.OpeningBalanceCash + c.TotalCashIn - c.TotalCashOut
	c.ClosingBalanceBank = c.OpeningBalanceBank + c.TotalBankIn - c.TotalBankOut
	c.TotalClosingBalance = c.ClosingBalanceCash + c.ClosingBalanceBank
}