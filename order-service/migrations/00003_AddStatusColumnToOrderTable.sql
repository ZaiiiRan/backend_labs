-- +goose Up

alter table orders add column status text not null default 'created';
alter type v1_order add attribute status text;

-- +goose Down
alter table orders drop column status;
alter type v1_order drop attribute status;
