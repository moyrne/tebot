create table reply
(
    id       bigint auto_increment,
    user_id  bigint       null,
    msg      varchar(128) null,
    matches  varchar(512) not null,
    function varchar(128),
    replies  text         null,
    constraint reply_pk
        primary key (id)
) default charset utf8mb4;

create index reply_user_id_msg_index
    on reply (user_id, msg);

