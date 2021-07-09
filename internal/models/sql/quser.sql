create table q_user
(
    id       bigserial             not null
        constraint q_user_pk
            primary key,
    quid     bigint                not null,
    nickname varchar,
    sex      varchar,
    age      integer,
    bind_area varchar,
    mode      integer,
    ban      boolean default false not null
);

alter table q_user
    owner to postgres;

create
unique index q_user_quid_uindex
    on q_user (quid);

