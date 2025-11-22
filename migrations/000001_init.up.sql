CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    team_name VARCHAR(255) UNIQUE NOT NULL
);
CREATE INDEX teams_team_name_index ON teams (team_name);

CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);
CREATE INDEX users_is_active_index ON users (is_active);

CREATE TABLE team_members (
    team_id INTEGER NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    PRIMARY KEY (team_id, user_id)
);
CREATE INDEX team_members_user_id_index ON team_members (user_id);
CREATE INDEX team_members_team_id_index ON team_members (team_id);

CREATE TABLE pull_requests (
    pull_request_id VARCHAR(255) PRIMARY KEY,
    pull_request_name VARCHAR(500) NOT NULL,
    author_id VARCHAR(255) NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    status pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMP NULL,
    CONSTRAINT check_merged_at CHECK (
        (status = 'MERGED' AND merged_at IS NOT NULL) OR
        (status = 'OPEN' AND merged_at IS NULL)
    )
);
CREATE INDEX pull_requests_author_id_index ON pull_requests (author_id);
CREATE INDEX pull_requests_status_index ON pull_requests (status);

CREATE TABLE reviewers (
    pull_request_id VARCHAR(255) NOT NULL REFERENCES pull_requests (pull_request_id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (pull_request_id, user_id)
);
CREATE INDEX reviewers_user_id_index ON reviewers (user_id);
CREATE INDEX reviewers_pull_request_id_index ON reviewers (pull_request_id);
