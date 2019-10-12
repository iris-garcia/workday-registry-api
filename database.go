package workday

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB returns a sql.DB connection using a given DBConfig.
func ConnectDB(config DBConfig) (*sql.DB, error) {
	ds := fmt.Sprintf("%v:%v@tcp(%v)/%v", config.User, config.Password, config.Host, config.Database)
	db, err := sql.Open("mysql", ds)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func insertNewEmployee(db *sql.DB, employee Employee) (id uint, rows uint, err error) {
	result, err := db.Exec(
		"INSERT INTO employee(firstname, lastname, role, password) values('?', '?', ?, '?')",
		employee.Firstname, employee.Lastname, employee.Role, employee.Password,
	)
	if err != nil {
		return 0, 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	return uint(lastID), uint(rowsAffected), nil
}

func findAllEmployees(db *sql.DB) ([]Employee, error) {
	employees := make([]Employee, 0)
	rows, err := db.Query("SELECT * FROM employee")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.ID, &employee.Firstname, &employee.Lastname, &employee.Role, &employee.Password)
		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}