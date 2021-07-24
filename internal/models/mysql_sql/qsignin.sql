create table q_sign_in
(
    id bigint auto_increment,
    quid bigint not null,
    create_at datetime null,
    day date null,
    constraint q_sign_in_pk
        primary key (id)
);

create index q_sign_in_quid_create_at_index
	on q_sign_in (quid, create_at);

create unique index q_sign_in_quid_day_uindex
	on q_sign_in (quid, day);

