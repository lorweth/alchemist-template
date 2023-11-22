package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/virsavik/alchemist-template/users/internal/adapters/repository/generator"
	"github.com/virsavik/alchemist-template/users/internal/adapters/repository/orm"
	"github.com/virsavik/alchemist-template/users/internal/core/domain"
	"github.com/virsavik/alchemist-template/users/internal/core/ports"
)

func (r Repository) GetAll(ctx context.Context, input ports.GetUserInput) ([]domain.User, error) {
	qms := []qm.QueryMod{
		orm.UserWhere.DeletedAt.IsNull(),
	}

	if input.ID != 0 {
		qms = append(qms, orm.UserWhere.ID.EQ(input.ID))
	}

	if input.Email != "" {
		qms = append(qms, orm.UserWhere.Email.EQ(input.Email))
	}

	users, err := orm.Users(qms...).All(ctx, r.db)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rs := make([]domain.User, len(users))
	for idx, user := range users {
		rs[idx] = domain.User{
			ID:        user.ID,
			Email:     user.Email,
			CreateAt:  user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: user.DeletedAt.Ptr(),
		}
	}

	return rs, nil
}

func (r Repository) GetOne(ctx context.Context, input ports.GetUserInput) (domain.User, error) {
	// Get all user by given input
	selectedUsers, err := r.GetAll(ctx, input)
	if err != nil {
		return domain.User{}, err
	}

	// Just user not found
	if len(selectedUsers) == 0 {
		return domain.User{}, nil
	}

	if len(selectedUsers) > 1 {
		return domain.User{}, ErrMoreThanOneRecord
	}

	return selectedUsers[0], nil
}

func (r Repository) Save(ctx context.Context, user domain.User) (domain.User, error) {
	// Generate ID if new user
	if user.ID == 0 {
		newID, err := generator.UserIDGenerator.NextID()
		if err != nil {
			return domain.User{}, errors.WithStack(err)
		}

		user.ID = int64(newID)
	}

	// Convert to user model
	userORM := orm.User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreateAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: null.TimeFromPtr(user.DeletedAt),
	}

	// Save to database
	if err := userORM.Upsert(ctx, r.db, true,
		[]string{orm.UserColumns.Email},
		boil.Whitelist(orm.UserColumns.DeletedAt),
		boil.Infer(),
	); err != nil {
		return domain.User{}, errors.WithStack(err)
	}

	return user, nil
}

func (r Repository) Delete(ctx context.Context, user domain.User) error {
	// Return error if user id is zero
	if user.ID == 0 {
		return ErrUserIDInvalid
	}

	// Update deleted at
	now := timeNowWrapper()
	user.DeletedAt = &now

	// Save to database
	if _, err := r.Save(ctx, user); err != nil {
		return err
	}

	return nil
}
