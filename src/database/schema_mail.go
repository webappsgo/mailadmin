package database

// schema_mail.go contains mail-specific database tables for users.db
// These tables support the mail infrastructure management features

// Mail infrastructure tables for users.db
var mailCreateStatements = []string{
	// Mail servers - per PLAN.AI.md line 1121
	`CREATE TABLE IF NOT EXISTS srv_mail_servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hostname TEXT NOT NULL UNIQUE,
		ip_address TEXT,
		server_role TEXT NOT NULL,
		agent_token TEXT NOT NULL,
		agent_port INTEGER DEFAULT 64100,
		primary_for_domains TEXT,
		relay_domains TEXT,
		relay_to_host TEXT,
		smarthost TEXT,
		postfix_running INTEGER DEFAULT 0,
		dovecot_running INTEGER DEFAULT 0,
		amavisd_running INTEGER DEFAULT 0,
		clamav_running INTEGER DEFAULT 0,
		spam_filter_running INTEGER DEFAULT 0,
		last_seen INTEGER,
		status TEXT DEFAULT 'unknown',
		load_average TEXT,
		disk_usage_percent INTEGER,
		queue_size INTEGER,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		enabled INTEGER DEFAULT 1
	)`,

	// Mail domains - virtual domains hosted
	`CREATE TABLE IF NOT EXISTS mail_domains (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT NOT NULL UNIQUE,
		description TEXT,
		owner_type TEXT NOT NULL,
		owner_id INTEGER NOT NULL,
		max_mailboxes INTEGER DEFAULT 0,
		max_aliases INTEGER DEFAULT 0,
		max_quota_mb INTEGER DEFAULT 0,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_domains_owner ON mail_domains(owner_type, owner_id)`,

	// Mail mailboxes - virtual mailboxes
	`CREATE TABLE IF NOT EXISTS mail_mailboxes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain_id INTEGER NOT NULL,
		local_part TEXT NOT NULL,
		password TEXT NOT NULL,
		name TEXT,
		quota_mb INTEGER DEFAULT 0,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		UNIQUE(domain_id, local_part)
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_mailboxes_domain ON mail_mailboxes(domain_id)`,

	// Mail aliases - email aliases
	`CREATE TABLE IF NOT EXISTS mail_aliases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain_id INTEGER NOT NULL,
		source TEXT NOT NULL,
		destination TEXT NOT NULL,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_aliases_domain ON mail_aliases(domain_id)`,
	`CREATE INDEX IF NOT EXISTS idx_mail_aliases_source ON mail_aliases(source)`,

	// Mail forwards - forwarding rules
	`CREATE TABLE IF NOT EXISTS mail_forwards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mailbox_id INTEGER NOT NULL,
		forward_to TEXT NOT NULL,
		keep_local INTEGER DEFAULT 1,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_forwards_mailbox ON mail_forwards(mailbox_id)`,

	// DKIM keys - per-domain DKIM signing keys
	`CREATE TABLE IF NOT EXISTS mail_dkim_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain_id INTEGER NOT NULL,
		selector TEXT NOT NULL,
		private_key TEXT NOT NULL,
		public_key TEXT NOT NULL,
		key_size INTEGER DEFAULT 2048,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		expires_at INTEGER,
		UNIQUE(domain_id, selector)
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_dkim_domain ON mail_dkim_keys(domain_id)`,

	// SPF records
	`CREATE TABLE IF NOT EXISTS mail_spf_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain_id INTEGER NOT NULL UNIQUE,
		policy TEXT NOT NULL,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	// DMARC policies
	`CREATE TABLE IF NOT EXISTS mail_dmarc_policies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain_id INTEGER NOT NULL UNIQUE,
		policy TEXT NOT NULL DEFAULT 'none',
		subdomain_policy TEXT,
		percentage INTEGER DEFAULT 100,
		rua_email TEXT,
		ruf_email TEXT,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	// PGP/GPG keys
	`CREATE TABLE IF NOT EXISTS mail_pgp_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mailbox_id INTEGER NOT NULL,
		fingerprint TEXT NOT NULL UNIQUE,
		public_key TEXT NOT NULL,
		private_key TEXT,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		expires_at INTEGER
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_pgp_mailbox ON mail_pgp_keys(mailbox_id)`,

	// Mail queue log
	`CREATE TABLE IF NOT EXISTS mail_queue_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		queue_id TEXT NOT NULL,
		mail_server_id INTEGER,
		sender TEXT NOT NULL,
		recipient TEXT NOT NULL,
		size_bytes INTEGER,
		status TEXT NOT NULL,
		timestamp INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_queue_timestamp ON mail_queue_log(timestamp)`,
	`CREATE INDEX IF NOT EXISTS idx_mail_queue_server ON mail_queue_log(mail_server_id)`,

	// Mail delivery log
	`CREATE TABLE IF NOT EXISTS mail_delivery_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		queue_id TEXT,
		mail_server_id INTEGER,
		sender TEXT NOT NULL,
		recipient TEXT NOT NULL,
		size_bytes INTEGER,
		status TEXT NOT NULL,
		delay_seconds INTEGER,
		dsn TEXT,
		timestamp INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_delivery_timestamp ON mail_delivery_log(timestamp)`,
	`CREATE INDEX IF NOT EXISTS idx_mail_delivery_recipient ON mail_delivery_log(recipient)`,
	`CREATE INDEX IF NOT EXISTS idx_mail_delivery_server ON mail_delivery_log(mail_server_id)`,

	// Mail stats hourly
	`CREATE TABLE IF NOT EXISTS mail_stats_hourly (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mail_server_id INTEGER NOT NULL,
		hour_timestamp INTEGER NOT NULL,
		sent_count INTEGER DEFAULT 0,
		received_count INTEGER DEFAULT 0,
		bounced_count INTEGER DEFAULT 0,
		deferred_count INTEGER DEFAULT 0,
		spam_count INTEGER DEFAULT 0,
		virus_count INTEGER DEFAULT 0,
		total_size_bytes INTEGER DEFAULT 0,
		UNIQUE(mail_server_id, hour_timestamp)
	)`,

	`CREATE INDEX IF NOT EXISTS idx_stats_hourly_server ON mail_stats_hourly(mail_server_id)`,
	`CREATE INDEX IF NOT EXISTS idx_stats_hourly_time ON mail_stats_hourly(hour_timestamp)`,

	// Mail stats daily
	`CREATE TABLE IF NOT EXISTS mail_stats_daily (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mail_server_id INTEGER NOT NULL,
		day_timestamp INTEGER NOT NULL,
		sent_count INTEGER DEFAULT 0,
		received_count INTEGER DEFAULT 0,
		bounced_count INTEGER DEFAULT 0,
		deferred_count INTEGER DEFAULT 0,
		spam_count INTEGER DEFAULT 0,
		virus_count INTEGER DEFAULT 0,
		total_size_bytes INTEGER DEFAULT 0,
		UNIQUE(mail_server_id, day_timestamp)
	)`,

	`CREATE INDEX IF NOT EXISTS idx_stats_daily_server ON mail_stats_daily(mail_server_id)`,
	`CREATE INDEX IF NOT EXISTS idx_stats_daily_time ON mail_stats_daily(day_timestamp)`,

	// Sieve scripts - mail filtering
	`CREATE TABLE IF NOT EXISTS mail_sieve_scripts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mailbox_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		script_content TEXT NOT NULL,
		priority INTEGER DEFAULT 0,
		active INTEGER DEFAULT 0,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_sieve_mailbox ON mail_sieve_scripts(mailbox_id)`,

	// Archive policies
	`CREATE TABLE IF NOT EXISTS mail_archive_policies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain_id INTEGER,
		mailbox_id INTEGER,
		archive_after_days INTEGER DEFAULT 365,
		delete_after_days INTEGER DEFAULT 0,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_archive_domain ON mail_archive_policies(domain_id)`,
	`CREATE INDEX IF NOT EXISTS idx_mail_archive_mailbox ON mail_archive_policies(mailbox_id)`,

	// Archived messages
	`CREATE TABLE IF NOT EXISTS mail_archived_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mailbox_id INTEGER NOT NULL,
		message_id TEXT NOT NULL,
		subject TEXT,
		sender TEXT,
		recipient TEXT,
		size_bytes INTEGER,
		archived_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		delete_at INTEGER,
		storage_path TEXT
	)`,

	`CREATE INDEX IF NOT EXISTS idx_mail_archived_mailbox ON mail_archived_messages(mailbox_id)`,
	`CREATE INDEX IF NOT EXISTS idx_mail_archived_delete ON mail_archived_messages(delete_at)`,
}

// CalDAV tables for users.db
var caldavCreateStatements = []string{
	// Calendars
	`CREATE TABLE IF NOT EXISTS cal_calendars (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner_type TEXT NOT NULL,
		owner_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		color TEXT DEFAULT '#3b82f6',
		timezone TEXT DEFAULT 'UTC',
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_cal_calendars_owner ON cal_calendars(owner_type, owner_id)`,

	// Calendar events
	`CREATE TABLE IF NOT EXISTS cal_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		calendar_id INTEGER NOT NULL,
		uid TEXT NOT NULL UNIQUE,
		summary TEXT NOT NULL,
		description TEXT,
		location TEXT,
		start_time INTEGER NOT NULL,
		end_time INTEGER NOT NULL,
		all_day INTEGER DEFAULT 0,
		recurring INTEGER DEFAULT 0,
		recurrence_rule TEXT,
		ical_data TEXT NOT NULL,
		etag TEXT NOT NULL,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_cal_events_calendar ON cal_events(calendar_id)`,
	`CREATE INDEX IF NOT EXISTS idx_cal_events_start ON cal_events(start_time)`,

	// Event attendees
	`CREATE TABLE IF NOT EXISTS cal_attendees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER NOT NULL,
		email TEXT NOT NULL,
		name TEXT,
		role TEXT DEFAULT 'REQ-PARTICIPANT',
		status TEXT DEFAULT 'NEEDS-ACTION',
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_cal_attendees_event ON cal_attendees(event_id)`,

	// Calendar shares
	`CREATE TABLE IF NOT EXISTS cal_shares (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		calendar_id INTEGER NOT NULL,
		shared_with_type TEXT NOT NULL,
		shared_with_id INTEGER NOT NULL,
		permission TEXT NOT NULL DEFAULT 'read',
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		UNIQUE(calendar_id, shared_with_type, shared_with_id)
	)`,

	`CREATE INDEX IF NOT EXISTS idx_cal_shares_calendar ON cal_shares(calendar_id)`,
}

// CardDAV tables for users.db
var carddavCreateStatements = []string{
	// Address books
	`CREATE TABLE IF NOT EXISTS card_addressbooks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner_type TEXT NOT NULL,
		owner_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_card_addressbooks_owner ON card_addressbooks(owner_type, owner_id)`,

	// Contacts
	`CREATE TABLE IF NOT EXISTS card_contacts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		addressbook_id INTEGER NOT NULL,
		uid TEXT NOT NULL UNIQUE,
		full_name TEXT NOT NULL,
		given_name TEXT,
		family_name TEXT,
		email TEXT,
		phone TEXT,
		organization TEXT,
		vcard_data TEXT NOT NULL,
		etag TEXT NOT NULL,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_card_contacts_addressbook ON card_contacts(addressbook_id)`,
	`CREATE INDEX IF NOT EXISTS idx_card_contacts_name ON card_contacts(full_name)`,
	`CREATE INDEX IF NOT EXISTS idx_card_contacts_email ON card_contacts(email)`,

	// Address book shares
	`CREATE TABLE IF NOT EXISTS card_shares (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		addressbook_id INTEGER NOT NULL,
		shared_with_type TEXT NOT NULL,
		shared_with_id INTEGER NOT NULL,
		permission TEXT NOT NULL DEFAULT 'read',
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		UNIQUE(addressbook_id, shared_with_type, shared_with_id)
	)`,

	`CREATE INDEX IF NOT EXISTS idx_card_shares_addressbook ON card_shares(addressbook_id)`,
}

// Mailing list tables for users.db
var mailingListCreateStatements = []string{
	// Mailing lists
	`CREATE TABLE IF NOT EXISTS list_lists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner_type TEXT NOT NULL,
		owner_id INTEGER NOT NULL,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		description TEXT,
		list_type TEXT NOT NULL DEFAULT 'discussion',
		moderated INTEGER DEFAULT 0,
		public_archive INTEGER DEFAULT 1,
		subscribe_policy TEXT DEFAULT 'open',
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_list_lists_owner ON list_lists(owner_type, owner_id)`,

	// List members
	`CREATE TABLE IF NOT EXISTS list_members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		list_id INTEGER NOT NULL,
		email TEXT NOT NULL,
		name TEXT,
		delivery_mode TEXT DEFAULT 'enabled',
		moderator INTEGER DEFAULT 0,
		subscribed_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		UNIQUE(list_id, email)
	)`,

	`CREATE INDEX IF NOT EXISTS idx_list_members_list ON list_members(list_id)`,

	// List archives
	`CREATE TABLE IF NOT EXISTS list_archives (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		list_id INTEGER NOT NULL,
		message_id TEXT NOT NULL,
		sender TEXT NOT NULL,
		subject TEXT NOT NULL,
		body TEXT NOT NULL,
		posted_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`,

	`CREATE INDEX IF NOT EXISTS idx_list_archives_list ON list_archives(list_id)`,
	`CREATE INDEX IF NOT EXISTS idx_list_archives_posted ON list_archives(posted_at)`,
}

// mailSchemaUpdates contains idempotent schema updates for mail tables
var mailSchemaUpdates = []string{
	// Future mail schema updates go here
}
