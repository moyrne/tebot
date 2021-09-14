create table user
(
    id        bigint auto_increment,
    user_id   bigint             not null,
    nickname  varchar(128)       null,
    sex       varchar(16)        null,
    age       int                null,
    bind_area varchar(128)       null,
    mode      int                null,
    ban       bool default false not null,
    constraint user_pk
        primary key (id)
) default charset utf8mb4;

create unique index user_user_id_uindex
    on user (user_id);

