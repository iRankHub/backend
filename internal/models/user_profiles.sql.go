// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user_profiles.sql

package models

import (
	"context"
	"database/sql"
)

const createUserProfile = `-- name: CreateUserProfile :one
INSERT INTO UserProfiles (UserID, Name, UserRole, Email, Password, VerificationStatus, Address, Phone, Bio, ProfilePicture)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING profileid, userid, name, userrole, email, password, address, phone, bio, profilepicture, verificationstatus
`

type CreateUserProfileParams struct {
	Userid             int32          `json:"userid"`
	Name               string         `json:"name"`
	Userrole           string         `json:"userrole"`
	Email              string         `json:"email"`
	Password           string         `json:"password"`
	Verificationstatus sql.NullBool   `json:"verificationstatus"`
	Address            sql.NullString `json:"address"`
	Phone              sql.NullString `json:"phone"`
	Bio                sql.NullString `json:"bio"`
	Profilepicture     []byte         `json:"profilepicture"`
}

func (q *Queries) CreateUserProfile(ctx context.Context, arg CreateUserProfileParams) (Userprofile, error) {
	row := q.db.QueryRowContext(ctx, createUserProfile,
		arg.Userid,
		arg.Name,
		arg.Userrole,
		arg.Email,
		arg.Password,
		arg.Verificationstatus,
		arg.Address,
		arg.Phone,
		arg.Bio,
		arg.Profilepicture,
	)
	var i Userprofile
	err := row.Scan(
		&i.Profileid,
		&i.Userid,
		&i.Name,
		&i.Userrole,
		&i.Email,
		&i.Password,
		&i.Address,
		&i.Phone,
		&i.Bio,
		&i.Profilepicture,
		&i.Verificationstatus,
	)
	return i, err
}

const deleteUserProfile = `-- name: DeleteUserProfile :exec
DELETE FROM UserProfiles
WHERE UserID = $1
`

func (q *Queries) DeleteUserProfile(ctx context.Context, userid int32) error {
	_, err := q.db.ExecContext(ctx, deleteUserProfile, userid)
	return err
}

const getUserProfile = `-- name: GetUserProfile :one
SELECT profileid, userid, name, userrole, email, password, address, phone, bio, profilepicture, verificationstatus FROM UserProfiles
WHERE UserID = $1
`

func (q *Queries) GetUserProfile(ctx context.Context, userid int32) (Userprofile, error) {
	row := q.db.QueryRowContext(ctx, getUserProfile, userid)
	var i Userprofile
	err := row.Scan(
		&i.Profileid,
		&i.Userid,
		&i.Name,
		&i.Userrole,
		&i.Email,
		&i.Password,
		&i.Address,
		&i.Phone,
		&i.Bio,
		&i.Profilepicture,
		&i.Verificationstatus,
	)
	return i, err
}

const updateUserProfile = `-- name: UpdateUserProfile :one
UPDATE UserProfiles
SET Name = $2, UserRole = $3, Email = $4, Password = $5, VerificationStatus = $6, Address = $7, Phone = $8, Bio = $9, ProfilePicture = $10
WHERE UserID = $1
RETURNING profileid, userid, name, userrole, email, password, address, phone, bio, profilepicture, verificationstatus
`

type UpdateUserProfileParams struct {
	Userid             int32          `json:"userid"`
	Name               string         `json:"name"`
	Userrole           string         `json:"userrole"`
	Email              string         `json:"email"`
	Password           string         `json:"password"`
	Verificationstatus sql.NullBool   `json:"verificationstatus"`
	Address            sql.NullString `json:"address"`
	Phone              sql.NullString `json:"phone"`
	Bio                sql.NullString `json:"bio"`
	Profilepicture     []byte         `json:"profilepicture"`
}

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) (Userprofile, error) {
	row := q.db.QueryRowContext(ctx, updateUserProfile,
		arg.Userid,
		arg.Name,
		arg.Userrole,
		arg.Email,
		arg.Password,
		arg.Verificationstatus,
		arg.Address,
		arg.Phone,
		arg.Bio,
		arg.Profilepicture,
	)
	var i Userprofile
	err := row.Scan(
		&i.Profileid,
		&i.Userid,
		&i.Name,
		&i.Userrole,
		&i.Email,
		&i.Password,
		&i.Address,
		&i.Phone,
		&i.Bio,
		&i.Profilepicture,
		&i.Verificationstatus,
	)
	return i, err
}
