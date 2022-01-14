package Postgres

import (
	"fmt"
	"reflect"
)

import "database/sql"
import _ "github.com/lib/pq"
import "strconv"

var db *sql.DB

type configuration struct {
	postgres_host string
	postgres_db_name string
	postgres_user string
	postgres_password string
	postgres_ssl string
}
var config configuration
var tableColumnOrders map[string][]string

func Init(conf map[string]string) {
	config.postgres_host = conf["postgres_host"]
	config.postgres_db_name = conf["postgres_db_name"]
	config.postgres_user = conf["postgres_user"]
	config.postgres_password = conf["postgres_password"]
	config.postgres_ssl = conf["postgres_ssl"]
	tableColumnOrders = make(map[string][]string)
//	fmt.Println("tablesMap: ", tableColumnOrders)
}

var connect_string string

func Connect() error {
	var err error
	connect_string = "host=" + config.postgres_host + " user=" + config.postgres_user + " password=" + config.postgres_password +
			" dbname=" + config.postgres_db_name + " sslmode=" + config.postgres_ssl
	db, err = sql.Open("postgres", connect_string)
//	db, err := sql.Open("postgres", "host=postgres user=warehouse password=warehouse1 dbname=warehouse sslmode=disable")
	if err != nil {
		fmt.Println("ERR: ", err)
		return err
	}
	err = Ping()
	return err
}

func Ping() error {
	err := db.Ping()
	if err != nil {
		fmt.Println("ERR: ", err)
		return err
	}
//	fmt.Println("Postgres successfully connected!")
	return nil
}
//
//

func TableCreate_(tableName string, keysOrder []string, fields map[string]string, primaryKey string) error {
	fields_cnt := len(keysOrder)
	tableColumnOrders[tableName] = make([]string, len(keysOrder))
	copy(tableColumnOrders[tableName], keysOrder)
	cnt := 0
	addFlag := ""
	query := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	for _, key := range keysOrder {
		if key == primaryKey {
			addFlag = " PRIMARY KEY"
		} else {
			addFlag = ""
		}
		cnt++
		query += key + " " + fields[key] + addFlag
		if cnt != fields_cnt {
			query += ", "
		}
	}
	query += ")"

	//fmt.Println(query)
	//*
	_, err := db.Exec(query)
	if err != nil{
		fmt.Println("ERR: ", err)
	}
	return err
	//*/
	return nil
}

func TableCreate__(tableName string, keysOrder []string, fields map[string]string, primaryKey string, foreignKey []string) error {
	fields_cnt := len(keysOrder)
	tableColumnOrders[tableName] = make([]string, len(keysOrder))
	copy(tableColumnOrders[tableName], keysOrder)
	cnt := 0
	addFlag := ""
	query := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	for _, key := range keysOrder {
		if key == primaryKey {
			addFlag = " PRIMARY KEY"
		} else {
			addFlag = ""
		}
		cnt++
		query += key + " " + fields[key] + addFlag
		if cnt != fields_cnt {
			query += ", "
		}
	}

	for _, fKey := range foreignKey {
		query += ", " + fKey
	}
	query += ")"

	//fmt.Println(query)
	//*
	_, err := db.Exec(query)
	if err != nil{
		fmt.Println("ERR: ", err)
	}
	return err
	//*/
	return nil
}

func AlterTable_AddColumn(tableName string, columnName string, columnType string, constraint string) error { // To-Do: handle constraint
	query := "ALTER TABLE " + tableName + " ADD COLUMN " + columnName + " " + columnType + ";"
	fmt.Println("query: ", query)
	_, err := db.Exec(query)
	if err != nil{
		fmt.Println("ERR: ", err)
	}
	return err

}

func AlterTable_AddConstraint(tableName string, columnName string, constraint string) error {
	query := "ALTER TABLE " + tableName + " ADD CONSTRAINT " + columnName + " " + constraint + ";"
	fmt.Println("query: ", query)
	_, err := db.Exec(query)
	if err != nil{
		fmt.Println("ERR: ", err)
	}
	return err

}

// ALTER TABLE distributors ADD CONSTRAINT distfk FOREIGN KEY (address) REFERENCES addresses (address) MATCH FULL;

func TableCreateBySnapshot(tableName string, keysOrder []string, example map[string]interface{}, primaryKey string) error {
//	fmt.Println("=================")
	fieldsCnt := len(keysOrder)
	tableColumnOrders[tableName] = make([]string, len(keysOrder))
	copy(tableColumnOrders[tableName], keysOrder)
//	fmt.Println("tablesMap: ", tableColumnOrders)
	cnt := 0
	//*
	query := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	addFlag := ""
	for _, key := range keysOrder {
		if key == primaryKey {
			addFlag = " PRIMARY KEY"
		} else {
			addFlag = ""
		}
		cnt++
		switch example[key].(type) {
		case string: 	query += key + " " + "text" + addFlag
		case float64:	query += key + " " + "real" + addFlag
		case float32:	query += key + " " + "real" + addFlag
		case int:		query += key + " "  + "integer" + addFlag
		default: fmt.Println("val: ", reflect.TypeOf(example[key]))
		}
		if cnt != fieldsCnt {
			query += ", "
		}
	}

	query += ")"
	//fmt.Println( query)
//*
	_, err := db.Exec(query)
	if err != nil{
		fmt.Println("ERR: ", err)
		return err
	}
	return err
//*/
	return nil
}

func TableCreate(name string, fields map[string]string) error {
	fields_cnt := len(fields)
	cnt := 0
	query := "CREATE TABLE IF NOT EXISTS " + name + " ("
	for key, val := range fields {
		cnt++
		query += key + " " + val 
		if cnt != fields_cnt {
			query += ", "
		}
	}
	query += ")"

	//fmt.Println(query)
	_, err := db.Exec(query)
    	if err != nil{
        	fmt.Println("ERR: ", err)
    	}
	return err
}

func TableDelete(name string) error {
	query := "DROP TABLE " + name
	//fmt.Println(query)
	_, err := db.Exec(query)
    	if err != nil{
        	fmt.Println("ERR: ", err)
    	}
	return err
}

func TableUpdateRow(tableName string, row map[string]interface{}, condition string) error {
	query1 := "UPDATE " + tableName + " SET "
	cnt := 0
	for key, val := range row {
		if cnt != 0 {
			query1 += ", "
		}
		cnt++
		query1 += key + " = '" + val.(string) + "'"
	}
	query1 += " WHERE " + condition
	fmt.Println(query1)
	_, err := db.Exec(query1)
	if err != nil{
		fmt.Println("ERR: ", err)
	}
	return err
}

func TableInsertRow(tableName string, row map[string]interface{}) error{
	columnsOrder := tableColumnOrders[tableName]
	//	fmt.Println("row: ", row)
	// fields_cnt := len(columnsOrder)
	cnt := 0
	query1 := "INSERT INTO " + tableName + " ("
	query2 := " VALUES ("
	for _, key := range columnsOrder {
		_, ok := row[key]
		if !ok {
			continue
		}
		if cnt != 0 {
			query1 += ", "
			query2 += ", "
		}
		cnt++
		query1 += key
		val := row[key]
		switch typ := val.(type) {
		default:
			fmt.Printf("unexpected type '%T' for key '%s'. Fill with empty string\r\n", typ, key)
		case nil:
//			fmt.Printf("INFO: empty value for key '%s'. Fill with empty string\r\n", key)
			query2 += "' '"
		case string:
			query2 += "'" + val.(string) + "'"
		case int:
			query2 += strconv.Itoa( val.(int) )
		case int16:
			query2 += strconv.FormatUint( uint64(val.(int16)), 10 )
		case int32:
			query2 += strconv.FormatUint( uint64(val.(int32)), 10 )
		case int64:
			query2 += strconv.FormatUint( uint64(val.(int64)), 10 )
		case int8:
			query2 += strconv.FormatUint( uint64(val.(int8)), 10 )
		case uint:
			query2 += strconv.FormatUint( uint64(val.(uint)), 10 )
		case uint8:
			query2 += strconv.FormatUint( uint64(val.(uint8)), 10 )
		case uint16:
			query2 += strconv.FormatUint( uint64(val.(uint16)), 10 )
		case uint32:
			query2 += strconv.FormatUint( uint64(val.(uint32)), 10 )
		case uint64:
			query2 += strconv.FormatUint( val.(uint64), 10 )
		case float64:
			query2 += fmt.Sprintf("%f", val)
		case float32:
			query2 += fmt.Sprintf("%f", val)
		}
		//		if cnt != fields_cnt {
		//			query1 += ", "
		//			query2 += ", "
		//		}
	}

	query1 += ") "
	query2 += ")"
	query := query1 + query2
	//fmt.Println("TableInsertRow query: ", query)

	//	query := "INSERT INTO components (id, description) VALUES (2, 'apple')"
	//*
	_, err := db.Exec(query)
	if err != nil{
		fmt.Println("ERR: ", err)
	}
	return err
	//*/
	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TableGetRow(tableName string, columnNames string, condition string) map[string]interface{}{
	rows, err := db.Query("SELECT " + columnNames + " FROM " + tableName + " " + condition)
	defer rows.Close()
	check(err)
	cols, err := rows.Columns()
	check(err)


	if rows.Next() {
		data := make(map[string]interface{})
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range cols {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)
		for i, colName := range cols {
			data[colName] = columns[i]
		}
		//fmt.Println("row: ", data)
		return data
	}
	return nil
}

func TableGetRows(tableName string, columnNames string, condition string) []map[string]interface{}{
//	fmt.Println("QUERY: ", "SELECT " + columnNames + " FROM " + tableName + " " + condition)
	rows, err := db.Query("SELECT " + columnNames + " FROM " + tableName + " " + condition)
	check(err)
	cols, err := rows.Columns()
	check(err)

	var data_ []map[string]interface{}

	for rows.Next() {
		data := make(map[string]interface{})
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)
		for i, colName := range cols {
			data[colName] = columns[i]
		}
		//fmt.Println("row: ", data)
		data_ = append(data_, data)
//		return data
	}
	return data_
}

func TableGetRow_map(tableName string, columnNames string) map[string]interface{}{
	rows, err := db.Query("SELECT " + columnNames + " FROM " + tableName + " limit 1")
	check(err)
	cols, err := rows.Columns()
	check(err)

	data := make(map[string]interface{})

	if rows.Next() {
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		
		rows.Scan(columnPointers...)
		for i, colName := range cols {
			data[colName] = columns[i]
		}
		return data
	}
	return nil
}

func TableGetRow_arr(tableName string, columnNames string) []interface{}{
	rows, err := db.Query("SELECT " + columnNames + " FROM " + tableName + " limit 1")
	check(err)
	cols, err := rows.Columns()
	check(err)

	data := make([]interface{}, len(cols))

	if rows.Next() {
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		
		rows.Scan(columnPointers...)
		for i, _ := range cols {
			data[i] = columns[i]
		}
		return data
	}
	return nil
}


