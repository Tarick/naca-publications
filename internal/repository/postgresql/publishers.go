package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tarick/naca-publications/internal/entity"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

// CreatePublisher inserts new publisher into db
func (repo *Repository) CreatePublisher(ctx context.Context, p *entity.Publisher) error {
	_, err := repo.pool.Exec(ctx, "insert into publishers (uuid, name, url) values ($1, $2, $3)", p.UUID, p.Name, p.URL)
	return err
}

// UpdatePublisher updates Publisher in db
func (repo *Repository) UpdatePublisher(ctx context.Context, p *entity.Publisher) error {
	_, err := repo.pool.Exec(ctx, "update publishers set name=$1, url=$2 where uuid=$3", p.Name, p.URL, p.UUID)
	return err
}

// DeletePublisher removes Publishers from db
func (repo *Repository) DeletePublisher(ctx context.Context, uuid uuid.UUID) error {
	result, err := repo.pool.Exec(ctx, "delete from publishers where uuid=$1", uuid)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return errors.New(fmt.Sprint("publisher delete from db execution didn't delete record for UUID ", uuid))
	}
	return err
}

// GetPublisher returns Publisher from db
func (repo *Repository) GetPublisher(ctx context.Context, uuid uuid.UUID) (*entity.Publisher, error) {
	p := &entity.Publisher{}
	err := repo.pool.QueryRow(ctx, "select uuid, name, url from publishers where uuid=$1", uuid).Scan(&p.UUID, &p.Name, &p.URL)
	if err != nil && err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetPublishers returns list of Publisher from db
func (repo *Repository) GetPublishers(ctx context.Context) ([]*entity.Publisher, error) {
	rows, err := repo.pool.Query(ctx, "select uuid, name, url from publishers")
	if err != nil {
		return nil, err
	}
	publishers := []*entity.Publisher{}
	for rows.Next() {
		p := &entity.Publisher{}
		if err := rows.Scan(&p.UUID, &p.Name, &p.URL); err != nil {
			return nil, err
		}
		publishers = append(publishers, p)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return publishers, nil
}

// GetPublicationsByPublisher returns list of Publication filterered by publisher uuid
func (repo *Repository) GetPublicationsByPublisher(ctx context.Context, publisherUUID uuid.UUID) ([]*entity.Publication, error) {
	rows, err := repo.pool.Query(ctx, "select uuid, name, description, language_code, publisher_uuid, pt.type from publications  join publication_types pt on type_id = id where publisher_uuid=$1", publisherUUID)
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
