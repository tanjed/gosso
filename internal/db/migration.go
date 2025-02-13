package db

var migrationQueries = make(map[string]string)

func createClientsTable() {
	migrationQueries["clients_create"] = `CREATE TABLE IF NOT EXISTS clients (
		client_id BIGSERIAL PRIMARY KEY,
		client_name VARCHAR(255) NOT NULL,
		client_secret VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
}


func RegisterMigrationQueries() map[string]string {
	createClientsTable()

	return migrationQueries
}
