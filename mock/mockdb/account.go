package mockdb

import (
	"context"

	"github.com/calvinchengx/gin-go-pg/model"
)

// Account database mock
type Account struct {
	CreateFn                  func(*model.User) (*model.User, error)
	CreateAndVerifyFn         func(*model.User) (*model.Verification, error)
	CreateWithMobileFn        func(*model.User) error
	ChangePasswordFn          func(context.Context, *model.User) error
	FindVerificationTokenFn   func(context.Context, string) (*model.Verification, error)
	DeleteVerificationTokenFn func(context.Context, *model.Verification) error
}

// Create mock
func (a *Account) Create(usr *model.User) (*model.User, error) {
	return a.CreateFn(usr)
}

// CreateAndVerify mock
func (a *Account) CreateAndVerify(usr *model.User) (*model.Verification, error) {
	return a.CreateAndVerifyFn(usr)
}

// CreateWithMobile mock
func (a *Account) CreateWithMobile(usr *model.User) error {
	return a.CreateWithMobileFn(usr)
}

// ChangePassword mock
func (a *Account) ChangePassword(c context.Context, usr *model.User) error {
	return a.ChangePasswordFn(c, usr)
}

// FindVerificationToken mock
func (a *Account) FindVerificationToken(c context.Context, token string) (*model.Verification, error) {
	return a.FindVerificationTokenFn(c, token)
}

// DeleteVerificationToken mock
func (a *Account) DeleteVerificationToken(c context.Context, v *model.Verification) error {
	return a.DeleteVerificationTokenFn(c, v)
}
