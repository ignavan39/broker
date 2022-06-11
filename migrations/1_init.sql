create extension if not exists "uuid-ossp";

CREATE TABLE users (
    email text not null,
    password text not null,
    first_name text not null,
    last_name text not null,
    id uuid not null default uuid_generate_v4() constraint user_pk primary key
);

create unique index user_email_idx on users(email);

CREATE TABLE chats (
    id uuid not null primary key,
    created_at timestamp not null default NOW()
);

CREATE TABLE user_chats (
    id uuid not null default uuid_generate_v4() primary key,
    user_id uuid not null constraint user_chat_id references users(id),
    chat_id uuid not null constraint chat_user_id references chats(id),
    is_blocked boolean default false
);
CREATE TABLE messages (
    created_at timestamp not null default NOW(),
    id uuid not null default uuid_generate_v4() primary key,
    text text,
    forwards jsonb default '[]',
    images jsonb default '[]',
    sender_id uuid not null constraint user_id references users(id),
    chat_id uuid not null constraint chat_id references chats(id)
);