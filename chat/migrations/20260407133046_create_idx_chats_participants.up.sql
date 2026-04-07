CREATE INDEX idx_chats_participants
ON chats
USING GIN (participant_ids);