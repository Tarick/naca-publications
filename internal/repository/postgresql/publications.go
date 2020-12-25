package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tarick/naca-publications/internal/entity"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

// CreatePublication inserts new publisher into db
func (repo *Repository) CreatePublication(ctx context.Context, p *entity.Publication) error {
	if repo.publicationExists(ctx, p) {
		return errors.New("publication already exists")
	}
	_, err := repo.pool.Exec(ctx, "insert into publications (uuid, name, description, type, publisher_uuid, language_code) values ($1, $2, $3, $4, $5, $6)",
		p.UUID, p.Name, p.Description, p.Type, p.PublisherUUID, p.LanguageCode)
	return err
}

func (repo *Repository) publicationExists(ctx context.Context, p *entity.Publication) bool {
	var exists bool
	row := repo.pool.QueryRow(ctx, "select exists (select 1 from publications where uuid=$1 or (publisher_uuid=$2 and name=$3))", p.UUID, p.PublisherUUID, p.Name)
	if err := row.Scan(&exists); err != nil {
		panic(err)
	}
	if exists == true {
		return true
	}
	return false
}

// UpdatePublication updates Publication in db
func (repo *Repository) UpdatePublication(ctx context.Context, p *entity.Publication) error {
	_, err := repo.pool.Exec(ctx, "update publications set name=$1, description=$2, language_code=$3 where uuid=$4", p.Name, p.Description, p.LanguageCode, p.UUID)
	return err
}

// DeletePublication removes Publications from db
func (repo *Repository) DeletePublication(ctx context.Context, uuid uuid.UUID) error {
	result, err := repo.pool.Exec(ctx, "delete from publications where uuid=$1", uuid)
	if err != nil {
		return err
	}
	// Should never be called since we find it before, but still (multiple requests could clash)
	if result.RowsAffected() != 1 {
		return fmt.Errorf("publication with UUID %v wasn't deleted from db", uuid)
	}
	return err
}

// GetPublication returns Publication from db
func (repo *Repository) GetPublication(ctx context.Context, uuid uuid.UUID) (*entity.Publication, error) {
	p := &entity.Publication{}
	err := repo.pool.QueryRow(ctx, "select uuid, name, description, language_code, publisher_uuid, type from publications where uuid=$1", uuid).
		Scan(&p.UUID, &p.Name, &p.Description, &p.LanguageCode, &p.PublisherUUID, &p.Type)
	if err != nil && err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetPublications returns list of Publication from db
func (repo *Repository) GetPublications(ctx context.Context) ([]*entity.Publication, error) {
	rows, err := repo.pool.Query(ctx, "select uuid, name, description, language_code, publisher_uuid, type from publications")
	if err != nil {
		return nil, err
	}
	publications := []*entity.Publication{}
	for rows.Next() {
		p := &entity.Publication{}
		if err := rows.Scan(&p.UUID, &p.Name, &p.Description, &p.LanguageCode, &p.PublisherUUID, &p.Type); err != nil {
			return nil, err
		}
		publications = append(publications, p)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return publications, nil
}

// Healthcheck is needed for application healtchecks
func (repo *Repository) Healthcheck(ctx context.Context) error {
	var exists bool
	row := repo.pool.QueryRow(ctx, "select exists (select 1 from publications limit 1)")
	if err := row.Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}
	return fmt.Errorf("failure checking access to 'publications' table")
}
