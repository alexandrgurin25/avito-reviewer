CREATE TABLE teams (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);


CREATE TABLE users (
    id         TEXT PRIMARY KEY,
    username   TEXT NOT NULL,
    is_active  BOOLEAN NOT NULL,
    team_id    BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE
);


CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE pull_requests (
    id             TEXT PRIMARY KEY,
    name           TEXT NOT NULL,
    author_id      TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status         pr_status NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at      TIMESTAMPTZ
);


CREATE TABLE pull_request_reviewers (
    pr_id        TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id  TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (pr_id, reviewer_id)
);
