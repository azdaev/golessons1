-- +goose Up
-- +goose StatementBegin
create table posts (
    id bigserial primary key,
    user_id bigint not null,
    title text not null,
    body text not null,
    views integer not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table users (
    id bigserial primary key,
    name text not null,
    email text not null,
    is_admin boolean not null default false,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table comments (
    id bigserial primary key,
    user_id bigint not null,
    post_id bigint not null,
    body text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
