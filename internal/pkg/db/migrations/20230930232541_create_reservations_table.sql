-- +goose Up
-- +goose StatementBegin
create table reservations(
      id BIGSERIAL PRIMARY KEY NOT NULL,
      start_date date NOT NULL,
      end_date date NOT NULL,
      room_id integer NOT NULL,
      created_at timestamp without time zone default now() NOT NULL,
      updated_at timestamp without time zone default now() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table reservations;
-- +goose StatementEnd
