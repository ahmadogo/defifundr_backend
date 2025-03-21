// internal/adapters/repositories/session_repository.go
package repositories

import (
	"context"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SessionRepository struct {
	store db.Store
}

func NewSessionRepository(store db.Store) *SessionRepository {
	return &SessionRepository{store: store}
}

// CreateSession creates a new session in the database
func (r *SessionRepository) CreateSession(ctx context.Context, session domain.Session) (*domain.Session, error) {
	// Convert optional fields to db nullable types
	webOAuthClientID := pgtype.Text{}
	if session.WebOAuthClientID != nil {
		webOAuthClientID.String = *session.WebOAuthClientID
		webOAuthClientID.Valid = true
	}

	oauthAccessToken := pgtype.Text{}
	if session.OAuthAccessToken != nil {
		oauthAccessToken.String = *session.OAuthAccessToken
		oauthAccessToken.Valid = true
	}

	oauthIDToken := pgtype.Text{}
	if session.OAuthIDToken != nil {
		oauthIDToken.String = *session.OAuthIDToken
		oauthIDToken.Valid = true
	}

	// Call the dbC-generated query
	dbSession, err := r.store.CreateSession(ctx, db.CreateSessionParams{
		ID:               session.ID,
		UserID:           session.UserID,
		RefreshToken:     session.RefreshToken,
		UserAgent:        session.UserAgent,
		WebOauthClientID: webOAuthClientID,
		OauthAccessToken: oauthAccessToken,
		OauthIDToken:     oauthIDToken,
		UserLoginType:    session.UserLoginType,
		MfaEnabled:       session.MFAEnabled,
		ClientIp:         session.ClientIP,
		IsBlocked:        session.IsBlocked,
		ExpiresAt:        session.ExpiresAt,
	})

	if err != nil {
		return nil, err
	}

	// Map the database model back to the domain model
	return mapDbSessionToDomain(dbSession), nil
}

// GetSessionByID retrieves a session by its ID
func (r *SessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	dbSession, err := r.store.GetSessionByID(ctx, id)
	if err != nil {
		if err == pgtype.ErrScanTargetTypeChanged {
			return nil, nil // Session not found
		}
		return nil, err
	}

	return mapDbSessionToDomain(dbSession), nil
}

// GetSessionByRefreshToken retrieves a session by its refresh token
func (r *SessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	dbSession, err := r.store.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		if err == pgtype.ErrScanTargetTypeChanged {
			return nil, nil // Session not found
		}
		return nil, err
	}

	return mapDbSessionToDomain(dbSession), nil
}

// DeleteSession deletes a session from the database
func (r *SessionRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteSessionsByUserID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// BlockSession marks a session as blocked
func (r *SessionRepository) BlockSession(ctx context.Context, id uuid.UUID) error {
	err := r.store.BlockSession(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

// Helper function to map DB model to domain model
func mapDbSessionToDomain(dbSession db.Sessions) *domain.Session {
	session := &domain.Session{
		ID:            dbSession.ID,
		UserID:        dbSession.UserID,
		RefreshToken:  dbSession.RefreshToken,
		UserAgent:     dbSession.UserAgent,
		UserLoginType: dbSession.UserLoginType,
		MFAEnabled:    dbSession.MfaEnabled,
		ClientIP:      dbSession.ClientIp,
		IsBlocked:     dbSession.IsBlocked,
		ExpiresAt:     dbSession.ExpiresAt,
		CreatedAt:     dbSession.CreatedAt,
	}

	// Handle nullable fields
	if dbSession.WebOauthClientID.Valid {
		webOAuthClientID := dbSession.WebOauthClientID.String
		session.WebOAuthClientID = &webOAuthClientID
	}

	if dbSession.OauthAccessToken.Valid {
		oauthAccessToken := dbSession.OauthAccessToken.String
		session.OAuthAccessToken = &oauthAccessToken
	}

	if dbSession.OauthIDToken.Valid {
		oauthIDToken := dbSession.OauthIDToken.String
		session.OAuthIDToken = &oauthIDToken
	}

	return session
}
