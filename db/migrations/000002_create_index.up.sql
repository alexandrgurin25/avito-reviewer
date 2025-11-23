CREATE INDEX idx_pr_reviewers_reviewer ON pull_request_reviewers(reviewer_id);

CREATE INDEX idx_pr_author ON pull_requests(author_id);
CREATE INDEX idx_pr_status ON pull_requests(status);

CREATE INDEX idx_users_team_id ON users(team_id);
CREATE INDEX idx_teams_name ON teams(name);
