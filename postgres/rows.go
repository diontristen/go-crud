package postgres

import (
	"fmt"

	"github.com/diontristen/go-crud/util"
	"github.com/jmoiron/sqlx"
)

func ForRows(query string, args interface{}, db sqlx.Ext, dest interface{}, onRow func() (bool, error)) error {
	rows, err := sqlx.NamedQuery(db, query, args)
	if err != nil {
		util.Errorf("%s:\n%s:\n%v", err, query, args)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(dest)
		if err != nil {
			err = fmt.Errorf("database rows iteration error: %s", err)
			util.Error(err)
			return err
		}

		goON, err := onRow()
		if err != nil {
			return err
		}

		if !goON {
			break
		}
	}

	if err := rows.Err(); err != nil {
		util.Error(err)
		return err
	}

	return nil

}
