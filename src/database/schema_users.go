package database

// usersCreateStatements contains all CREATE TABLE statements for users.db
// Per AI.md PART 10: All tables use CREATE TABLE IF NOT EXISTS (idempotent)
var usersCreateStatements = []string{
	// Admins - server admin accounts
	`CREATE TABLE IF NOT EXISTS admins (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		username    TEXT NOT NULL UNIQUE,
		password    TEXT NOT NULL,
		email       TEXT,
		role        TEXT NOT NULL DEFAULT 'admin',
		enabled     INTEGER NOT NULL DEFAULT 1,
		api_token_hash TEXT,
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		last_login  INTEGER,
		failed_attempts INTEGER NOT NULL DEFAULT 0,
		locked_until INTEGER,
		source      TEXT NOT NULL DEFAULT 'local',
		external_id TEXT,
		groups      TEXT,
		last_sync   INTEGER
	)`,

	`CREATE UNIQUE INDEX IF NOT EXISTS idx_admins_username ON admins(username)`,

	// Admin preferences
	`CREATE TABLE IF NOT EXISTS admin_preferences (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		admin_id    INTEGER NOT NULL UNIQUE,
		theme           TEXT NOT NULL DEFAULT 'dark',
		font_size       TEXT NOT NULL DEFAULT 'medium',
		reduce_motion   INTEGER NOT NULL DEFAULT 0,
		date_format     TEXT NOT NULL DEFAULT 'YYYY-MM-DD',
		time_format     TEXT NOT NULL DEFAULT '24h',
		email_security  INTEGER NOT NULL DEFAULT 1,
		email_server    INTEGER NOT NULL DEFAULT 1,
		email_backups   INTEGER NOT NULL DEFAULT 1,
		email_users     INTEGER NOT NULL DEFAULT 1,
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE UNIQUE INDEX IF NOT EXISTS idx_admin_preferences_admin ON admin_preferences(admin_id)`,

	// Regular users (if project has user accounts)
	`CREATE TABLE IF NOT EXISTS users (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		username    TEXT NOT NULL UNIQUE,
		email       TEXT UNIQUE,
		password    TEXT NOT NULL,
		display_name TEXT,
		bio         TEXT,
		location    TEXT,
		website     TEXT,
		avatar_type TEXT NOT NULL DEFAULT 'gravatar',
		avatar_url  TEXT,
		visibility  TEXT NOT NULL DEFAULT 'public',
		org_visibility INTEGER NOT NULL DEFAULT 1,
		timezone    TEXT,
		language    TEXT NOT NULL DEFAULT 'en',
		role        TEXT NOT NULL DEFAULT 'user',
		enabled     INTEGER NOT NULL DEFAULT 1,
		verified    INTEGER NOT NULL DEFAULT 0,
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		last_login  INTEGER,
		failed_attempts INTEGER NOT NULL DEFAULT 0,
		locked_until INTEGER,
		metadata    TEXT
	)`,

	`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
	`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)`,

	// User preferences
	`CREATE TABLE IF NOT EXISTS user_preferences (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id     INTEGER NOT NULL UNIQUE,
		show_email      INTEGER NOT NULL DEFAULT 0,
		show_activity   INTEGER NOT NULL DEFAULT 1,
		show_orgs       INTEGER NOT NULL DEFAULT 1,
		searchable      INTEGER NOT NULL DEFAULT 1,
		email_security  INTEGER NOT NULL DEFAULT 1,
		email_mentions  INTEGER NOT NULL DEFAULT 1,
		email_updates   INTEGER NOT NULL DEFAULT 1,
		email_digest    TEXT NOT NULL DEFAULT 'weekly',
		push_enabled    INTEGER NOT NULL DEFAULT 0,
		push_mentions   INTEGER NOT NULL DEFAULT 1,
		theme           TEXT NOT NULL DEFAULT 'dark',
		font_size       TEXT NOT NULL DEFAULT 'medium',
		reduce_motion   INTEGER NOT NULL DEFAULT 0,
		date_format     TEXT NOT NULL DEFAULT 'YYYY-MM-DD',
		time_format     TEXT NOT NULL DEFAULT '24h',
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE UNIQUE INDEX IF NOT EXISTS idx_user_preferences_user ON user_preferences(user_id)`,

	// Organizations
	`CREATE TABLE IF NOT EXISTS orgs (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		slug        TEXT NOT NULL UNIQUE,
		name        TEXT NOT NULL,
		description TEXT,
		avatar_type TEXT NOT NULL DEFAULT 'gravatar',
		avatar_url  TEXT,
		visibility  TEXT NOT NULL DEFAULT 'public',
		owner_id    INTEGER NOT NULL,
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		metadata    TEXT
	)`,

	`CREATE UNIQUE INDEX IF NOT EXISTS idx_orgs_slug ON orgs(slug)`,
	`CREATE INDEX IF NOT EXISTS idx_orgs_owner ON orgs(owner_id)`,

	// Organization members
	`CREATE TABLE IF NOT EXISTS org_members (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		org_id      INTEGER NOT NULL,
		user_id     INTEGER NOT NULL,
		role        TEXT NOT NULL DEFAULT 'member',
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		UNIQUE(org_id, user_id)
	)`,

	`CREATE INDEX IF NOT EXISTS idx_org_members_org ON org_members(org_id)`,
	`CREATE INDEX IF NOT EXISTS idx_org_members_user ON org_members(user_id)`,

	// Organization preferences
	`CREATE TABLE IF NOT EXISTS org_preferences (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		org_id      INTEGER NOT NULL UNIQUE,
		default_member_role TEXT NOT NULL DEFAULT 'member',
		require_2fa     INTEGER NOT NULL DEFAULT 0,
		allow_invites   INTEGER NOT NULL DEFAULT 1,
		show_members    INTEGER NOT NULL DEFAULT 1,
		show_activity   INTEGER NOT NULL DEFAULT 1,
		notify_new_member   INTEGER NOT NULL DEFAULT 1,
		notify_member_leave INTEGER NOT NULL DEFAULT 1,
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE UNIQUE INDEX IF NOT EXISTS idx_org_preferences_org ON org_preferences(org_id)`,

	// Password reset tokens
	`CREATE TABLE IF NOT EXISTS password_resets (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		token_hash  TEXT NOT NULL UNIQUE,
		user_type   TEXT NOT NULL,
		user_id     INTEGER NOT NULL,
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		expires_at  INTEGER NOT NULL,
		used_at     INTEGER
	)`,

	`CREATE INDEX IF NOT EXISTS idx_password_resets_expires ON password_resets(expires_at)`,

	// Email verification tokens
	`CREATE TABLE IF NOT EXISTS email_verifications (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		token_hash  TEXT NOT NULL UNIQUE,
		user_type   TEXT NOT NULL,
		user_id     INTEGER NOT NULL,
		email       TEXT NOT NULL,
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		expires_at  INTEGER NOT NULL,
		verified_at INTEGER
	)`,

	`CREATE INDEX IF NOT EXISTS idx_email_verifications_expires ON email_verifications(expires_at)`,

	// TOTP secrets for 2FA
	`CREATE TABLE IF NOT EXISTS totp_secrets (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		user_type   TEXT NOT NULL,
		user_id     INTEGER NOT NULL UNIQUE,
		secret      TEXT NOT NULL,
		enabled     INTEGER NOT NULL DEFAULT 0,
		backup_codes TEXT,
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		last_used   INTEGER
	)`,

	`CREATE UNIQUE INDEX IF NOT EXISTS idx_totp_user ON totp_secrets(user_type, user_id)`,

	// User sessions (for regular users)
	`CREATE TABLE IF NOT EXISTS user_sessions (
		id          TEXT PRIMARY KEY,
		user_id     INTEGER NOT NULL,
		ip_address  TEXT NOT NULL,
		user_agent  TEXT,
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		expires_at  INTEGER NOT NULL,
		last_active INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_user_sessions_user ON user_sessions(user_id)`,
	`CREATE INDEX IF NOT EXISTS idx_user_sessions_expires ON user_sessions(expires_at)`,

	// Passkeys (WebAuthn/FIDO2)
	`CREATE TABLE IF NOT EXISTS passkeys (
		id              TEXT PRIMARY KEY,
		user_type       TEXT NOT NULL,
		user_id         INTEGER NOT NULL,
		name            TEXT NOT NULL,
		public_key      TEXT NOT NULL,
		sign_count      INTEGER NOT NULL DEFAULT 0,
		transports      TEXT,
		aaguid          TEXT,
		created_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		last_used       INTEGER
	)`,

	`CREATE INDEX IF NOT EXISTS idx_passkeys_user ON passkeys(user_type, user_id)`,

	// Trusted devices (skip 2FA)
	`CREATE TABLE IF NOT EXISTS trusted_devices (
		id          TEXT PRIMARY KEY,
		user_type   TEXT NOT NULL,
		user_id     INTEGER NOT NULL,
		device_hash TEXT NOT NULL,
		name        TEXT,
		created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		expires_at  INTEGER NOT NULL,
		last_used   INTEGER
	)`,

	`CREATE INDEX IF NOT EXISTS idx_trusted_devices_user ON trusted_devices(user_type, user_id)`,
}

// Combine all users.db table creation statements
func init() {
	// Append mail-specific tables
	usersCreateStatements = append(usersCreateStatements, mailCreateStatements...)
	// Append CalDAV tables
	usersCreateStatements = append(usersCreateStatements, caldavCreateStatements...)
	// Append CardDAV tables
	usersCreateStatements = append(usersCreateStatements, carddavCreateStatements...)
	// Append mailing list tables
	usersCreateStatements = append(usersCreateStatements, mailingListCreateStatements...)
	
	// Combine schema updates
	usersSchemaUpdates = append(usersSchemaUpdates, mailSchemaUpdates...)
}

// usersSchemaUpdates contains idempotent schema updates for users.db
// Per AI.md PART 10: Add new columns/indexes here, safe to run multiple times
var usersSchemaUpdates = []string{
	// Future schema updates go here
	// Example (commented):
	// `ALTER TABLE users ADD COLUMN verified_at INTEGER`,
}
