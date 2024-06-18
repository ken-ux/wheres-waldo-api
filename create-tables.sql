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
    "user" ( -- This is in quotes because user is a reserved word.
        user_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        game_id SMALLINT NOT NULL,
        user_name VARCHAR(20) NOT NULL,
        user_score INTEGER,
        FOREIGN KEY (game_id) REFERENCES game (game_id)
    );