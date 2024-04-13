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

func (br *BannerRepo) GetForUser(b *models.Banner) (*models.Content, error) {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	var contentJSON []byte

	err := br.data.Master().QueryRowContext(ctx,
		`SELECT bc.content
		FROM banner_content bc
		JOIN banner b ON bc.banner_id = b.id
		JOIN banner_feature_tag bft ON b.id = bft.banner_id
		WHERE bft.tag_id = $1
		AND bft.feature_id = $2
		AND b.is_active = true
		AND b.active_version = bc.version`, b.TagIDs[0], b.FeatureID).Scan(&contentJSON)
	if err != nil {
		return nil, errs.WithMessagef(err, "failed to get banner content with bannerID %d", b.ID)
	}

	var content models.Content
	if err = json.Unmarshal(contentJSON, &content); err != nil {
		return nil, errs.WithMessagef(err, "failed to unmarshal content with bannerID %d", b.ID)
	}

	return &content, nil
}

func (br *BannerRepo) GetForAdmin(b *models.Banner, limit, offset int) ([]*models.Banner, error) {
	var banners []*models.Banner

	tagID := b.TagIDs[0]
	featureID := b.FeatureID

	var queryLimit, queryOffset interface{}
	if limit != 0 {
		queryLimit = limit
	}
	if offset != 0 {
		queryOffset = offset
	}

	rows, err := br.data.Master().QueryContext(br.Ctx, `
        SELECT
            b.id,
            b.created_at,
            b.updated_at,
            array_agg(bft.tag_id) AS tag_ids,
            bft.feature_id,
            b.is_active,
            bc.content,
            bc.version
        FROM
            banner b
        JOIN
            banner_content bc ON b.id = bc.banner_id
        JOIN
            banner_feature_tag bft ON b.id = bft.banner_id
        WHERE
            (bft.feature_id = $1 OR $1 IS NULL)
            AND (bft.tag_id = $2 OR $2 IS NULL)
        GROUP BY
            b.id,
            b.created_at,
            b.updated_at,
            bft.feature_id,
            b.is_active,
            bc.content,
            bc.version
        ORDER BY
            b.id
        LIMIT $3 OFFSET $4;
    `, featureID, tagID, queryLimit, queryOffset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var banner models.Banner
		var contentJSON []byte
		if err = rows.Scan(&banner.ID, &banner.CreatedAt, &banner.UpdatedAt, pq.Array(&banner.TagIDs), &banner.FeatureID, &banner.IsActive, &contentJSON, &banner.Version); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(contentJSON, &banner.Content); err != nil {
			return nil, err
		}
		banners = append(banners, &banner)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (br *BannerRepo) CreateBanner(tx *sql.Tx, b *models.Banner) (int, error) {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	var (
		createdBannerID, activeVersion int
	)

	if b.IsActive == true {
		b.Version = 1
		activeVersion = 1
	} else {
		b.Version = 1
		activeVersion = 0
	}

	err := tx.QueryRowContext(ctx,
		`INSERT INTO banner (created_at, updated_at, is_active, active_version, last_version)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id`, b.CreatedAt, b.UpdatedAt, b.IsActive, activeVersion, b.Version).Scan(&createdBannerID)
	if err != nil {
		return 0, errs.New("fail to insert into banner table while exec Create")
	}

	return createdBannerID, nil
}

func (br *BannerRepo) CreateContent(tx *sql.Tx, b *models.Banner) error {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	contentJSON, err := json.Marshal(b.Content)
	if err != nil {
		return errs.WithMessagef(err, "fail to marshal content to JSON, content: %v", b.Content)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO banner_content (banner_id, version, content, updated_at)
				VALUES ($1, $2, $3, $4)`, b.ID, b.Version, contentJSON, b.UpdatedAt)
	if err != nil {
		return errs.New("fail to insert into banner_content table while exec Create")
	}

	return nil
}

func (br *BannerRepo) CreateFeatureTags(tx *sql.Tx, b *models.Banner) error {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	for _, tagID := range b.TagIDs {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO banner_feature_tag(banner_id, feature_id, tag_id, version, updated_at)
			VALUES($1, $2, $3, $4, $5)`, b.ID, b.FeatureID, tagID, b.Version, b.UpdatedAt)
		if err != nil {
			return errs.WithMessagef(err, "fail to insert into banner_feature_tag table while exec Create")
		}
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
		return errs.WithMessagef(err, "fail to create banner with id %d", b.ID)
	}

	err = br.CreateContent(tx, b)
	if err != nil {
		return errs.WithMessagef(err, "fail to create version with id %d", b.ID)
	}

	err = br.CreateFeatureTags(tx, b)
	if err != nil {
		return errs.WithMessagef(err, "fail to create tags with id %d", b.ID)
	}

	if err = tx.Commit(); err != nil {
		logger.Errf("failed to commit transaction CreateBanner: %s", err)
		return err
	}

	return nil
}

func (br *BannerRepo) UpdateBanner(tx *sql.Tx, b *models.Banner, lastVersion int) error {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	var activeVersion int

	if b.IsActive == true {
		activeVersion = b.Version
	} else {
		activeVersion = 0
	}

	_, err := tx.ExecContext(ctx,
		`UPDATE banner
				SET updated_at = $1, is_active = $2, active_version = $3, last_version = $4
				WHERE id = $5`, b.UpdatedAt, b.IsActive, activeVersion, b.Version, b.ID)
	if err != nil {
		return errs.New("fail to exec query: Update")
	}

	return nil
}

func (br *BannerRepo) UpdateBannerContent(tx *sql.Tx, b *models.Banner) error {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	contentJSON, err := json.Marshal(b.Content)
	if err != nil {
		return errs.WithMessagef(err, "fail to marshal content to JSON, content: %v", b.Content)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO banner_content (banner_id, version, content, updated_at)
				VALUES ($1, $2, $3, $4)`, b.ID, b.Version, contentJSON, b.UpdatedAt)
	if err != nil {
		return errs.New("fail to exec query: UpdateBannerContent")
	}

	return nil
}

func (br *BannerRepo) UpdateFeatureTag(tx *sql.Tx, b *models.Banner) error {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	for _, tagID := range b.TagIDs {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO banner_feature_tag (banner_id, feature_id, tag_id, version, updated_at) 
			VALUES ($1, $2, $3, $4, $5)`, b.ID, b.FeatureID, tagID, b.Version, b.UpdatedAt)
		if err != nil {
			return errs.WithMessagef(err, "fail to insert tag %d for banner %d", tagID, b.ID)
		}
	}

	return nil
}

func (br *BannerRepo) GetLastVersion(tx *sql.Tx, b *models.Banner) (error, int) {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	var lastVersion int
	err := tx.QueryRowContext(ctx,
		`SELECT last_version FROM banner WHERE id = $1`, b.ID).Scan(&lastVersion)
	if err != nil {
		return errs.WithMessagef(err, "failed to get last version for banner %d", b.ID), 0
	}

	return nil, lastVersion
}

func (br *BannerRepo) MergeUpdateVersion(tx *sql.Tx, b *models.Banner, lastVersion int) (*models.Banner, error) {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	var contentJSON []byte
	var featureID int
	err := tx.QueryRowContext(ctx,
		`SELECT bc.content, bft.feature_id
		FROM banner b 
		JOIN banner_content bc on b.id = bc.banner_id
		JOIN banner_feature_tag bft on bc.banner_id = bft.banner_id
		WHERE b.id = $1 
		AND bc.version = $2
		AND bft.version = $3`, b.ID, lastVersion, lastVersion).Scan(&contentJSON, &featureID)
	if err != nil {
		return nil, errs.WithMessagef(err, "fail to get old version for banner %d", b.ID)
	}

	var oldContent models.Content
	err = json.Unmarshal(contentJSON, &oldContent)
	if err != nil {
		return nil, errs.WithMessagef(err, "fail to unmarshal old content for banner %d", b.ID)
	}

	var oldTags []int
	rows, err := tx.QueryContext(ctx,
		`SELECT bft.tag_id
		FROM banner b 
		JOIN banner_content bc ON b.id = bc.banner_id
		JOIN banner_feature_tag bft ON bc.banner_id = bft.banner_id
		WHERE b.id = $1
		AND bft.version = $2`, b.ID, lastVersion)
	if err != nil {
		return nil, errs.WithMessagef(err, "failed to get old version for banner %d", b.ID)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var tagID int
		if err = rows.Scan(&tagID); err != nil {
			return nil, errs.WithMessagef(err, "failed to scan tag ID")
		}
		oldTags = append(oldTags, tagID)
	}

	if b.Content.Text == "" {
		b.Content.Text = oldContent.Text
	}
	if b.Content.Title == "" {
		b.Content.Title = oldContent.Title
	}
	if b.Content.URL == "" {
		b.Content.URL = oldContent.URL
	}
	if b.FeatureID == 0 {
		b.FeatureID = featureID
	}
	if b.TagIDs == nil {
		b.TagIDs = oldTags
	}

	return b, nil
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

	err, lastVersion := br.GetLastVersion(tx, b)
	if err != nil {
		return errs.WithMessagef(err, "fail to get last version of banner %d", b.ID)
	}
	b.Version = lastVersion + 1

	b, err = br.MergeUpdateVersion(tx, b, lastVersion)
	if err != nil {
		return errs.WithMessagef(err, "fail to merge update banner %d", b.ID)
	}

	err = br.UpdateBanner(tx, b, lastVersion)
	if err != nil {
		return errs.WithMessagef(err, "fail to update banner with id %d", b.ID)
	}

	err = br.UpdateBannerContent(tx, b)
	if err != nil {
		return errs.WithMessagef(err, "fail to update banner with id %d", b.ID)
	}

	err = br.UpdateFeatureTag(tx, b)
	if err != nil {
		return errs.WithMessagef(err, "fail to insert new tags with id %d", b.ID)
	}

	if err = tx.Commit(); err != nil {
		logger.Errf("failed to commit transaction UpdateBanner: %s", err)
		return err
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
		return errs.WithMessagef(err, "fail to exec query: DeleteBanner")
	}

	return nil
}

func (br *BannerRepo) CheckTagFeatureOverlap(b *models.Banner) (int, error) {
	ctx, cancel := context.WithTimeout(br.Ctx, time.Second*5)
	defer cancel()

	var bannerID int

	err := br.data.Master().QueryRowContext(ctx,
		`
	SELECT b.id
	FROM banner b
	JOIN banner_feature_tag bft on b.id = bft.banner_id 
	WHERE b.is_active = true
	AND b.active_version = bft.version
	AND bft.tag_id = ANY($1)
	AND bft.feature_id = $2`, pq.Array(b.TagIDs), b.FeatureID).Scan(&bannerID)
	if err != nil {
		return 0, errs.WithMessagef(err, "active banner with this combination of tagIDs and featureID is not found")
	}

	return bannerID, nil
}
