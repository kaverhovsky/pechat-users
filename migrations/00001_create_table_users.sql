-- +goose Up
-- +goose StatementBegin

create table users (
    id uuid primary key,
    nickname text not null,
    password_hash text not null,
    firstname text,
    lastname text,
    email text not null,
    bio text
)

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

drop table users;

-- +goose StatementEnd
