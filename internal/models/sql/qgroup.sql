create table q_group
(
    id   bigserial not null
        constraint q_group_pk
            primary key,
    qgid bigint    not null,
    name varchar
);

alter table q_group
    owner to postgres;

create
unique index q_group_qgid_uindex
    on q_group (qgid);

