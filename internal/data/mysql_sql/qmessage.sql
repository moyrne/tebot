create table q_message
(
    id bigint auto_increment,
    time integer null,
    self_id bigint null,
    post_type varchar(128) not null,
    message_type varchar(128) not null,
    sub_type varchar(128) null,
    temp_source varchar(128) null,
    message_id bigint not null,
    group_id bigint null,
    user_id bigint null,
    message text null,
    raw_message text null,
    font int null,
    reply varchar(256) null,
    constraint q_message_pk
        primary key (id)
) default charset utf8mb4;

create index q_message_user_id_index
	on q_message (user_id);

