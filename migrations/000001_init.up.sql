CREATE TABLE IF NOT EXISTS users_sessions(
  id                  UUID          PRIMARY KEY,
  user_id             UUID          NOT NULL,
  refresh_token_hash  TEXT          UNIQUE                NOT NULL,
  user_agent          TEXT          NOT NULL,
  ip                  INET          NOT NULL,
  created_at          TIMESTAMPTZ   DEFAULT now()         NOT NULL,
  expires_at          TIMESTAMPTZ   NOT NULL  
);