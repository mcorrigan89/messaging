package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcorrigan89/messaging/internal/entities"
	"github.com/mcorrigan89/messaging/internal/repositories/models"
)

type EmailRepository struct {
	utils   ServicesUtils
	DB      *pgxpool.Pool
	queries *models.Queries
}

func NewEmailRepository(utils ServicesUtils, db *pgxpool.Pool, queries *models.Queries) *EmailRepository {
	return &EmailRepository{
		utils:   utils,
		DB:      db,
		queries: queries,
	}
}

func (repo *EmailRepository) GetEmailByID(ctx context.Context, id uuid.UUID) (*entities.Email, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	row, err := repo.queries.GetEmailByID(ctx, id)
	if err != nil {
		repo.utils.logger.Err(err).Ctx(ctx).Msg("Get email by ID")
		return nil, err
	}

	entity := &entities.Email{
		ID:        row.Email.ID,
		MessageID: row.Email.MessageID,
	}

	return entity, nil
}
