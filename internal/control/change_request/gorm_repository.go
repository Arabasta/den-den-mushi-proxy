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
			return nil, nil
		}
		r.log.Error("DB error while fetching change request", zap.String("id", num), zap.Error(err))
		return nil, err
	}

	return dto.FromModel(&m)
}

func (r *GormRepository) FindChangeRequests(filter filters.ListCR) ([]*dto.Record, error) {
	var models []dto.Model
	query := r.db.Model(&dto.Model{})

	if filter.TicketIDs != nil && len(*filter.TicketIDs) > 0 {
		query = query.Where("Ticketnumber IN ?", *filter.TicketIDs)
	}

	query = query.Where("State = ?", "Approved")

	if filter.ImplementorGroups != nil {
		for _, group := range *filter.ImplementorGroups {
			query = query.Or("Implementorgroups LIKE ?", "%"+group+"%")
		}
	}

	if filter.LOB != nil {
		query = query.Where("Lob = ?", *filter.LOB)
	}

	if filter.Country != nil {
		query = query.Where("CountryImpacted LIKE ?", "%"+*filter.Country+"%")
	}

	// note that this is stored as text in DB
	if filter.StartTime != nil {
		query = query.Where("ChangeSchedStartTime >= ?", filter.StartTime.Format("2006-01-02 15:04:05"))
	}

	if filter.EndTime != nil {
		query = query.Where("ChangeSchedEndTime <= ?", filter.EndTime.Format("2006-01-02 15:04:05"))
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}

	pageSize := filter.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	err := query.Find(&models).Error
	if err != nil {
		r.log.Error("DB error while fetching change requests", zap.Error(err))
		return nil, err
	}

	r.log.Debug("Fetched change requests", zap.Any("models", models))
	return dto.FromModels(models)
}
