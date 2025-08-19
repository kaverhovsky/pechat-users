-- +goose Up
-- +goose StatementBegin

create table users (
    id uuid primary key,
    nickname text,
    password_hash text,
    firstname text,
    lastname text,
    email text,
    bio text
)

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

drop table users;

-- +goose StatementEnd
