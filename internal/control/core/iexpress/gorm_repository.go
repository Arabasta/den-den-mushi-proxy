package iexpress

import (
	"den-den-mushi-Go/internal/control/filters"
	dto "den-den-mushi-Go/pkg/dto/iexpress"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
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
	err := r.db.Where("Ticket = ?", num).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No ticket found for ID", zap.String("id", num))
			// must return err
			return nil, errors.New(num + " not found")
		}
		r.log.Error("DB error while fetching iexpress request", zap.String("id", num), zap.Error(err))
		return nil, err
	}

	return dto.FromModel(&m)
}

func (r *GormRepository) FindApprovedByFilter(f filters.ListIexpress) ([]*dto.Record, error) {
	var models []dto.Model
	query := r.db.Model(&dto.Model{})

	if f.TicketIDs != nil && len(*f.TicketIDs) > 0 {
		query = query.Where("Ticket IN ?", *f.TicketIDs)
	}

	query = query.Where("State = ?", "Approved")

	if f.ImplementorGroups != nil && len(*f.ImplementorGroups) > 0 {
		query = query.Where(
			"("+
				"COALESCE(Approver_Group_1, '') IN ? OR "+
				"COALESCE(Approver_Group_2, '') IN ? OR "+
				"COALESCE(MD_Approver_Group, '') IN ?"+
				")",
			*f.ImplementorGroups, *f.ImplementorGroups, *f.ImplementorGroups,
		)
	}

	if f.AppImpacted != nil && len(*f.AppImpacted) > 0 {
		var likeClauses []string
		var args []interface{}
		for _, g := range *f.AppImpacted {
			likeClauses = append(likeClauses, "Application_Impacted LIKE ?")
			args = append(args, "%"+g+"%")
		}
		query = query.Where("("+strings.Join(likeClauses, " OR ")+")", args...)
	}

	if f.Requestor != nil {
		query = query.Where("Requestor = ?", *f.Requestor)
	}

	if f.LOB != nil {
		query = query.Where("Unit = ?", *f.LOB)
	}

	if f.Country != nil {
		query = query.Where("Country_of_Origin = ?", *f.Country)
	}

	// note that this is stored as text in DB
	if f.StartTime != nil && f.EndTime != nil {
		query = query.Where(
			"Schedule_Start <= ? AND Schedule_End >= ?",
			f.EndTime.Format("2006-01-02 15:04:05"),
			f.StartTime.Format("2006-01-02 15:04:05"),
		)
	} else if f.StartTime != nil {
		query = query.Where("Schedule_End >= ?", f.StartTime.Format("2006-01-02 15:04:05"))
	} else if f.EndTime != nil {
		query = query.Where("Schedule_Start <= ?", f.EndTime.Format("2006-01-02 15:04:05"))
	}

	query = query.Order("Schedule_Start ASC")

	page := f.Page
	if page < 1 {
		page = 1
	}

	pageSize := f.PageSize
	if pageSize <= 0 || pageSize > 1000 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	err := query.Find(&models).Error
	if err != nil {
		r.log.Error("DB error while fetching iexpress requests", zap.Error(err))
		return nil, err
	}
	if len(models) == 0 {
		return []*dto.Record{}, nil
	}

	r.log.Debug("Fetched iexpress requests", zap.Int("Count", len(models)))
	return dto.FromModels(models)
}

func (r *GormRepository) CountApprovedByFilter(f filters.ListIexpress) (int, error) {
	var count int64
	query := r.db.Model(&dto.Model{})

	if f.TicketIDs != nil && len(*f.TicketIDs) > 0 {
		query = query.Where("Ticket IN ?", *f.TicketIDs)
	}

	query = query.Where("State = ?", "Approved")

	if f.ImplementorGroups != nil && len(*f.ImplementorGroups) > 0 {
		query = query.Where(
			"("+
				"COALESCE(Approver_Group_1, '') IN ? OR "+
				"COALESCE(Approver_Group_2, '') IN ? OR "+
				"COALESCE(MD_Approver_Group, '') IN ?"+
				")",
			*f.ImplementorGroups, *f.ImplementorGroups, *f.ImplementorGroups,
		)
	}

	if f.AppImpacted != nil && len(*f.AppImpacted) > 0 {
		var likeClauses []string
		var args []interface{}
		for _, g := range *f.AppImpacted {
			likeClauses = append(likeClauses, "Application_Impacted LIKE ?")
			args = append(args, "%"+g+"%")
		}
		query = query.Where("("+strings.Join(likeClauses, " OR ")+")", args...)
	}

	if f.LOB != nil {
		query = query.Where("Unit = ?", *f.LOB)
	}

	if f.Country != nil {
		query = query.Where("Country_of_Origin = ?", *f.Country)
	}

	// note that this is stored as text in DB
	if f.StartTime != nil && f.EndTime != nil {
		query = query.Where(
			"Schedule_Start <= ? AND Schedule_End >= ?",
			f.EndTime.Format("2006-01-02 15:04:05"),
			f.StartTime.Format("2006-01-02 15:04:05"),
		)
	} else if f.StartTime != nil {
		query = query.Where("Schedule_End >= ?", f.StartTime.Format("2006-01-02 15:04:05"))
	} else if f.EndTime != nil {
		query = query.Where("Schedule_Start <= ?", f.EndTime.Format("2006-01-02 15:04:05"))
	}

	err := query.Count(&count).Error
	if err != nil {
		r.log.Error("DB error while counting change requests", zap.Error(err))
		return 0, err
	}

	r.log.Debug("Counted iexpress requests", zap.Int64("Count", count))
	return int(count), nil
}
