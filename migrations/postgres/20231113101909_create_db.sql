-- +goose Up
-- +goose StatementBegin
create table users
(
    id         uuid primary key,
    name       varchar,
    username   varchar not null,
    password   varchar not null,
    role       varchar not null,
    created_at timestamp default now(),
    updated_at timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd

