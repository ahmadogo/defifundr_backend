package sqlc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

// SeedDB populates the database with sample data for development and testing
func SeedDB(ctx context.Context, q *Queries) error {
	// Clear existing data if needed
	if err := clearTables(ctx, q); err != nil {
		return fmt.Errorf("error clearing tables: %w", err)
	}

	// Create sample users
	users, err := createSampleUsers(ctx, q)
	if err != nil {
		return fmt.Errorf("error creating sample users: %w", err)
	}

	// Create wallets for users
	if err := createSampleWallets(ctx, q, users); err != nil {
		return fmt.Errorf("error creating sample wallets: %w", err)
	}

	// Create organizations
	orgs, err := createSampleOrganizations(ctx, q, users)
	if err != nil {
		return fmt.Errorf("error creating sample organizations: %w", err)
	}

	// Add members to organizations
	if err := addSampleOrganizationMembers(ctx, q, orgs, users); err != nil {
		return fmt.Errorf("error adding sample org members: %w", err)
	}

	// Create payrolls
	payrolls, err := createSamplePayrolls(ctx, q, users, orgs)
	if err != nil {
		return fmt.Errorf("error creating sample payrolls: %w", err)
	}

	// Add employees to payrolls
	if err := addSamplePayrollEmployees(ctx, q, payrolls, users); err != nil {
		return fmt.Errorf("error adding sample payroll employees: %w", err)
	}

	// Create invoices
	if err := createSampleInvoices(ctx, q, users); err != nil {
		return fmt.Errorf("error creating sample invoices: %w", err)
	}

	// Create sample transactions
	if err := createSampleTransactions(ctx, q, users); err != nil {
		return fmt.Errorf("error creating sample transactions: %w", err)
	}

	// Create sample notifications
	if err := createSampleNotifications(ctx, q, users); err != nil {
		return fmt.Errorf("error creating sample notifications: %w", err)
	}

	log.Println("Database seeded successfully!")
	return nil
}

// Helper functions for creating sample data

func clearTables(ctx context.Context, q *Queries) error {
	// Clear tables in reverse order of dependencies
	tables := []string{
		"notifications",
		"transactions",
		"payroll_employees",
		"payrolls",
		"invoices",
		"organization_members",
		"organizations",
		"wallets",
		"sessions",
		"kyc",
		"users",
	}

	for _, table := range tables {
		if _, err := q.db.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return err
		}
	}
	return nil
}

func createSampleUsers(ctx context.Context, q *Queries) (map[string]uuid.UUID, error) {
	users := make(map[string]uuid.UUID)

	// Create a variety of user types

	// Employers
	employerParams := []CreateUserParams{
		{
			Email:               "employer1@example.com",
			PasswordHash:        hashPassword("password123"),
			AccountType:         "business",
			PersonalAccountType: "business",
			FirstName:           "John",
			LastName:            "Smith",
			Nationality:         "United States",
			ResidencialCountry:  pgtype.Text{String: "United States", Valid: true},
			CompanyWebsite:      pgtype.Text{String: "techcorp.com", Valid: true},
		},
		{
			Email:               "employer2@example.com",
			PasswordHash:        hashPassword("password123"),
			AccountType:         "business",
			PersonalAccountType: "business",
			FirstName:           "Sarah",
			LastName:            "Johnson",
			Nationality:         "Canada",
			ResidencialCountry:  pgtype.Text{String: "Canada", Valid: true},
			CompanyWebsite:      pgtype.Text{String: "innovatech.ca", Valid: true},
		},
	}

	for _, params := range employerParams {
		user, err := q.CreateUser(ctx, params)
		if err != nil {
			return nil, err
		}

		if params.Email == "employer1@example.com" {
			users["employer1"] = user.ID
		} else {
			users["employer2"] = user.ID
		}

		// Create KYC records
		_, err = q.CreateKYC(ctx, CreateKYCParams{
			UserID:               user.ID,
			FaceVerification:     true,
			IdentityVerification: true,
		})
		if err != nil {
			return nil, err
		}
	}

	// Employees
	employeeParams := []CreateUserParams{
		{
			Email:               "employee1@example.com",
			PasswordHash:        hashPassword("password123"),
			AccountType:         "personal",
			PersonalAccountType: "contractor",
			FirstName:           "Michael",
			LastName:            "Brown",
			Nationality:         "United States",
			ResidencialCountry:  pgtype.Text{String: "United States", Valid: true},
			JobRole:             pgtype.Text{String: "Software Developer", Valid: true},
			EmploymentType:      pgtype.Text{String: "Full-time", Valid: true},
		},
		{
			Email:               "employee2@example.com",
			PasswordHash:        hashPassword("password123"),
			AccountType:         "personal",
			PersonalAccountType: "contractor",
			FirstName:           "Emma",
			LastName:            "Davis",
			Nationality:         "United Kingdom",
			ResidencialCountry:  pgtype.Text{String: "United Kingdom", Valid: true},
			JobRole:             pgtype.Text{String: "UI/UX Designer", Valid: true},
			EmploymentType:      pgtype.Text{String: "Full-time", Valid: true},
		},
		{
			Email:               "employee3@example.com",
			PasswordHash:        hashPassword("password123"),
			AccountType:         "personal",
			PersonalAccountType: "contractor",
			FirstName:           "Carlos",
			LastName:            "Rodriguez",
			Nationality:         "Spain",
			ResidencialCountry:  pgtype.Text{String: "Spain", Valid: true},
			JobRole:             pgtype.Text{String: "Project Manager", Valid: true},
			EmploymentType:      pgtype.Text{String: "Full-time", Valid: true},
		},
	}

	for i, params := range employeeParams {
		user, err := q.CreateUser(ctx, params)
		if err != nil {
			return nil, err
		}

		users[fmt.Sprintf("employee%d", i+1)] = user.ID

		// Create KYC records
		_, err = q.CreateKYC(ctx, CreateKYCParams{
			UserID:               user.ID,
			FaceVerification:     true,
			IdentityVerification: true,
		})
		if err != nil {
			return nil, err
		}
	}

	// Freelancers
	freelancerParams := []CreateUserParams{
		{
			Email:               "freelancer1@example.com",
			PasswordHash:        hashPassword("password123"),
			AccountType:         "personal",
			PersonalAccountType: "contractor",
			FirstName:           "Jessica",
			LastName:            "Lee",
			Nationality:         "Singapore",
			ResidencialCountry:  pgtype.Text{String: "Singapore", Valid: true},
			JobRole:             pgtype.Text{String: "Content Writer", Valid: true},
			EmploymentType:      pgtype.Text{String: "Freelance", Valid: true},
		},
		{
			Email:               "freelancer2@example.com",
			PasswordHash:        hashPassword("password123"),
			AccountType:         "personal",
			PersonalAccountType: "contractor",
			FirstName:           "David",
			LastName:            "Kim",
			Nationality:         "South Korea",
			ResidencialCountry:  pgtype.Text{String: "South Korea", Valid: true},
			JobRole:             pgtype.Text{String: "Mobile Developer", Valid: true},
			EmploymentType:      pgtype.Text{String: "Freelance", Valid: true},
		},
	}

	for i, params := range freelancerParams {
		user, err := q.CreateUser(ctx, params)
		if err != nil {
			return nil, err
		}

		users[fmt.Sprintf("freelancer%d", i+1)] = user.ID

		// Create KYC records
		_, err = q.CreateKYC(ctx, CreateKYCParams{
			UserID:               user.ID,
			FaceVerification:     true,
			IdentityVerification: true,
		})
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

func hashPassword(password string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedBytes)
}

func createSampleWallets(ctx context.Context, q *Queries, users map[string]uuid.UUID) error {
	wallets := []struct {
		UserKey       string
		WalletAddress string
		Chain         string
		IsPrimary     bool
	}{
		{
			UserKey:       "employer1",
			WalletAddress: "0x1234567890123456789012345678901234567890",
			Chain:         "ethereum",
			IsPrimary:     true,
		},
		{
			UserKey:       "employer1",
			WalletAddress: "9XyJa3o1SP5vU5Vr6EcBR9LS7oJ2aNabWHqpgAqiwLi4",
			Chain:         "solana",
			IsPrimary:     false,
		},
		{
			UserKey:       "employer2",
			WalletAddress: "0x2345678901234567890123456789012345678901",
			Chain:         "ethereum",
			IsPrimary:     true,
		},
		{
			UserKey:       "employee1",
			WalletAddress: "0x3456789012345678901234567890123456789012",
			Chain:         "ethereum",
			IsPrimary:     true,
		},
		{
			UserKey:       "employee2",
			WalletAddress: "0x4567890123456789012345678901234567890123",
			Chain:         "ethereum",
			IsPrimary:     true,
		},
		{
			UserKey:       "employee3",
			WalletAddress: "0x5678901234567890123456789012345678901234",
			Chain:         "ethereum",
			IsPrimary:     true,
		},
		{
			UserKey:       "freelancer1",
			WalletAddress: "0x6789012345678901234567890123456789012345",
			Chain:         "ethereum",
			IsPrimary:     true,
		},
		{
			UserKey:       "freelancer1",
			WalletAddress: "8FLJDzPMU4Az3E5Wj2JzxUBSgDqDpwzxiyV9eKn7GR6S",
			Chain:         "solana",
			IsPrimary:     false,
		},
		{
			UserKey:       "freelancer2",
			WalletAddress: "0x7890123456789012345678901234567890123456",
			Chain:         "ethereum",
			IsPrimary:     true,
		},
	}

	for _, w := range wallets {
		// PIN hash would be handled properly in a real application
		pinHash, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)

		_, err := q.CreateWallet(ctx, CreateWalletParams{
			UserID:        users[w.UserKey],
			WalletAddress: w.WalletAddress,
			Chain:         w.Chain,
			IsPrimary:     w.IsPrimary,
			PinHash:       string(pinHash),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func createSampleOrganizations(ctx context.Context, q *Queries, users map[string]uuid.UUID) (map[string]uuid.UUID, error) {
	orgs := make(map[string]uuid.UUID)

	organizations := []struct {
		Name        string
		EmployerKey string
	}{
		{
			Name:        "TechCorp Inc.",
			EmployerKey: "employer1",
		},
		{
			Name:        "InnovaTech Solutions",
			EmployerKey: "employer2",
		},
	}

	for i, org := range organizations {
		newOrg, err := q.CreateOrganization(ctx, CreateOrganizationParams{
			Name:       org.Name,
			EmployerID: users[org.EmployerKey],
		})
		if err != nil {
			return nil, err
		}

		orgs[fmt.Sprintf("org%d", i+1)] = newOrg.ID
	}

	return orgs, nil
}

func addSampleOrganizationMembers(ctx context.Context, q *Queries, orgs map[string]uuid.UUID, users map[string]uuid.UUID) error {
	members := []struct {
		OrgKey      string
		EmployeeKey string
		Role        string
	}{
		{
			OrgKey:      "org1",
			EmployeeKey: "employee1",
			Role:        "developer",
		},
		{
			OrgKey:      "org1",
			EmployeeKey: "employee2",
			Role:        "designer",
		},
		{
			OrgKey:      "org2",
			EmployeeKey: "employee3",
			Role:        "manager",
		},
	}

	for _, member := range members {
		_, err := q.AddOrganizationMember(ctx, AddOrganizationMemberParams{
			OrganizationID: orgs[member.OrgKey],
			EmployeeID:     users[member.EmployeeKey],
			Role:           member.Role,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func createSamplePayrolls(ctx context.Context, q *Queries, users map[string]uuid.UUID, orgs map[string]uuid.UUID) (map[string]uuid.UUID, error) {
	payrolls := make(map[string]uuid.UUID)

	payrollData := []struct {
		EmployerKey      string
		OrgKey           string
		PaymentFrequency string
		SalaryAmount     decimal.Decimal
		Currency         string
		Status           string
	}{
		{
			EmployerKey:      "employer1",
			OrgKey:           "org1",
			PaymentFrequency: "monthly",
			SalaryAmount:     decimal.NewFromFloat(5000.00),
			Currency:         "USDC",
			Status:           "active",
		},
		{
			EmployerKey:      "employer2",
			OrgKey:           "org2",
			PaymentFrequency: "bi-weekly",
			SalaryAmount:     decimal.NewFromFloat(2500.00),
			Currency:         "ETH",
			Status:           "active",
		},
	}

	for i, p := range payrollData {
		orgID := pgtype.UUID{Valid: false}

		payroll, err := q.CreatePayroll(ctx, CreatePayrollParams{
			EmployerID:       users[p.EmployerKey],
			OrganizationID:   orgID,
			PaymentFrequency: p.PaymentFrequency,
			SalaryAmount:     p.SalaryAmount,
			Currency:         p.Currency,
			Status:           p.Status,
		})
		if err != nil {
			return nil, err
		}

		payrolls[fmt.Sprintf("payroll%d", i+1)] = payroll.ID
	}

	return payrolls, nil
}

func addSamplePayrollEmployees(ctx context.Context, q *Queries, payrolls map[string]uuid.UUID, users map[string]uuid.UUID) error {
	assignments := []struct {
		PayrollKey  string
		EmployeeKey string
	}{
		{
			PayrollKey:  "payroll1",
			EmployeeKey: "employee1",
		},
		{
			PayrollKey:  "payroll1",
			EmployeeKey: "employee2",
		},
		{
			PayrollKey:  "payroll2",
			EmployeeKey: "employee3",
		},
	}

	for _, a := range assignments {
		_, err := q.AddPayrollEmployee(ctx, AddPayrollEmployeeParams{
			PayrollID:  payrolls[a.PayrollKey],
			EmployeeID: users[a.EmployeeKey],
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func createSampleInvoices(ctx context.Context, q *Queries, users map[string]uuid.UUID) error {
	invoices := []struct {
		FreelancerKey string
		EmployerKey   string
		Amount        decimal.Decimal
		Currency      string
		Status        string
	}{
		{
			FreelancerKey: "freelancer1",
			EmployerKey:   "employer1",
			Amount:        decimal.NewFromFloat(1200.00),
			Currency:      "USDC",
			Status:        "pending",
		},
		{
			FreelancerKey: "freelancer1",
			EmployerKey:   "employer2",
			Amount:        decimal.NewFromFloat(800.00),
			Currency:      "USDC",
			Status:        "approved",
		},
		{
			FreelancerKey: "freelancer2",
			EmployerKey:   "employer1",
			Amount:        decimal.NewFromFloat(2500.00),
			Currency:      "ETH",
			Status:        "paid",
		},
	}

	for _, inv := range invoices {
		_, err := q.CreateInvoice(ctx, CreateInvoiceParams{
			FreelancerID: users[inv.FreelancerKey],
			EmployerID:   users[inv.EmployerKey],
			Amount:       inv.Amount,
			Currency:     inv.Currency,
			Status:       inv.Status,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func createSampleTransactions(ctx context.Context, q *Queries, users map[string]uuid.UUID) error {
	transactions := []struct {
		UserKey  string
		TxHash   string
		Amount   decimal.Decimal
		Currency string
		Type     string
		Status   string
	}{
		{
			UserKey:  "employee1",
			TxHash:   "0xabcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			Amount:   decimal.NewFromFloat(5000.00),
			Currency: "USDC",
			Type:     "payroll",
			Status:   "success",
		},
		{
			UserKey:  "employee2",
			TxHash:   "0xbcde2345bcde2345bcde2345bcde2345bcde2345bcde2345bcde2345bcde2345",
			Amount:   decimal.NewFromFloat(5000.00),
			Currency: "USDC",
			Type:     "payroll",
			Status:   "success",
		},
		{
			UserKey:  "employee3",
			TxHash:   "0xcdef3456cdef3456cdef3456cdef3456cdef3456cdef3456cdef3456cdef3456",
			Amount:   decimal.NewFromFloat(2500.00),
			Currency: "ETH",
			Type:     "payroll",
			Status:   "success",
		},
		{
			UserKey:  "freelancer2",
			TxHash:   "0xefgh5678efgh5678efgh5678efgh5678efgh5678efgh5678efgh5678efgh5678",
			Amount:   decimal.NewFromFloat(2500.00),
			Currency: "ETH",
			Type:     "invoice",
			Status:   "success",
		},
	}

	for _, tx := range transactions {
		_, err := q.CreateTransaction(ctx, CreateTransactionParams{
			UserID:   users[tx.UserKey],
			TxHash:   tx.TxHash,
			Amount:   tx.Amount,
			Currency: tx.Currency,
			Type:     tx.Type,
			Status:   tx.Status,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func createSampleNotifications(ctx context.Context, q *Queries, users map[string]uuid.UUID) error {
	notifications := []struct {
		UserKey string
		Message string
		Type    string
	}{
		{
			UserKey: "employee1",
			Message: "Your monthly payroll has been processed",
			Type:    "payroll",
		},
		{
			UserKey: "employee2",
			Message: "Your monthly payroll has been processed",
			Type:    "payroll",
		},
		{
			UserKey: "employee3",
			Message: "Your bi-weekly payroll has been processed",
			Type:    "payroll",
		},
		{
			UserKey: "freelancer1",
			Message: "Your invoice has been approved",
			Type:    "invoice",
		},
		{
			UserKey: "freelancer2",
			Message: "Your invoice has been paid",
			Type:    "invoice",
		},
		{
			UserKey: "employer1",
			Message: "New invoice received from Jessica Lee",
			Type:    "invoice",
		},
		{
			UserKey: "employer2",
			Message: "New invoice received from David Kim",
			Type:    "invoice",
		},
	}

	for _, notif := range notifications {
		_, err := q.CreateNotification(ctx, CreateNotificationParams{
			UserID:  users[notif.UserKey],
			Message: notif.Message,
			Type:    notif.Type,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
