// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package models

import (
	"context"
	"database/sql"
)

const clearResetToken = `-- name: ClearResetToken :exec
UPDATE Users SET reset_token = NULL, reset_token_expires = NULL WHERE UserID = $1
`

func (q *Queries) ClearResetToken(ctx context.Context, userid int32) error {
	_, err := q.db.ExecContext(ctx, clearResetToken, userid)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO Users (Name, Email, Password, UserRole, Status)
VALUES ($1, $2, $3, $4, $5)
RETURNING userid, name, email, password, userrole, status, verificationstatus, deactivatedat, two_factor_secret, two_factor_enabled, failed_login_attempts, last_login_attempt, last_logout, reset_token, reset_token_expires, biometric_token, created_at, updated_at
`

type CreateUserParams struct {
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Userrole string         `json:"userrole"`
	Status   sql.NullString `json:"status"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.Userrole,
		arg.Status,
	)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Userrole,
		&i.Status,
		&i.Verificationstatus,
		&i.Deactivatedat,
		&i.TwoFactorSecret,
		&i.TwoFactorEnabled,
		&i.FailedLoginAttempts,
		&i.LastLoginAttempt,
		&i.LastLogout,
		&i.ResetToken,
		&i.ResetTokenExpires,
		&i.BiometricToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deactivateAccount = `-- name: DeactivateAccount :exec
UPDATE Users
SET DeactivatedAt = CURRENT_TIMESTAMP
WHERE UserID = $1
`

func (q *Queries) DeactivateAccount(ctx context.Context, userid int32) error {
	_, err := q.db.ExecContext(ctx, deactivateAccount, userid)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM Users
WHERE UserID = $1
`

func (q *Queries) DeleteUser(ctx context.Context, userid int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, userid)
	return err
}

const disableTwoFactor = `-- name: DisableTwoFactor :exec
UPDATE Users SET two_factor_enabled = FALSE WHERE UserID = $1
`

func (q *Queries) DisableTwoFactor(ctx context.Context, userid int32) error {
	_, err := q.db.ExecContext(ctx, disableTwoFactor, userid)
	return err
}

const enableTwoFactor = `-- name: EnableTwoFactor :exec
UPDATE Users SET two_factor_enabled = TRUE WHERE UserID = $1
`

func (q *Queries) EnableTwoFactor(ctx context.Context, userid int32) error {
	_, err := q.db.ExecContext(ctx, enableTwoFactor, userid)
	return err
}

const getAccountStatus = `-- name: GetAccountStatus :one
SELECT
    CASE
        WHEN DeactivatedAt IS NULL THEN 'active'
        ELSE 'deactivated'
    END AS status
FROM Users
WHERE UserID = $1
`

func (q *Queries) GetAccountStatus(ctx context.Context, userid int32) (string, error) {
	row := q.db.QueryRowContext(ctx, getAccountStatus, userid)
	var status string
	err := row.Scan(&status)
	return status, err
}

const getPendingUsers = `-- name: GetPendingUsers :many
SELECT userid, name, email, password, userrole, status, verificationstatus, deactivatedat, two_factor_secret, two_factor_enabled, failed_login_attempts, last_login_attempt, last_logout, reset_token, reset_token_expires, biometric_token, created_at, updated_at FROM Users
WHERE Status = 'pending'
`

func (q *Queries) GetPendingUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getPendingUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.Userid,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.Userrole,
			&i.Status,
			&i.Verificationstatus,
			&i.Deactivatedat,
			&i.TwoFactorSecret,
			&i.TwoFactorEnabled,
			&i.FailedLoginAttempts,
			&i.LastLoginAttempt,
			&i.LastLogout,
			&i.ResetToken,
			&i.ResetTokenExpires,
			&i.BiometricToken,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByBiometricToken = `-- name: GetUserByBiometricToken :one
SELECT userid, name, email, password, userrole, status, verificationstatus, deactivatedat, two_factor_secret, two_factor_enabled, failed_login_attempts, last_login_attempt, last_logout, reset_token, reset_token_expires, biometric_token, created_at, updated_at FROM Users WHERE biometric_token = $1 LIMIT 1
`

func (q *Queries) GetUserByBiometricToken(ctx context.Context, biometricToken sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByBiometricToken, biometricToken)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Userrole,
		&i.Status,
		&i.Verificationstatus,
		&i.Deactivatedat,
		&i.TwoFactorSecret,
		&i.TwoFactorEnabled,
		&i.FailedLoginAttempts,
		&i.LastLoginAttempt,
		&i.LastLogout,
		&i.ResetToken,
		&i.ResetTokenExpires,
		&i.BiometricToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT userid, name, email, password, userrole, status, verificationstatus, deactivatedat, two_factor_secret, two_factor_enabled, failed_login_attempts, last_login_attempt, last_logout, reset_token, reset_token_expires, biometric_token, created_at, updated_at FROM Users
WHERE Email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Userrole,
		&i.Status,
		&i.Verificationstatus,
		&i.Deactivatedat,
		&i.TwoFactorSecret,
		&i.TwoFactorEnabled,
		&i.FailedLoginAttempts,
		&i.LastLoginAttempt,
		&i.LastLogout,
		&i.ResetToken,
		&i.ResetTokenExpires,
		&i.BiometricToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT userid, name, email, password, userrole, status, verificationstatus, deactivatedat, two_factor_secret, two_factor_enabled, failed_login_attempts, last_login_attempt, last_logout, reset_token, reset_token_expires, biometric_token, created_at, updated_at FROM Users
WHERE UserID = $1
`

func (q *Queries) GetUserByID(ctx context.Context, userid int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, userid)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Userrole,
		&i.Status,
		&i.Verificationstatus,
		&i.Deactivatedat,
		&i.TwoFactorSecret,
		&i.TwoFactorEnabled,
		&i.FailedLoginAttempts,
		&i.LastLoginAttempt,
		&i.LastLogout,
		&i.ResetToken,
		&i.ResetTokenExpires,
		&i.BiometricToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByResetToken = `-- name: GetUserByResetToken :one
SELECT userid, name, email, password, userrole, status, verificationstatus, deactivatedat, two_factor_secret, two_factor_enabled, failed_login_attempts, last_login_attempt, last_logout, reset_token, reset_token_expires, biometric_token, created_at, updated_at FROM Users WHERE reset_token = $1 AND reset_token_expires > NOW() LIMIT 1
`

func (q *Queries) GetUserByResetToken(ctx context.Context, resetToken sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByResetToken, resetToken)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Userrole,
		&i.Status,
		&i.Verificationstatus,
		&i.Deactivatedat,
		&i.TwoFactorSecret,
		&i.TwoFactorEnabled,
		&i.FailedLoginAttempts,
		&i.LastLoginAttempt,
		&i.LastLogout,
		&i.ResetToken,
		&i.ResetTokenExpires,
		&i.BiometricToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserWithAuthDetails = `-- name: GetUserWithAuthDetails :one
SELECT userid, name, email, password, userrole, status, verificationstatus, deactivatedat, two_factor_secret, two_factor_enabled, failed_login_attempts, last_login_attempt, last_logout, reset_token, reset_token_expires, biometric_token, created_at, updated_at FROM Users WHERE UserID = $1
`

func (q *Queries) GetUserWithAuthDetails(ctx context.Context, userid int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserWithAuthDetails, userid)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Userrole,
		&i.Status,
		&i.Verificationstatus,
		&i.Deactivatedat,
		&i.TwoFactorSecret,
		&i.TwoFactorEnabled,
		&i.FailedLoginAttempts,
		&i.LastLoginAttempt,
		&i.LastLogout,
		&i.ResetToken,
		&i.ResetTokenExpires,
		&i.BiometricToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUsersByStatus = `-- name: GetUsersByStatus :many
SELECT userid, name, email, password, userrole, status, verificationstatus, deactivatedat, two_factor_secret, two_factor_enabled, failed_login_attempts, last_login_attempt, last_logout, reset_token, reset_token_expires, biometric_token, created_at, updated_at FROM Users
WHERE Status = $1
`

func (q *Queries) GetUsersByStatus(ctx context.Context, status sql.NullString) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsersByStatus, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.Userid,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.Userrole,
			&i.Status,
			&i.Verificationstatus,
			&i.Deactivatedat,
			&i.TwoFactorSecret,
			&i.TwoFactorEnabled,
			&i.FailedLoginAttempts,
			&i.LastLoginAttempt,
			&i.LastLogout,
			&i.ResetToken,
			&i.ResetTokenExpires,
			&i.BiometricToken,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const incrementFailedLoginAttempts = `-- name: IncrementFailedLoginAttempts :exec
UPDATE Users
SET failed_login_attempts = failed_login_attempts + 1,
    last_login_attempt = NOW()
WHERE UserID = $1
`

func (q *Queries) IncrementFailedLoginAttempts(ctx context.Context, userid int32) error {
	_, err := q.db.ExecContext(ctx, incrementFailedLoginAttempts, userid)
	return err
}

const reactivateAccount = `-- name: ReactivateAccount :exec
UPDATE Users
SET DeactivatedAt = NULL
WHERE UserID = $1
`

func (q *Queries) ReactivateAccount(ctx context.Context, userid int32) error {
	_, err := q.db.ExecContext(ctx, reactivateAccount, userid)
	return err
}

const resetFailedLoginAttempts = `-- name: ResetFailedLoginAttempts :exec
UPDATE Users SET failed_login_attempts = 0 WHERE UserID = $1
`

func (q *Queries) ResetFailedLoginAttempts(ctx context.Context, userid int32) error {
	_, err := q.db.ExecContext(ctx, resetFailedLoginAttempts, userid)
	return err
}

const setBiometricToken = `-- name: SetBiometricToken :exec
UPDATE Users SET biometric_token = $2 WHERE UserID = $1
`

type SetBiometricTokenParams struct {
	Userid         int32          `json:"userid"`
	BiometricToken sql.NullString `json:"biometric_token"`
}

func (q *Queries) SetBiometricToken(ctx context.Context, arg SetBiometricTokenParams) error {
	_, err := q.db.ExecContext(ctx, setBiometricToken, arg.Userid, arg.BiometricToken)
	return err
}

const setResetToken = `-- name: SetResetToken :exec
UPDATE Users SET reset_token = $2, reset_token_expires = $3 WHERE UserID = $1
`

type SetResetTokenParams struct {
	Userid            int32          `json:"userid"`
	ResetToken        sql.NullString `json:"reset_token"`
	ResetTokenExpires sql.NullTime   `json:"reset_token_expires"`
}

func (q *Queries) SetResetToken(ctx context.Context, arg SetResetTokenParams) error {
	_, err := q.db.ExecContext(ctx, setResetToken, arg.Userid, arg.ResetToken, arg.ResetTokenExpires)
	return err
}

const updateLastLoginAttempt = `-- name: UpdateLastLoginAttempt :exec
UPDATE Users SET last_login_attempt = NOW() WHERE UserID = $1
`

func (q *Queries) UpdateLastLoginAttempt(ctx context.Context, userid int32) error {
	_, err := q.db.ExecContext(ctx, updateLastLoginAttempt, userid)
	return err
}

const updateLastLogout = `-- name: UpdateLastLogout :exec
UPDATE Users
SET last_logout = $2
WHERE UserID = $1
`

type UpdateLastLogoutParams struct {
	Userid     int32        `json:"userid"`
	LastLogout sql.NullTime `json:"last_logout"`
}

func (q *Queries) UpdateLastLogout(ctx context.Context, arg UpdateLastLogoutParams) error {
	_, err := q.db.ExecContext(ctx, updateLastLogout, arg.Userid, arg.LastLogout)
	return err
}

const updateUser = `-- name: UpdateUser :one
UPDATE Users
SET Name = $2, Email = $3, Password = $4, UserRole = $5, VerificationStatus = $6, Status = $7
WHERE UserID = $1
RETURNING userid, name, email, password, userrole, status, verificationstatus, deactivatedat, two_factor_secret, two_factor_enabled, failed_login_attempts, last_login_attempt, last_logout, reset_token, reset_token_expires, biometric_token, created_at, updated_at
`

type UpdateUserParams struct {
	Userid             int32          `json:"userid"`
	Name               string         `json:"name"`
	Email              string         `json:"email"`
	Password           string         `json:"password"`
	Userrole           string         `json:"userrole"`
	Verificationstatus sql.NullBool   `json:"verificationstatus"`
	Status             sql.NullString `json:"status"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Userid,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.Userrole,
		arg.Verificationstatus,
		arg.Status,
	)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Userrole,
		&i.Status,
		&i.Verificationstatus,
		&i.Deactivatedat,
		&i.TwoFactorSecret,
		&i.TwoFactorEnabled,
		&i.FailedLoginAttempts,
		&i.LastLoginAttempt,
		&i.LastLogout,
		&i.ResetToken,
		&i.ResetTokenExpires,
		&i.BiometricToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE Users
SET Password = $2
WHERE UserID = $1
`

type UpdateUserPasswordParams struct {
	Userid   int32  `json:"userid"`
	Password string `json:"password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.Userid, arg.Password)
	return err
}

const updateUserStatus = `-- name: UpdateUserStatus :exec
UPDATE Users
SET Status = $2
WHERE UserID = $1
`

type UpdateUserStatusParams struct {
	Userid int32          `json:"userid"`
	Status sql.NullString `json:"status"`
}

func (q *Queries) UpdateUserStatus(ctx context.Context, arg UpdateUserStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateUserStatus, arg.Userid, arg.Status)
	return err
}

const updateUserTwoFactorSecret = `-- name: UpdateUserTwoFactorSecret :exec
UPDATE Users SET two_factor_secret = $2 WHERE UserID = $1
`

type UpdateUserTwoFactorSecretParams struct {
	Userid          int32          `json:"userid"`
	TwoFactorSecret sql.NullString `json:"two_factor_secret"`
}

func (q *Queries) UpdateUserTwoFactorSecret(ctx context.Context, arg UpdateUserTwoFactorSecretParams) error {
	_, err := q.db.ExecContext(ctx, updateUserTwoFactorSecret, arg.Userid, arg.TwoFactorSecret)
	return err
}
