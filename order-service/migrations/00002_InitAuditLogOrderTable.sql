-- +goose Up
create table if not exists audit_log_order (
    id bigserial not null primary key,
    order_id bigint not null,
    order_item_id bigint not null,
    customer_id bigint not null,
    order_status text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);

create type v1_audit_log_order as (
    id bigint,
    order_id bigint,
    order_item_id bigint,
    customer_id bigint,
    order_status text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

-- +goose Down
drop table if exists audit_log_order;
drop type if exists v1_audit_log_order;