package postgre

func IsExist(Name string) (bool, error) {
	response, err := Database.Query("SELECT EXISTS (SELECT * FROM database_leaks WHERE \"Name\" = $1)",
		Name)
	if err != nil {
		return true, err
	}

	var result bool
	err = response.Scan(result)

	if err != nil {
		return result, err
	}

	return result, nil

}

func insert(Name, Size, Date, Price, Buy, Source string) (int64, error) {
	//check, err := IsExist(Name)
	//if err != nil {
	//	return 0, err
	//}
	//if check == true {
	//	return 0, nil
	//}
	res, err := Database.Exec("INSERT INTO database_leaks VALUES ($1, $2,$3,$4,$5,$6)",
		Name, Size, Date, Price, Buy, Source)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
