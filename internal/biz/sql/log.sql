create table `log`
(
    id        bigint auto_increment
        primary key,
    create_at datetime     not null,
    detail    varchar(512) null
) default charset utf8mb4;

