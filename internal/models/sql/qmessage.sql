create table q_message
(
    id serial not null
        constraint q_message_pk
            primary key,
    time int,
    self_id bigint,
    post_type varchar not null,
    message_type varchar not null,
    sub_type varchar not null,
    temp_source varchar,
    message_id bigint not null,
    user_id bigint,
    message varchar,
    raw_message varchar,
    font int,
    reply varchar
);

create unique index q_message_message_id_uindex
	on q_message (message_id);

