package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type WebAuthnUser struct {
	id          []byte
	name        string
	displayName string
	credentials []webauthn.Credential
}

func (u WebAuthnUser) WebAuthnID() []byte                         { return u.id }
func (u WebAuthnUser) WebAuthnName() string                       { return u.name }
func (u WebAuthnUser) WebAuthnDisplayName() string                { return u.displayName }
func (u WebAuthnUser) WebAuthnCredentials() []webauthn.Credential { return u.credentials }
func (u WebAuthnUser) WebAuthnIcon() string                       { return "" }

type BiometricService struct {
	db       *sql.DB
	webauthn *webauthn.WebAuthn
}

func NewBiometricService(db *sql.DB, w *webauthn.WebAuthn) *BiometricService {
	return &BiometricService{db: db, webauthn: w}
}

func (s *BiometricService) BeginRegistration(ctx context.Context, userID int32) ([]byte, error) {
	user, err := s.getUserForWebAuthn(ctx, userID)
	if err != nil {
		return nil, err
	}

	options, sessionData, err := s.webauthn.BeginRegistration(user)
	if err != nil {
		return nil, fmt.Errorf("failed to begin registration: %v", err)
	}

	err = s.storeSessionData(ctx, userID, sessionData)
	if err != nil {
		return nil, fmt.Errorf("failed to store session data: %v", err)
	}

	return json.Marshal(options)
}

func (s *BiometricService) FinishRegistration(ctx context.Context, userID int32, credentialJSON []byte) error {
	user, err := s.getUserForWebAuthn(ctx, userID)
	if err != nil {
		return err
	}

	sessionData, err := s.getSessionData(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get session data: %v", err)
	}

	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(credentialJSON))
	if err != nil {
		return fmt.Errorf("failed to parse credential: %v", err)
	}

	credential, err := s.webauthn.CreateCredential(user, *sessionData, parsedResponse)
	if err != nil {
		return fmt.Errorf("failed to finish registration: %v", err)
	}

	err = s.storeCredential(ctx, userID, parsedResponse, credential)
	if err != nil {
		return fmt.Errorf("failed to store credential: %v", err)
	}

	return nil
}

func (s *BiometricService) BeginLogin(ctx context.Context, email string) ([]byte, error) {
	user, err := s.getUserForWebAuthnByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	options, sessionData, err := s.webauthn.BeginLogin(user)
	if err != nil {
		return nil, fmt.Errorf("failed to begin login: %v", err)
	}

	err = s.storeSessionData(ctx, utils.WebAuthnIDToInt32(user.WebAuthnID()), sessionData)
	if err != nil {
		return nil, fmt.Errorf("failed to store session data: %v", err)
	}

	return json.Marshal(options)
}

func (s *BiometricService) FinishLogin(ctx context.Context, email string, credentialJSON []byte) error {
	user, err := s.getUserForWebAuthnByEmail(ctx, email)
	if err != nil {
		return err
	}

	sessionData, err := s.getSessionData(ctx, utils.WebAuthnIDToInt32(user.WebAuthnID()))
	if err != nil {
		return fmt.Errorf("failed to get session data: %v", err)
	}

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(credentialJSON))
	if err != nil {
		return fmt.Errorf("failed to parse credential: %v", err)
	}

	_, err = s.webauthn.ValidateLogin(user, *sessionData, parsedResponse)
	if err != nil {
		return fmt.Errorf("failed to finish login: %v", err)
	}

	return nil
}

func (s *BiometricService) getUserForWebAuthn(ctx context.Context, userID int32) (*WebAuthnUser, error) {
	queries := models.New(s.db)
	user, err := queries.GetUserForWebAuthn(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	credentials, err := s.getCredentials(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %v", err)
	}

	return &WebAuthnUser{
		id:          user.Webauthnuserid,
		name:        user.Email,
		displayName: user.Name,
		credentials: credentials,
	}, nil
}

func (s *BiometricService) getUserForWebAuthnByEmail(ctx context.Context, email string) (*WebAuthnUser, error) {
	queries := models.New(s.db)
	user, err := queries.GetUserForWebAuthnByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	credentials, err := s.getCredentials(ctx, user.Userid)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %v", err)
	}

	return &WebAuthnUser{
		id:          user.Webauthnuserid,
		name:        user.Email,
		displayName: user.Name,
		credentials: credentials,
	}, nil
}

func (s *BiometricService) getCredentials(ctx context.Context, userID int32) ([]webauthn.Credential, error) {
	queries := models.New(s.db)
	dbCredentials, err := queries.GetWebAuthnCredentials(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %v", err)
	}

	var credentials []webauthn.Credential
	for _, cred := range dbCredentials {
		credentials = append(credentials, webauthn.Credential{
			ID:              cred.Credentialid,
			PublicKey:       cred.Publickey,
			AttestationType: cred.Attestationtype,
			Transport:       nil,                        // Add transport if you're storing it
			Flags:           webauthn.CredentialFlags{}, // Set appropriate flags
			Authenticator: webauthn.Authenticator{
				AAGUID:    cred.Aaguid,
				SignCount: uint32(cred.Signcount),
			},
		})
	}

	return credentials, nil
}

func (s *BiometricService) storeSessionData(ctx context.Context, userID int32, sessionData *webauthn.SessionData) error {
	queries := models.New(s.db)
	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %v", err)
	}

	err = queries.StoreWebAuthnSessionData(ctx, models.StoreWebAuthnSessionDataParams{
		Userid:      userID,
		Sessiondata: sessionDataJSON,
	})
	if err != nil {
		return fmt.Errorf("failed to store session data: %v", err)
	}

	return nil
}

func (s *BiometricService) getSessionData(ctx context.Context, userID int32) (*webauthn.SessionData, error) {
	queries := models.New(s.db)
	sessionDataJSON, err := queries.GetWebAuthnSessionData(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session data: %v", err)
	}

	var sessionData webauthn.SessionData
	err = json.Unmarshal(sessionDataJSON, &sessionData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %v", err)
	}

	return &sessionData, nil
}

func (s *BiometricService) storeCredential(ctx context.Context, userID int32, parsedResponse *protocol.ParsedCredentialCreationData, credential *webauthn.Credential) error {
	queries := models.New(s.db)
	err := queries.StoreWebAuthnCredential(ctx, models.StoreWebAuthnCredentialParams{
		Userid:          userID,
		Credentialid:    credential.ID,
		Publickey:       credential.PublicKey,
		Attestationtype: parsedResponse.Response.AttestationObject.Format,
		Aaguid:          credential.Authenticator.AAGUID,
		Signcount:       int64(credential.Authenticator.SignCount),
	})
	if err != nil {
		return fmt.Errorf("failed to store credential: %v", err)
	}

	return nil
}
