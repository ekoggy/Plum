package postgre

import (
	"database/sql"
	"fmt"
)

type Record struct {
	Name   string
	Size   string
	Date   string
	Price  string
	Buy    string
	Source string
}

var Database *sql.DB

func IsExist(Name string) (bool, error) {
	query := fmt.Sprintf("SELECT \"Name\" FROM database_leaks where \"Name\" = $1 limit 1")
	row := Database.QueryRow(query, Name)
	var tmp interface{}
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

func Insert(Name, Size, Date, Price, Buy, Source string) (int64, error) {
	check, err := IsExist(Name)
	if err != nil {
		return 0, err
	}
	if check == true {
		return 0, nil
	}
	res, err := Database.Exec("INSERT INTO database_leaks VALUES ($1, $2,$3,$4,$5,$6)",
		Name, Size, Date, Price, Buy, Source)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
