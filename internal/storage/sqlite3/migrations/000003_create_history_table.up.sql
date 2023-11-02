CREATE TABLE IF NOT EXISTS history (
    id            INTEGER NOT NULL,
    submission_id INTEGER,
    date          DATETIME DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    FOREIGN KEY(submission_id) REFERENCES submissions (id)
);
