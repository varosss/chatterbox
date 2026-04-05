CREATE TABLE chats (
    id UUID PRIMARY KEY,
    participant_ids UUID[] NOT NULL
);