create table q_reply
(
    id bigint auto_increment,
    quid bigint null,
    msg varchar(256) null,
    reply text null,
    constraint q_reply_pk
        primary key (id)
) default charset utf8mb4;

create index q_reply_quid_msg_index
	on q_reply (quid, msg);

