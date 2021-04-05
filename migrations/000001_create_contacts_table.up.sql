CREATE TABLE "contacts" (
    id bigserial primary key,
    name varchar(255) NOT NULL,
    phone varchar(255) NOT NULL,
    address varchar(255) NOT NULL,
    favorites JSONB NOT NULL DEFAULT '{}',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create unique index contact_name on "contacts" (name);