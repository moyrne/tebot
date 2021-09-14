create table `group`
(
    id      bigint auto_increment
        primary key,
    group_id bigint not null,
    name    varchar(20) null,
    constraint group_user_id_uindex
        unique (group_id)
) default charset utf8mb4;

