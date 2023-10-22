package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sony/sonyflake"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/virsavik/alchemist-template/pkg/postgres"
	"github.com/virsavik/alchemist-template/users/internal/model"
)

type userRepository struct {
	dbConn postgres.ContextExecutor
	idSNF  *sonyflake.Sonyflake
}

func New(dbConn postgres.ContextExecutor) UserRepository {
	return &userRepository{
		dbConn: dbConn,
		idSNF:  sonyflake.NewSonyflake(sonyflake.Settings{}),
	}
}

func (u userRepository) FindOne(ctx context.Context, email string) (*model.User, error) {
	qms := []qm.QueryMod{
		model.UserWhere.DeletedAt.IsNull(),
		model.UserWhere.Email.EQ(email),
	}

	rs, err := model.Users(qms...).One(ctx, u.dbConn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("%w", err)
	}

	return rs, nil
}

func (u userRepository) Save(ctx context.Context, user *model.User) error {
	if user.ID == 0 {
		userID, err := u.idSNF.NextID()
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		user.ID = int64(userID)
	}

	if err := user.Upsert(
		ctx,
		u.dbConn,
		true,
		[]string{model.UserColumns.ID},
		boil.Whitelist(
			model.UserColumns.Email,
			model.UserColumns.IsActive,
			model.UserColumns.CreatedAt,
			model.UserColumns.UpdatedAt,
		),
		boil.Infer(),
	); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (u userRepository) Delete(ctx context.Context, email string) error {
	//TODO implement me
	panic("implement me")
}
