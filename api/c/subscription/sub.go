package subscription

import (
	"database/sql"
	"strings"

	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/db/models/sub_models"
	"github.com/johnyeocx/usual/server2/db/sub_db"
	"github.com/johnyeocx/usual/server2/errors/sub_errors"
	"github.com/johnyeocx/usual/server2/utils/enums"
)

func CreateSub(sqlDB *sql.DB, cusId int, brandId int, status enums.SubStatus, plan sub_models.Plan) (*int, *models.RequestError) {
	s := sub_db.SubDB{DB: sqlDB}

	sId, err := s.InsertSub(cusId, brandId, status, &plan)
	if err != nil {
		if strings.Contains(err.Error(), `"subscription" violates foreign key constraint`) {
			return nil, sub_errors.InvalidBrandIDErr(err)
		} else if strings.Contains(err.Error(), `unique constraint "subscription_brand_id_customer_id_key"`) {
			return nil, sub_errors.DuplicateSubErr(err)
		}

		return nil, sub_errors.CreateSubFailedErr(err)
	}

	return sId, nil
}