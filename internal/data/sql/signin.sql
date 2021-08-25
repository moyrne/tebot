create table sign_in
(
    id bigint auto_increment,
    user_id bigint not null,
    create_at datetime null,
    day date null,
    constraint sign_in_pk
        primary key (id)
) default charset utf8mb4;

create index sign_in_user_id_create_at_index
	on sign_in (user_id, create_at);

create unique index sign_in_user_id_day_uindex
	on sign_in (user_id, day);

