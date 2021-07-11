create table q_reply
(
    id     bigserial not null
        constraint q_reply_pk
            primary key,
    quid   bigint,
    msg    varchar,
    reply  json
);

alter table q_reply
    owner to postgres;

create index q_reply_quid_msg_index
    on q_reply (quid, msg);

