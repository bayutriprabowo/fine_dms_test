-- table m_user
CREATE TABLE m_user(
	id		BIGSERIAL PRIMARY KEY,
	username	VARCHAR(64) UNIQUE NOT NULL,
	password	VARCHAR(255) NOT NULL,
	email		VARCHAR(128) UNIQUE NOT NULL,
	first_name	VARCHAR(64) NOT NULL,
	last_name	VARCHAR(64),
	created_at	TIMESTAMP NOT NULL,
	updated_at	TIMESTAMP NOT NULL
);

-- table m_file
CREATE TABLE m_file(
	id		BIGSERIAL PRIMARY KEY,
	path		VARCHAR(255) NOT NULL,
	ext		VARCHAR(64) NOT NULL,
	user_id		BIGINT REFERENCES m_user(id) NOT NULL,
	created_at	TIMESTAMP NOT NULL,
	updated_at	TIMESTAMP NOT NULL
);

-- table m_tags
CREATE TABLE m_tags(
	id		BIGSERIAL PRIMARY KEY,
	name		VARCHAR(16) NOT NULL,
	created_at	TIMESTAMP NOT NULL,
	updated_at	TIMESTAMP NOT NULL
);

-- table t_tags
CREATE TABLE t_tags(
	id		BIGSERIAL PRIMARY KEY,
	tags_id		BIGINT REFERENCES m_tags(id) NOT NULL,
	file_id		BIGINT REFERENCES m_file(id) NOT NULL
);

-- table t_file_upload
CREATE TABLE t_file_upload(
	id		BIGSERIAL PRIMARY KEY,
	file_id		BIGINT REFERENCES m_file(id) NOT NULL,
	user_id		BIGINT REFERENCES m_user(id) NOT NULL,
	created_at	TIMESTAMP NOT NULL
);

-- table t_file_download
CREATE TABLE t_file_download(
	id		BIGSERIAL PRIMARY KEY,
	file_id		BIGINT REFERENCES m_file(id) NOT NULL,
	user_id		BIGINT REFERENCES m_user(id) NOT NULL,
	created_at	TIMESTAMP NOT NULL
);
