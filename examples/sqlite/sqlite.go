package main

import (
	"os"
	"fmt"
	"errors"
	"github.com/rpucella/go-neutralino-extension"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	connInfo, err := neutralinoext.ReadConnInfo(os.Stdin)
	if err != nil {
		panic(err)
	}

	if err := connInfo.StartMessageLoop(processMsg); err != nil {
		panic(err)
	}
}

func extractStrings(m map[string]any, keys []string) ([]string, error) {
	result := make([]string, len(keys))
	for i, key := range keys {
		ifc, ok := m[key]
		if !ok {
			return nil, fmt.Errorf("no key %s", key)
		}
		s, ok := ifc.(string)
		if !ok {
			return nil, fmt.Errorf("not a string %v", ifc)
		}
		result[i] = s
	}
	return result, nil
}

func processMsg(event string, data any) (map[string]any, error) {
	if event != "query" {
		return nil, nil
	}
	dataObj, ok := data.(map[string]any)
	if !ok {
		return nil, errors.New("data not an object")
	}
	params, err := extractStrings(dataObj, []string{"database", "query"})
	if err != nil {
		return nil, err
	}
	rows, err := queryDatabase(params[0], params[1])
	if err != nil {
		return nil, err
	}
	result := make(map[string]any)
	result["rows"] = rows
	return result, nil
}

func queryDatabase(database string, query string) ([][]any, error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, fmt.Errorf("cannot open db file: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("cannot query: %w", err)
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("cannot get columns: %w", err)
	}
	result := make([][]any, 0)
	for rows.Next() {
		row := make([]any, len(cols))
		rowAddrs := make([]any, len(cols))
		for i := range row {
			rowAddrs[i] = &row[i]
		}
		err = rows.Scan(rowAddrs...)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		result = append(result, row)
	}
	return result, nil
}
