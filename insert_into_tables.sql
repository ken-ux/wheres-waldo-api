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
    -- Insert goals for easy mode.
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'easy'
        ),
        'Man with Snorkel',
        1175,
        640
    ),
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'easy'
        ),
        'Waldo',
        835,
        900
    ),
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'easy'
        ),
        'Man with Brown Jacket',
        550,
        585
    ),
    -- Insert goals for medium mode.
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'medium'
        ),
        'Waldo',
        1095,
        600
    ),
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'medium'
        ),
        'Knocked-Out Skier',
        415,
        270
    ),
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'medium'
        ),
        'Resting Skier',
        555,
        520
    ),
    -- Insert goals for hard mode.
    (
        (
            SELECT
                game_id
            FROM
                game
            WHERE
                game_name = 'hard'
        ),
        'Cowboy on Horse',
        957,
        236
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
        'Waldo',
        789,
        295
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
        'Man with Camera',
        248,
        416
    );

-- Insert data into user table as a test for leaderboard submissions.
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