create table tebot.q_group
(
    id   bigint auto_increment
        primary key,
    qgid bigint      not null,
    name varchar(20) null,
    constraint q_group_qgid_uindex
        unique (qgid)
) default charset utf8mb4;

