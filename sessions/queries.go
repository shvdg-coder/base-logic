package sessions

const createSessionsTableQuery = `
	CREATE TABLE IF NOT EXISTS sessions  (
		token TEXT PRIMARY KEY,
		data BYTEA NOT NULL,
		expiry TIMESTAMPTZ NOT NULL
	);
	`

const createSessionExpiryIndexQuery = `
	DO $$ BEGIN
		IF NOT EXISTS (
			SELECT 1
			FROM   pg_class c
			JOIN   pg_namespace n ON n.oid = c.relnamespace
			WHERE  c.relname = 'sessions_expiry_idx'
			AND    n.nspname = 'public' 
		) THEN
			CREATE INDEX sessions_expiry_idx ON sessions (expiry);
		END IF;
	END $$`
