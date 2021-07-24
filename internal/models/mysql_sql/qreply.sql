create table q_reply
(
    id bigint auto_increment,
    quid bigint null,
    msg varchar(256) null,
    reply json null,
    constraint q_reply_pk
        primary key (id)
);

create index q_reply_quid_msg_index
	on q_reply (quid, msg);

