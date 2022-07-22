create extension if not exists "uuid-ossp";

CREATE TABLE users (
    id uuid NOT NULL DEFAULT uuid_generate_v4() CONSTRAINT user_pk PRIMARY KEY,
    email TEXT NOT NULL,
    nickname TEXT NOT NULL,
    "password" TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    avatar_url TEXT NOT NULL DEFAULT 'https://vk.com/images/camera_c.gif'
);

CREATE UNIQUE INDEX user_email_idx ON users(email);
CREATE UNIQUE INDEX user_nickname_idx ON users(nickname);

CREATE TABLE workspaces (
    id uuid NOT NULL DEFAULT uuid_generate_v4() CONSTRAINT workspace_pk PRIMARY KEY,
    "name" TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_private BOOLEAN DEFAULT TRUE
);

CREATE TYPE workspace_access_type AS ENUM (
    'ADMIN',
    'MODERATOR',
    'USER'
);

CREATE TABLE workspace_accesses (
    id uuid NOT NULL DEFAULT uuid_generate_v4() CONSTRAINT workspace_access_pk PRIMARY KEY,
    workspace_id uuid NOT NULL CONSTRAINT wa_workspace_id_fk REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id uuid NOT NULL CONSTRAINT wa_user_id_fk REFERENCES users(id) ON DELETE CASCADE,
    "type" workspace_access_type NOT NULL DEFAULT 'USER'
);

CREATE UNIQUE INDEX user_workspace_access_idx ON workspace_accesses(user_id, workspace_id);

CREATE TABLE peers (
    id uuid NOT NULL DEFAULT uuid_generate_v4() CONSTRAINT peer_pk PRIMARY KEY,
    "name" TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    workspace_id uuid NOT NULL CONSTRAINT peer_workspace_id_fk REFERENCES workspaces(id) ON DELETE CASCADE
);

CREATE TABLE user_peers (
    id uuid NOT NULL DEFAULT uuid_generate_v4() CONSTRAINT user_peers_pk PRIMARY KEY,
    user_id uuid NOT NULL CONSTRAINT user_peers_user_id_fk REFERENCES users(id) ON DELETE CASCADE,
    peer_id uuid NOT NULL CONSTRAINT user_peers_peer_id_fk REFERENCES peers(id) ON DELETE CASCADE,
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX user_peers_id_idx ON user_peers(user_id, peer_id);

CREATE TABLE messages (
    id uuid NOT NULL DEFAULT uuid_generate_v4() CONSTRAINT messages_pk PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    user_id uuid NOT NULL CONSTRAINT messages_user_id_fk REFERENCES users(id) ON DELETE CASCADE,
    peer_id uuid NOT NULL CONSTRAINT messages_peer_id_fk REFERENCES peers(id) ON DELETE CASCADE,
    "text" TEXT NOT NULL,
    parent_id uuid DEFAULT NULL CONSTRAINT messages_parent_fk REFERENCES messages(id)
);

CREATE TYPE invitation_status AS ENUM (
    'PENDING',
    'ACCEPTED',
    'CANCELED',
    'EXPIRED'
);

CREATE TYPE invitation_system_status AS ENUM (
    'CREATED',
    'SEND',
    'DELIVERED',
    'REJECT'
);

CREATE TABLE invitations (
    id uuid NOT NULL DEFAULT uuid_generate_v4() CONSTRAINT invites_pk PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sender_id uuid NOT NULL CONSTRAINT invitations_sender_id_fk REFERENCES users(id),
    recipient_email TEXT NOT NULL,
    workspace_id uuid NOT NULL CONSTRAINT invitations_workspace_id_fk REFERENCES workspaces(id) ON DELETE CASCADE,
    "status" invitation_status NOT NULL DEFAULT 'PENDING',
    system_status invitation_system_status NOT NULL DEFAULT 'CREATED',
    code text NOT NULL
);

CREATE UNIQUE INDEX invitation_idx ON invitations(recipient_email, workspace_id, "status")