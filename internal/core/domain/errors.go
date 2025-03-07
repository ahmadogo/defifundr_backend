package domain

import "errors"

// Domain-level errors
var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrWalletAddressExists  = errors.New("wallet address already exists")
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrPayrollNotFound      = errors.New("payroll not found")
	ErrInvoiceNotFound      = errors.New("invoice not found")
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrUnauthorized         = errors.New("unauthorized access")
	ErrForbidden            = errors.New("forbidden access")
	ErrInvalidBlockchain    = errors.New("invalid blockchain type")
	ErrContractDeployment   = errors.New("error deploying smart contract")
	ErrContractExecution    = errors.New("error executing smart contract")
	ErrInvalidWalletAddress = errors.New("invalid wallet address")
	ErrNoEmployeesInPayroll = errors.New("no employees in payroll")
)
