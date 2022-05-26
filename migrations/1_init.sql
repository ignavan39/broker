create extension if not exists "uuid-ossp";

create table users (
    email text not null,
    password text not null,
    id uuid not null default uuid_generate_v4() constraint user_pk primary key
);

create unique index user_email_idx on users(email);

create table events (
    id uuid not null default uuid_generate_v4() constraint event_pk primary key,
    summary text not null,
    description text,
    start_date timestamp not null,
    end_date timestamp,
    time_zone text default 'Europe/Moscow'
);

create type remind_type as enum ('tg', 'vk', 'email', 'google_calendar_popup');

create table reminds (
    id uuid not null default uuid_generate_v4() constraint reminds_pk primary key,
    type remind_type not null,
    event_id uuid constraint reminds_events_fk references events (id)
);

create unique index reminds_events_idx on reminds(type, event_id);