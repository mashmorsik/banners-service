package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/mashmorsik/banners-service/infrastructure/data"
	"github.com/mashmorsik/banners-service/pkg/models"
	"github.com/mashmorsik/logger"
	errs "github.com/pkg/errors"
	"time"
)

type BannerRepo struct {
	Ctx  context.Context
	data *data.Data
}

func NewBannerRepo(ctx context.Context, data *data.Data) *BannerRepo {
	return &BannerRepo{Ctx: ctx, data: data}
}

func (br *BannerRepo) GetForUser(bannerID int) (*models.Content, error) {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	var contentJSON []byte

	err := br.data.Master().QueryRowContext(ctx,
		`SELECT content
		FROM banner_version
		WHERE banner_id = $1
		AND is_active = true`, bannerID).Scan(&contentJSON)
	if err != nil {
		return nil, errs.WithMessagef(err, "failed to get banner content with bannerID %d", bannerID)
	}

	var content models.Content
	if err = json.Unmarshal(contentJSON, &content); err != nil {
		return nil, errs.WithMessagef(err, "failed to unmarshal content with bannerID %d", bannerID)
	}

	return &content, nil
}

func (br *BannerRepo) GetForAdmin(b *models.Banner, limit, offset int) ([]*models.Banner, error) {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	query := `
        SELECT b.id, b.created_at, b.updated_at, b.tag_ids, b.feature_id, b.is_active, bv.content
		FROM banner b
		JOIN (
			SELECT banner_id, MAX(version) AS max_version
			FROM banner_version
			GROUP BY banner_id
		) latest_version ON bv.banner_id = b.id
		JOIN banner_version bv ON bv.banner_id = latest_version.banner_id AND bv.version = latest_version.max_version
		WHERE ($1::int[] IS NULL OR b.tag_ids @> $1)
		AND ($2::int = 0 OR b.feature_id = $2)
		LIMIT $3 OFFSET $4;

    `

	rows, err := br.data.Master().QueryContext(ctx, query, pq.Array(b.TagIDs), b.FeatureID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			logger.Errf("failed to close rows, err: %v", err)
			return
		}
	}(rows)

	var banners []*models.Banner
	for rows.Next() {
		var banner models.Banner
		var contentJSON []byte
		if err = rows.Scan(&banner.ID, &banner.CreatedAt, &banner.UpdatedAt, &banner.TagIDs, &banner.FeatureID, &banner.IsActive, &contentJSON); err != nil {
			return nil, errs.WithMessagef(err, "failed to scan rows")
		}
		if err = json.Unmarshal(contentJSON, &banner.Content); err != nil {
			return nil, errs.WithMessagef(err, "failed to unmarshal content JSON")
		}
		banners = append(banners, &banner)
	}
	if err = rows.Err(); err != nil {
		return nil, errs.WithMessagef(err, "failed to fetch rows")
	}

	return banners, nil
}

func (br *BannerRepo) CreateBanner(tx *sql.Tx, b *models.Banner) (int, error) {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	var createdBannerID int

	err := tx.QueryRowContext(ctx,
		`INSERT INTO banner (created_at, updated_at, tag_ids, feature_id, is_active, version)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id`, b.CreatedAt, b.UpdatedAt, b.TagIDs, b.FeatureID, true, 1).Scan(&createdBannerID)
	if err != nil {
		return 0, errs.New("failed to exec query: Create")
	}

	return createdBannerID, nil
}
func (br *BannerRepo) CreateVersion(tx *sql.Tx, b *models.Banner) error {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	contentJSON, err := json.Marshal(b.Content)
	if err != nil {
		return errs.WithMessagef(err, "failed to marshal content to JSON, content: %v", b.Content)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO banner_version (banner_id, version, content, updated_at, is_active)
				VALUES ($1, $2, $3, $4, $5)`, b.ID, contentJSON, b.UpdatedAt, true)
	if err != nil {
		return errs.New("failed to exec query: Create")
	}

	return nil
}

func (br *BannerRepo) Create(b *models.Banner) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	tx, err := br.data.Master().Begin()
	if err != nil {
		logger.Errf("can't begin transaction, err: %s", err)
	}
	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	b.ID, err = br.CreateBanner(tx, b)
	if err != nil {
		return errs.WithMessagef(err, "failed to create banner with id %d", b.ID)
	}

	err = br.CreateVersion(tx, b)
	if err != nil {
		return errs.WithMessagef(err, "failed to create banner with id %d", b.ID)
	}

	return nil
}

func (br *BannerRepo) UpdateBanner(tx *sql.Tx, b *models.Banner) error {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	_, err := tx.ExecContext(ctx,
		`UPDATE banner
			SET updated_at = $1, tag_ids = $2, feature_id = $3, version = version + 1
			WHERE id = $4`, b.UpdatedAt, b.TagIDs, b.FeatureID, b.ID)
	if err != nil {
		return errs.New("failed to exec query: Create")
	}

	return nil
}

func (br *BannerRepo) UpdateBannerVersion(tx *sql.Tx, b *models.Banner) error {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	contentJSON, err := json.Marshal(b.Content)
	if err != nil {
		return errs.WithMessagef(err, "failed to marshal content to JSON, content: %v", b.Content)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO banner_version (banner_id, content, version, updated_at, is_active)
				VALUES ($1, $2, version + 1, $4, $5)`, b.ID, contentJSON, b.UpdatedAt, true)
	if err != nil {
		return errs.New("failed to exec query: Create")
	}

	return nil
}

func (br *BannerRepo) SetOldVersionInactive(tx *sql.Tx, b *models.Banner) error {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	_, err := tx.ExecContext(ctx,
		`UPDATE banner_version
				 SET is_active = false
				 WHERE id = (
					 SELECT id FROM banner_version
					 WHERE banner_id = $1
					 ORDER BY updated_at DESC
					 LIMIT 1
				 )`, b.ID)
	if err != nil {
		return errs.WithMessagef(err, "failed to deactivate last updated version of banner")
	}

	return nil
}

func (br *BannerRepo) Update(b *models.Banner) error {
	b.UpdatedAt = time.Now()

	tx, err := br.data.Master().Begin()
	if err != nil {
		logger.Errf("can't begin transaction, err: %s", err)
	}
	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	err = br.UpdateBanner(tx, b)
	if err != nil {
		return errs.WithMessagef(err, "failed to update banner with id %d", b.ID)
	}

	err = br.SetOldVersionInactive(tx, b)
	if err != nil {
		return errs.WithMessagef(err, "failed to set old version inactive with id %d", b.ID)
	}

	err = br.UpdateBannerVersion(tx, b)
	if err != nil {
		return errs.WithMessagef(err, "failed to update banner with id %d", b.ID)
	}

	return nil
}

func (br *BannerRepo) Delete(bannerID int) error {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	_, err := br.data.Master().ExecContext(ctx,
		`DELETE 
		FROM banner
		WHERE id = $1`, bannerID)
	if err != nil {
		return errs.WithMessagef(err, "failed to exec query: DeleteBanner")
	}

	return nil
}

func (br *BannerRepo) CheckTagFeatureOverlap(b *models.Banner) (int, error) {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	var bannerID int

	err := br.data.Master().QueryRowContext(ctx,
		`SELECT id
		FROM banner
		WHERE (tag_ids @> $1)
		AND (feature_id = $2)
		AND is_active = true`, pq.Array(b.TagIDs), b.FeatureID).Scan(&bannerID)
	if err != nil {
		return 0, errs.WithMessagef(err, "active banner with this combination of tagIDs and featureID is not found")
	}

	return bannerID, nil
}
