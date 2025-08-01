package change_request

import (
	"den-den-mushi-Go/internal/control/filters"
	dto "den-den-mushi-Go/pkg/dto/change_request"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewGormRepository(db *gorm.DB, log *zap.Logger) *GormRepository {
	return &GormRepository{
		db:  db,
		log: log,
	}
}

func (r *GormRepository) FindByTicketNumber(num string) (*dto.Record, error) {
	var m dto.Model
	err := r.db.Where("TicketNumber = ?", num).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No change request found for ID", zap.String("id", num))
			// must return err
			return nil, errors.New(num + " not found")
		}
		r.log.Error("DB error while fetching change request", zap.String("id", num), zap.Error(err))
		return nil, err
	}

	return dto.FromModel(&m)
}

func (r *GormRepository) FindChangeRequestsByFilter(f filters.ListCR) ([]*dto.Record, error) {
	var models []dto.Model
	query := r.db.Model(&dto.Model{})

	if f.TicketIDs != nil && len(*f.TicketIDs) > 0 {
		query = query.Where("TicketNumber IN ?", *f.TicketIDs)
	}

	query = query.Where("State = ?", "Approved")

	if f.ImplementorGroups != nil && len(*f.ImplementorGroups) > 0 {
		query = query.Where("ImplementerGroup IN ?", *f.ImplementorGroups)
	}

	if f.LOB != nil {
		query = query.Where("LOB = ?", *f.LOB)
	}

	if f.Country != nil {
		query = query.Where("CountryImpacted LIKE ?", "%"+*f.Country+"%")
	}

	// note that this is stored as text in DB
	if f.StartTime != nil {
		query = query.Where("ChangeSchedStartDateTime >= ?", f.StartTime.Format("2006-01-02 15:04:05"))
	}

	if f.EndTime != nil {
		query = query.Where("ChangeSchedEndDateTime <= ?", f.EndTime.Format("2006-01-02 15:04:05"))
	}

	query = query.Order("ChangeSchedStartDateTime ASC")

	page := f.Page
	if page < 1 {
		page = 1
	}

	pageSize := f.PageSize
	if pageSize <= 0 || pageSize > 100000 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	err := query.Find(&models).Error
	if err != nil {
		r.log.Error("DB error while fetching change requests", zap.Error(err))
		return nil, err
	}
	if len(models) == 0 {
		return []*dto.Record{}, nil
	}

	r.log.Debug("Fetched change requests", zap.Int("Count", len(models)))
	return dto.FromModels(models)
}
