package store

import (
	"context"
	"database/sql"
	"urlify/model"
)

type Visit interface {
	LogVisit(ctx context.Context, v model.Visit) error
	GetAnalytics(ctx context.Context, code string) ([]model.Visit, error)
}

type visitStore struct {
	db *sql.DB
}

func NewVisitStore(db *sql.DB) Visit {
	return &visitStore{db: db}
}

func (s *visitStore) LogVisit(ctx context.Context, v model.Visit) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO visits (code, timestamp, ip, country, browser, device) VALUES ($1, $2, $3, $4, $5, $6)`,
		v.Code, v.Timestamp, v.IP, v.Country, v.Browser, v.Device)
	return err
}

func (s *visitStore) GetAnalytics(ctx context.Context, code string) ([]model.Visit, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT timestamp, ip, country, browser, device FROM visits WHERE code = $1`, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visits []model.Visit
	for rows.Next() {
		var v model.Visit
		err = rows.Scan(&v.Timestamp, &v.IP, &v.Country, &v.Browser, &v.Device)
		if err != nil {
			return nil, err
		}
		v.Code = code
		visits = append(visits, v)
	}

	return visits, nil
}
