create table q_user
(
    id bigint auto_increment,
    quid bigint not null,
    nickname varchar(128) null,
    sex varchar(16) null,
    age int null,
    bind_area varchar(128) null,
    mode int null,
    ban bool default false not null,
    constraint q_user_pk
        primary key (id)
) default charset utf8mb4;

create unique index q_user_quid_uindex
	on q_user (quid);

