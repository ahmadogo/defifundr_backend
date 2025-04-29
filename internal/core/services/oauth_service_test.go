package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// MockUserRepository is a mock of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	args := m.Called(ctx, userID, newPassword)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// OAuthServiceTestSuite is a test suite for the OAuth service
type OAuthServiceTestSuite struct {
	suite.Suite
	mockRepo    *MockUserRepository
	config      config.Config
	service     *services.OAuthServiceImpl
	testUser    *domain.User
	testContext context.Context
}

// SetupTest runs before each test
func (s *OAuthServiceTestSuite) SetupTest() {
	s.mockRepo = new(MockUserRepository)
	s.config = config.Config{}
	// s.service = services.NewOAuthService(s.mockRepo, s.config).(*services.OAuthServiceImpl)

	// Create test user
	s.testUser = &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.testContext = context.Background()
}

// TestGetUserInfoByEmail tests the GetUserInfoByEmail method
func (s *OAuthServiceTestSuite) TestGetUserInfoByEmail() {
	// Arrange
	s.mockRepo.On("GetUserByEmail", s.testContext, "test@example.com").Return(s.testUser, nil)

	// Act
	result, err := s.service.GetUserInfoByEmail(s.testContext, "test@example.com")

	// Assert
	require.NoError(s.Suite.T(), err, "GetUserInfoByEmail should not return an error")
	require.NotNil(s.Suite.T(), result, "Result should not be nil")

	assert.Equal(s.Suite.T(), s.testUser.ID, result["user_id"], "User ID should match")
	assert.Equal(s.Suite.T(), s.testUser.Email, result["email"], "Email should match")
	assert.Equal(s.Suite.T(), s.testUser.FirstName+" "+s.testUser.LastName, result["name"], "Name should match")
	assert.NotNil(s.Suite.T(), result["created_at"], "Created at should be set")
	assert.NotNil(s.Suite.T(), result["updated_at"], "Updated at should be set")

	// Verify mock expectations
	s.mockRepo.AssertExpectations(s.Suite.T())
}

// TestGetUserInfoByEmail_UserNotFound tests the error case when user is not found
func (s *OAuthServiceTestSuite) TestGetUserInfoByEmail_UserNotFound() {
	// Arrange
	notFoundErr := errors.New("user not found")
	s.mockRepo.On("GetUserByEmail", s.testContext, "nonexistent@example.com").Return(nil, notFoundErr)

	// Act
	result, err := s.service.GetUserInfoByEmail(s.testContext, "nonexistent@example.com")

	// Assert
	assert.Error(s.Suite.T(), err, "Should return an error when user not found")
	assert.Nil(s.Suite.T(), result, "Result should be nil when error occurs")
	assert.Contains(s.Suite.T(), err.Error(), "failed to get user info", "Error message should indicate user retrieval failure")

	// Verify mock expectations
	s.mockRepo.AssertExpectations(s.Suite.T())
}

// TestGetUserInfo tests the GetUserInfo method with a valid token
func (s *OAuthServiceTestSuite) TestGetUserInfo() {
	// Arrange
	// Mock ValidateWebAuthToken to return mock claims
	mockClaims := map[string]interface{}{
		"email": "test@example.com",
		"name":  "Web3Auth User",
	}
	s.service.ValidateWebAuthTokenFunc = func(ctx context.Context, tokenString string) (map[string]interface{}, error) {
		return mockClaims, nil
	}

	// Mock user repository
	s.mockRepo.On("GetUserByEmail", s.testContext, "test@example.com").Return(s.testUser, nil)

	// Act
	result, err := s.service.GetUserInfo(s.testContext, "valid.test.token")

	// Assert
	require.NoError(s.Suite.T(), err, "GetUserInfo should not return an error with valid token")
	require.NotNil(s.Suite.T(), result, "Result should not be nil")

	assert.Equal(s.Suite.T(), s.testUser.ID, result["user_id"], "User ID should match")
	assert.Equal(s.Suite.T(), "test@example.com", result["email"], "Email should match")
	assert.Equal(s.Suite.T(), mockClaims["name"], result["name"], "Name should come from token claims")

	// Verify mock expectations
	s.mockRepo.AssertExpectations(s.Suite.T())
}

// TestGetUserInfo_InvalidToken tests the GetUserInfo method with an invalid token
func (s *OAuthServiceTestSuite) TestGetUserInfo_InvalidToken() {
	// Arrange
	tokenErr := errors.New("invalid token")
	s.service.ValidateWebAuthTokenFunc = func(ctx context.Context, tokenString string) (map[string]interface{}, error) {
		return nil, tokenErr
	}

	// Act
	result, err := s.service.GetUserInfo(s.testContext, "invalid.token")

	// Assert
	assert.Error(s.Suite.T(), err, "Should return an error with invalid token")
	assert.Nil(s.Suite.T(), result, "Result should be nil when token is invalid")
	assert.Equal(s.Suite.T(), tokenErr, err, "Error should be passed through from token validation")
}

// TestGetUserInfo_MissingEmail tests when email is missing from token claims
func (s *OAuthServiceTestSuite) TestGetUserInfo_MissingEmail() {
	// Arrange
	// Mock ValidateWebAuthToken to return claims without email
	mockClaims := map[string]interface{}{
		"name": "Web3Auth User",
		// No email field
	}
	s.service.ValidateWebAuthTokenFunc = func(ctx context.Context, tokenString string) (map[string]interface{}, error) {
		return mockClaims, nil
	}

	// Act
	result, err := s.service.GetUserInfo(s.testContext, "missing-email.token")

	// Assert
	assert.Error(s.Suite.T(), err, "Should return an error when email is missing")
	assert.Nil(s.Suite.T(), result, "Result should be nil when email is missing")
	assert.Contains(s.Suite.T(), err.Error(), "email not found", "Error message should indicate missing email")
}

// TestValidateWebAuthToken tests the token validation logic
func (s *OAuthServiceTestSuite) TestValidateWebAuthToken_InvalidToken() {
	// This test would normally require mocking the JWT library,
	// but we can still test the invalid token case

	// Act - use the real implementation (not the mock function)
	s.service.ValidateWebAuthTokenFunc = nil
	result, err := s.service.ValidateWebAuthToken(s.testContext, "invalid.token.format")

	// Assert
	assert.Error(s.Suite.T(), err, "Should return an error for invalid token")
	assert.Nil(s.Suite.T(), result, "Result should be nil for invalid token")
}

// Run the test suite
func TestOAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(OAuthServiceTestSuite))
}

// Individual test case without the suite
func TestOAuthService_IndividualTest(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockUserRepository)

	// Create a test user
	testUser := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(testUser, nil)

	// Create the service with the mock repository
	cfg := config.Config{}
	service := services.NewOAuthService(mockRepo, cfg)

	// Call the method being tested
	result, err := service.GetUserInfo(context.Background(), "test@example.com")

	// Assert expectations
	require.NoError(t, err, "GetUserInfoByEmail should not return an error")
	require.NotNil(t, result, "Result should not be nil")

	assert.Equal(t, testUser.ID, result["user_id"], "User ID should match")
	assert.Equal(t, testUser.Email, result["email"], "Email should match")
	assert.Equal(t, testUser.FirstName+" "+testUser.LastName, result["name"], "Name should match")

	// Verify mocks
	mockRepo.AssertExpectations(t)
}
