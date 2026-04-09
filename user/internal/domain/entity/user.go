package entity

import (
	"chatterbox/user/internal/domain/valueobject"
	"errors"
)

type User struct {
	id           valueobject.UserID
	email        valueobject.Email
	username     string
	displayName  string
	passwordHash valueobject.PasswordHash
	status       valueobject.Status
}

func NewUser(
	email valueobject.Email,
	username string,
	displayName string,
	passwordHash valueobject.PasswordHash,
) *User {
	return &User{
		id:           valueobject.NewUserID(),
		email:        email,
		username:     username,
		displayName:  displayName,
		passwordHash: passwordHash,
		status:       valueobject.ActiveStatus,
	}
}

func UserFromPrimitives(
	id valueobject.UserID,
	email valueobject.Email,
	username string,
	displayName string,
	passwordHash valueobject.PasswordHash,
	status valueobject.Status,
) *User {
	return &User{
		id:           id,
		email:        email,
		username:     username,
		displayName:  displayName,
		passwordHash: passwordHash,
		status:       status,
	}
}

func (u *User) ID() valueobject.UserID {
	return u.id
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) Username() string {
	return u.username
}

func (u *User) DisplayName() string {
	return u.displayName
}

func (u *User) PasswordHash() valueobject.PasswordHash {
	return u.passwordHash
}

func (u *User) Status() valueobject.Status {
	return u.status
}

func (u *User) Block() error {
	if u.status == valueobject.DeletedStatus {
		return errors.New("cannot block deleted user")
	}
	u.status = valueobject.BlockedStatus
	return nil
}

func (u *User) Activate() {
	u.status = valueobject.ActiveStatus
}

func (u *User) Delete() {
	u.status = valueobject.DeletedStatus
}
