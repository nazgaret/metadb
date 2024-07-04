package catalog

import (
	"context"
	"fmt"

	"github.com/metadb-project/metadb/cmd/metadb/dbx"
)

func Origins(db *dbx.DB) ([]string, error) {
	dc, err := db.Connect()
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	defer dbx.Close(dc)

	//var tenants string
	//err = dc.QueryRow(context.TODO(), "SELECT name FROM metadb.origin").Scan(&tenants)
	//switch {
	//case err == sql.ErrNoRows:
	//	return "", fmt.Errorf("select tenants: %w", err)
	//case err != nil:
	//	return "", err
	//default:
	//	return tenants, nil
	//}

	rows, err := dc.Query(context.TODO(), "SELECT name FROM metadb.origin")
	if err != nil {
		return nil, fmt.Errorf("selecting origins: %w", err)
	}
	defer rows.Close()
	origins := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("reading origin: %w", err)
		}
		origins = append(origins, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("last error check: %w", err)
	}
	return origins, nil
}
