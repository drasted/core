package dwh

// SQL statements for SQLite backend.
var (
	createTableOrdersSQLite = `
	CREATE TABLE IF NOT EXISTS orders (
		id					TEXT PRIMARY KEY,
		type				UNSIGNED INTEGER NOT NULL,
		author				TEXT NOT NULL,
		counter_agent 		TEXT NOT NULL,
		duration 			UNSIGNED INTEGER NOT NULL,
		price				BIGINT NOT NULL
	);`
	createTableDealsSQLite = `
	CREATE TABLE IF NOT EXISTS deals (
		id					TEXT PRIMARY KEY,
		status				UNSIGNED INTEGER NOT NULL,
		supplier 			TEXT NOT NULL,
		consumer 			TEXT NOT NULL,
		duration 			UNSIGNED INTEGER NOT NULL,
		price				BIGINT NOT NULL,
		startTime			UNSIGNED INTEGER NOT NULL
	);`
	createTableChangesSQLite = `
	CREATE TABLE IF NOT EXISTS changes (
		duration 			UNSIGNED INTEGER NOT NULL,
		price				BIGINT NOT NULL,
		deal				BIGINT NOT NULL,
		FOREIGN KEY (deal)	REFERENCES deals(TEXT) ON DELETE CASCADE
	);`
)
