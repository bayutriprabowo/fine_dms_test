-- table m_user
create table m_user (
    id bigserial primary key,
    username varchar(64) not null unique,
    password varchar(255) not null,
    email varchar(128) not null unique,
    first_name varchar(64) not null,
    last_name varchar(64),
    created_at timestamp not null,
    updated_at timestamp not null
);


-- table m_file
create table m_file (
    id bigserial primary key,
    path varchar(255) not null,
    ext varchar(64) not null,
    user_id bigint not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint fk_file_user
		foreign key (user_id)
			references m_user (id)
);

-- table m_tags
create table m_tags (
    id bigserial primary key,
    name varchar(64) not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- table t_tags
create table t_tags (
    id bigserial primary key,
    tags_id bigint not null,
    file_id bigint not null,

    constraint fk_tags_tags
		foreign key (tags_id)
			references m_tags (id)
    constraint fk_tags_file
		foreign key (file_id)
			references m_file (id)
);

-- table t_upload
create table t_upload (
    id bigserial primary key,
    file_id bigint not null,
    user_id bigint not null,
    date timestamp not null,
    constraint fk_upload_file
		foreign key (file_id)
			references m_file (id),
	constraint fk_upload_user
		foreign key (user_id)
			references m_user (id)
);

-- table t_download
create table t_download (
    id bigserial primary key,
    file_id bigint not null,
    user_id bigint not null,
    date timestamp not null,
    constraint fk_download_file
		foreign key (file_id)
			references m_file (id),
	constraint fk_download_user
		foreign key (user_id)
			references m_user (id)
);
