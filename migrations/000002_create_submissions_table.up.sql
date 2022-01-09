CREATE TABLE IF NOT EXISTS submissions (
    id               INTEGER NOT NULL,
    subreddit_id     INTEGER,
    post_id          VARCHAR,
    author           VARCHAR,
    created_utc      DATETIME,
    domain           VARCHAR,
    over_18          BOOLEAN,
    permalink        VARCHAR,
    score            INTEGER,
    title            VARCHAR,
    url              VARCHAR,
    image_filename   VARCHAR,
    image_height_px  INTEGER,
    image_width_px   INTEGER,

    PRIMARY KEY (id),
    FOREIGN KEY(subreddit_id) REFERENCES subreddits (id),
    CHECK (over_18 IN (0, 1))
);
