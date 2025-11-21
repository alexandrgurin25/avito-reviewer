
CREATE TABLE teams (
    name TEXT PRIMARY KEY
);


CREATE TABLE users (
    id         TEXT PRIMARY KEY,           
    username   TEXT NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    team_name  TEXT NOT NULL REFERENCES teams(name) ON DELETE CASCADE
);



CREATE TABLE pull_requests (
    id             TEXT PRIMARY KEY,        
    name           TEXT NOT NULL,            
    author_id      TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status         TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at      TIMESTAMPTZ
);


CREATE TABLE pull_request_reviewers (
    pr_id        TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id  TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (pr_id, reviewer_id)
);

