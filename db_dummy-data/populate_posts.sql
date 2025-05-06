BEGIN TRANSACTION;

-- Create a temporary table to cache 5 users for use in post generation.
-- This improves performance by reducing repeated queries to the Users table.
CREATE TEMP TABLE TempUsers AS
SELECT ID AS UserID, Username, Avatar FROM Users LIMIT 5;

-- Prepare indexed users by assigning a row number to each user in the temporary table.
-- This allows efficient mapping of users to posts in subsequent queries.
WITH IndexedUsers AS (
    SELECT UserID, Username, Avatar, ROW_NUMBER() OVER () AS rn
    FROM TempUsers
),
-- Generate a sequence of numbers from 1 to 1000 using a recursive CTE.
-- This sequence is used to create multiple posts for testing purposes.
Counter(x) AS (
    SELECT 1
    UNION ALL
    SELECT x+1 FROM Counter WHERE x < 1000
)
INSERT INTO Posts (Title, Content, Images, Created, IsCommentable, IsFlagged, Author, AuthorID, AuthorAvatar)
SELECT
    'Discussion: ' || c.Name || ' #' || x,
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque vel sem eget justo consequat convallis. Integer porta purus at egestas tincidunt.

     Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Curabitur id nunc id nulla dapibus fermentum.',
    'noimage',
  -- This calculates a pseudo-random date within the year 2025, capped at the current date.
    MIN(datetime('2025-01-01', '+' || (x % 365) || ' days'), CURRENT_DATE),
    (x % 2),
    0,
    iu.Username,
    iu.UserID,
    iu.Avatar
FROM Counter
-- This maps posts to channels in a round-robin fashion.
JOIN Channels c ON c.ID = ((x - 1) % 25) + 177
JOIN IndexedUsers iu ON ((x - 1) % 5) + 1 = iu.rn;

-- Link each post to a channel
INSERT INTO PostChannels (PostID, ChannelID, Created)
SELECT p.ID, ((p.ID - 1) % 25) + 177, MIN(p.Created, CURRENT_DATE)
FROM Posts p;

COMMIT;
