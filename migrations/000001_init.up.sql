CREATE TABLE IF NOT EXISTS users(
  id         UUID        PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users_sessions(
  user_id       UUID          REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  refresh_token TEXT          UNIQUE                NOT NULL,
  user_agent    TEXT          NOT NULL,
  ip            INET          NOT NULL,
  created_at    TIMESTAMPTZ   DEFAULT now()         NOT NULL,
  expires_at    TIMESTAMPTZ   NOT NULL  
);