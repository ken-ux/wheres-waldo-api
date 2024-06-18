-- Create game table for managing game modes.
CREATE TABLE
    game (
        game_id SMALLINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        game_name VARCHAR(10) NOT NULL
    );

-- Create goal table for storing coordinates.
CREATE TABLE
    goal (
        goal_id SMALLINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        game_id SMALLINT NOT NULL,
        goal_desc VARCHAR(50) NOT NULL,
        goal_pos_x INTEGER NOT NULL,
        goal_pos_y INTEGER NOT NULL,
        FOREIGN KEY (game_id) REFERENCES game (game_id)
    );

-- Create user table for leaderboard scores.
CREATE TABLE
    "user" ( -- This is in double quotes because user is a reserved word.
        user_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        game_id SMALLINT NOT NULL,
        user_name VARCHAR(20) NOT NULL,
        user_score INTEGER,
        FOREIGN KEY (game_id) REFERENCES game (game_id)
    );

-- Insert data into game table since it has no FKs.
INSERT INTO
    game (game_name)
VALUES
    ('easy'), -- Can't use double quotes otherwise throws a syntax error.
    ('medium'),
    ('hard');

-- Insert data into goal table.
INSERT INTO
    goal (game_id, goal_desc, goal_pos_x, goal_pos_y)
VALUES
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'medium'
        ),
        'This is the first goal.',
        125,
        711
    ),
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'hard'
        ),
        'This is the second goal.',
        500,
        212
    );

-- Insert data into user table as a test.
INSERT INTO
    "user" (game_id, user_name, user_score)
VALUES
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'easy'
        ),
        'Michael',
        120
    ),
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'hard'
        ),
        'Kelly',
        500
    );