package respositories

import (
	"context"

	"gorm.io/gorm"
)

type dbRepo struct {
	db *gorm.DB
}

func (r *dbRepo) getDB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}
