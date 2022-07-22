CREATE TABLE sessions
(
    session_id   UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    expired_at     TIMESTAMP with time zone,
    id       UUID ,
    FOREIGN KEY (id)
        references user_details(id)
);