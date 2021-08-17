create table q_sign_in
(
    id        bigserial not null
        constraint q_sign_in_pk
            primary key,
    quid      bigint    not null,
    create_at date,
    day       date,
    constraint q_sign_in_pk_2
        unique (quid, day)
);

alter table q_sign_in
    owner to postgres;

create index q_sign_in_quid_create_at_index
    on q_sign_in (quid, create_at);

