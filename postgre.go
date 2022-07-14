package main

func IsExist(Name string) (bool, error) {
	response, err := Database.Exec("SELECT EXISTS (SELECT * FROM database_leaks WHERE \"Name\" = $1)",
		Name)
	if err != nil {
		return false, err
	}

	result, err := response.RowsAffected()

	if err != nil {
		return false, err
	}

	if result == 1 {
		return true, nil
	}
	return false, nil

}

func insert(Name, Size, Date, Price, Buy, Source string) (int64, error) {
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
