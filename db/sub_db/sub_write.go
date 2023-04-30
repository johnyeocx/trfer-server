package sub_db

import (
	"database/sql"

	"github.com/johnyeocx/usual/server2/db/models/sub_models"
	"github.com/johnyeocx/usual/server2/utils/enums"
)

type SubDB struct {
	DB	*sql.DB
}

func (s *SubDB) InsertSub(cusId int, brandId int, status enums.SubStatus, plan *sub_models.Plan) (*int, error){
	query := `INSERT into subscription 
	(customer_id, brand_id, status, unit_amount, interval, interval_count, currency, bill_date) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING sub_id`

	var subId int
	err := s.DB.QueryRow(
		query, 
		cusId, 
		brandId, 
		status,
		plan.UnitAmount, 
		plan.Interval, 
		plan.IntervalCount, 
		plan.Currency,
		plan.BillDate,
	).Scan(&subId)

	if err != nil {
		return nil, err
	}
	return &subId, nil
}