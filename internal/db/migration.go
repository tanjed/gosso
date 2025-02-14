package db

var migrationQueries = make(map[string]string)

func createClientsTable() {
	migrationQueries["clients_table_create"] = `CREATE TABLE IF NOT EXISTS clients (
    client_id UUID PRIMARY KEY,
    client_name TEXT,
    client_secret TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
	`
    migrationQueries["client_table_client_name_index_create"] = `CREATE INDEX IF NOT EXISTS client_name_idx ON clients(client_name);`
}

func createUsersTable() {
	migrationQueries["user_table_create"] = `CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    mobile_number TEXT,
    password TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
`
    migrationQueries["user_table_client_name_index_create"] = `CREATE INDEX IF NOT EXISTS mobile_number_idx ON users(mobile_number);`
}

func createClientAuthorizationCodesTable() {
    migrationQueries["client_authorization_codes_table_create"] = `CREATE TABLE IF NOT EXISTS client_authorization_codes (
        client_id UUID,
        client_code TEXT,
        redirect_uri TEXT,
        generated_at TIMESTAMP,
        expired_at TIMESTAMP,
        PRIMARY KEY((client_id), redirect_uri, expired_at)
    );
    `
}

func RegisterMigrationQueries() map[string]string {
	createClientsTable()
    createUsersTable()
    createClientAuthorizationCodesTable()

	return migrationQueries
}
