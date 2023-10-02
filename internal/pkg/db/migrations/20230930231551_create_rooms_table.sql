-- +goose Up
-- +goose StatementBegin
create table rooms(
     id SERIAL PRIMARY KEY NOT NULL,
     name character varying(255) NOT NULL DEFAULT '' UNIQUE,
     cost double precision NOT NULL DEFAULT 0,
     created_at timestamp without time zone default now() NOT NULL,
     updated_at timestamp without time zone default now() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table rooms;
-- +goose StatementEnd
