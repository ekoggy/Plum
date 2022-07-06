package main

func insert(Name, Size, Date, Price, Buy, Source string) (int64, error) {
	res, err := Database.Exec("INSERT INTO database_leaks VALUES ($1, $2,$3,$4,$5,$6)",
		Name, Size, Date, Price, Buy, Source)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func remove(addr string) error {
	stmt, err := Database.Prepare("DELETE FROM database_leaks WHERE address=$1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(addr)
	if err != nil {
		return err
	}
	return nil
}

func show(arg string) ([]Record, error) {
	rows, err := Database.Query("SELECT * FROM database_leaks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rs = make([]Record, 0)
	var rec Record
	for rows.Next() {
		err = rows.Scan(&rec.Name, &rec.Size,
			&rec.Date, &rec.Price, &rec.Buy, &rec.Source)
		if err != nil {
			return nil, err
		}
		rs = append(rs, rec)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return rs, nil
}
