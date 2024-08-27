package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/rpc"
	"strconv"

	lib "github.com/owner/root/lib/libs"

	_ "github.com/lib/pq"
	"github.com/spiral/goridge"
)

type App struct{}

type Result struct {
	Data string `json:"data"`
}

func initDBConnection(driver string, credentials string) (*sql.DB, error) {
	db, err := sql.Open(driver, credentials)
	return db, err
}
func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func sqlQuery(db *sql.DB, expression string, limit string, offset string) (int, *sql.Rows, error) {
	c, err := db.Query(fmt.Sprintf("SELECT COUNT(*) as count FROM (%s) as x", expression))
	if err != nil {
		fmt.Println(err)
	}
	stmt := expression
	count := checkCount(c)
	if limit != "0" {
		stmt += fmt.Sprintf(" LIMIT %s", limit)
	}
	if offset != "0" {
		stmt += fmt.Sprintf(" OFFSET %s", offset)
	}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM (%s) as x", stmt))
	if err != nil {
		fmt.Println(err)
	}
	return count, rows, err
}

func getColumnTypes(rows *sql.Rows) []map[string]string {
	types, err := rows.ColumnTypes()
	if err != nil {
		panic(err.Error())
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	data_types := make([]map[string]string, len(types))
	for i := 0; i < len(types); i++ {
		results := make(map[string]string)
		results["column_name"] = columns[i]
		results["column_type"] = types[i].DatabaseTypeName()
		data_types[i] = results
	}

	return data_types
}

func makeResponse(rows *sql.Rows, data_types []map[string]string, limit string) map[string]interface{} {
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]interface{}, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	lim, _ := strconv.Atoi(limit)
	c := 0
	results := make(map[string]interface{})
	data := make([]map[string]interface{}, lim)

	for rows.Next() {

		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		for i, value := range values {
			switch value.(type) {
			case nil:
				results[columns[i]] = nil

			case []byte:
				s := string(value.([]byte))
				x, err := strconv.Atoi(s)

				if err != nil {
					results[columns[i]] = s
				} else {
					results[columns[i]] = x
				}

			default:
				results[columns[i]] = value
			}
		}
		data[c] = results
		c++
	}

	response := make(map[string]interface{})
	response["columns"] = data_types
	response["rows"] = data
	return response
}

func (s *App) MakeRequest(args [5]string, responce *map[string]interface{}) error {
	token, err := lib.GetRequestToken()
	if err != nil {
		fmt.Println(err)
		// panic(err.Error())
	}
	db, err := initDBConnection(args[0], args[1])
	if err != nil {
		fmt.Println(err)
		// panic(err.Error())
	}
	count, rows, err := sqlQuery(db, args[2], args[3], args[4])
	if err != nil {
		fmt.Println(err)
		// panic(err.Error())
	}
	data := make(map[string]interface{})
	data_types := getColumnTypes(rows)
	if args[3] != "0" {
		data = makeResponse(rows, data_types, args[3])
		// maps.Copy(data, temp)
	} else {
		data = makeResponse(rows, data_types, strconv.Itoa(count))
		// maps.Copy(data, temp)
	}
	data["count"] = count
	data["token"] = token
	*responce = data
	defer rows.Close()
	defer db.Close()
	return nil
}

func main() {
	ln, err := net.Listen("tcp", ":6001")
	if err != nil {
		fmt.Println(err)
	}

	rpc.Register(new(App))

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go rpc.ServeCodec(goridge.NewCodec(conn))
	}
}
