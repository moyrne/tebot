create table q_message
(
    id           serial  not null
        constraint q_message_pk
            primary key,
    time         integer,
    self_id      bigint,
    post_type    varchar not null,
    message_type varchar not null,
    sub_type     varchar not null,
    temp_source  varchar,
    message_id   bigint  not null,
    group_id     bigint,
    user_id      bigint,
    message      varchar,
    raw_message  varchar,
    font         integer,
    reply        varchar
);

alter table q_message
    owner to postgres;

create
index q_message_user_id_index
    on q_message (user_id);

