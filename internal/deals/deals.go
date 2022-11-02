package deals

import (
	"database/sql"
	"fmt"
)

func Create(db *sql.DB, deal Deal, user string) error {

	insertQuery := `INSERT INTO deals (username, weight,volume,load_datetime,unload_datetime,load_adress,unload_adress,comment) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`

	res, err := db.Exec(insertQuery, user, deal.Weight, deal.Volume, deal.LoadDateTime, deal.UnLoadDateTime, deal.LoadAdress, deal.UnloadAdress, deal.Comment)
	if err != nil {
		fmt.Println(err)
		return err
	}
	rowsNum, _ := res.RowsAffected()
	fmt.Println("rows affected:", rowsNum)
	return nil
}

func List(db *sql.DB, user string, perPage, offset, page int) (list ListObj, err error) {

	selectQuery := `SELECT id,username, weight,volume,load_datetime,unload_datetime,load_adress,unload_adress,comment FROM deals WHERE username = $1 ORDER BY id ASC LIMIT $2 OFFSET $3`
	selectTotalCountQuery := `SELECT COUNT (*) FROM deals WHERE username = $1`

	rows, err := db.Query(selectQuery, user, perPage, offset)
	defer rows.Close()

	if err != nil {
		return ListObj{}, err
	}

	for rows.Next() {
		var deal Deal
		err = rows.Scan(&deal.Id, &deal.User, &deal.Weight, &deal.Volume, &deal.LoadDateTime, &deal.UnLoadDateTime,
			&deal.LoadAdress, &deal.UnloadAdress, &deal.Comment)
		if err != nil {
			return ListObj{}, err
		}
		list.Deals = append(list.Deals, deal)
	}

	err = db.QueryRow(selectTotalCountQuery, user).Scan(&list.Total)
	if err != nil {
		return ListObj{}, err
	}

	list.Page = page
	list.PerPage = perPage

	return list, nil
}

func Info(db *sql.DB, username, dealId string) (deal Deal, err error) {

	selectQuery := `SELECT id,username, weight,volume,load_datetime,unload_datetime,load_adress,unload_adress,comment FROM deals WHERE username = $1 AND id = $2 LIMIT 1`

	err = db.QueryRow(selectQuery, username, dealId).Scan(&deal.Id, &deal.User, &deal.Weight, &deal.Volume, &deal.LoadDateTime, &deal.UnLoadDateTime,
		&deal.LoadAdress, &deal.UnloadAdress, &deal.Comment)

	return deal, err
}
