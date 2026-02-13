package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"
	"time"

	"gorm.io/gorm"
)

type ReportRepository struct {
	q *query.Query
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{
		q: query.Use(db),
	}
}

func (r *ReportRepository) Create(report *model.Report) error {
	return r.q.Report.Create(report)
}

func (r *ReportRepository) GetByID(id int64) (*model.Report, error) {
	rp := r.q.Report
	return rp.Where(rp.ID.Eq(id)).First()
}

type ReportFilter struct {
	Type      string
	StationID int64
	CreatedBy int64
	IsAdmin   bool
}

func (r *ReportRepository) List(filter ReportFilter, offset, limit int) ([]*model.Report, int64, error) {
	rp := r.q.Report
	db := rp.UnderlyingDB().Model(&model.Report{}).Where("deleted_at IS NULL")

	if filter.Type != "" {
		db = db.Where("type = ?", filter.Type)
	}

	if !filter.IsAdmin {
		db = db.Where("(station_id = ? OR created_by = ?)", filter.StationID, filter.CreatedBy)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reports []*model.Report
	if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *ReportRepository) Delete(id int64) error {
	rp := r.q.Report
	_, err := rp.Where(rp.ID.Eq(id)).Delete()
	return err
}

func (r *ReportRepository) DeleteOlderThan(days int) ([]string, error) {
	rp := r.q.Report
	cutoff := time.Now().AddDate(0, 0, -days)

	var filePaths []string
	db := rp.UnderlyingDB().Model(&model.Report{}).
		Where("created_at < ? AND deleted_at IS NULL", cutoff).
		Pluck("file_path", &filePaths)
	if db.Error != nil {
		return nil, db.Error
	}

	_, err := rp.Where(rp.CreatedAt.Lt(cutoff)).Delete()
	if err != nil {
		return nil, err
	}

	return filePaths, nil
}
