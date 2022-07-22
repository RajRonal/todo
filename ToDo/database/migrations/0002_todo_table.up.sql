ALTER TABLE todo  RENAME TO user_details;
CREATE TABLE todo
(
    task_id   UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    task     TEXT NOT NULL ,
    status   BOOLEAN NOT NULL DEFAULT FALSE,
    id       UUID ,
        FOREIGN KEY (id)
        references user_details(id)
);