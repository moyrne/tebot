create table group
(
    id      bigint auto_increment
        primary key,
    user_id bigint not null,
    name    varchar(20) null,
    constraint group_user_id_uindex
        unique (user_id)
) default charset utf8mb4;

