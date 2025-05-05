-- This SQL file populates Channels, Memberships, Posts, and PostChannels with sample data.
-- Assumes 5 existing users in Users table with columns: ID (BLOB), Username (TEXT), Avatar (TEXT)
BEGIN TRANSACTION;

-- Temporary table to cache 5 users
CREATE TEMP TABLE TempUsers AS
SELECT ID AS UserID, Username, Avatar FROM Users LIMIT 5;

-- Descriptive channel names
CREATE TEMP TABLE ChannelNames(name TEXT);
INSERT INTO ChannelNames (name) VALUES
('Go Concurrency Deep Dive'),
('CSS Advanced Techniques'),
('C# for Beginners'),
('Java Memory Management'),
('JavaScript Event Loop Explained'),
('Rust Ownership Patterns'),
('Python Data Science Hub'),
('Node.js Performance Tuning'),
('React State Management'),
('Vue Composition API'),
('Machine Learning in Practice'),
('Docker and Containerization'),
('Kubernetes for Devs'),
('PostgreSQL Performance Tips'),
('SQLite for Mobile Apps'),
('Microservices Architecture'),
('API Security Fundamentals'),
('OAuth2 and Identity'),
('Web Accessibility Matters'),
('TypeScript Type Systems'),
('DevOps CI/CD Strategies'),
('Clean Code Principles'),
('Design Patterns in Go'),
('Functional Programming in JS'),
('SvelteKit for Web Apps');


INSERT INTO Channels (OwnerID, Name, Avatar, Banner, Description, Created, Privacy, IsMuted, IsFlagged)
SELECT
  (SELECT UserID FROM TempUsers ORDER BY RANDOM() LIMIT 1),
  name,
  'noimage',
  'default.png',
  'Community focused on ' || name || '. Join discussions, share ideas, and collaborate with fellow tech enthusiasts.',
  MIN(datetime('2024-01-01', '+' || ROWID || ' days'), CURRENT_DATE),
  (ROWID % 2),
  ((ROWID + 1) % 2),
  (ROWID % 2)
FROM ChannelNames;

-- Insert memberships: Each channel gets the owner + 2 random users
INSERT INTO Memberships (UserID, ChannelID, Created)
SELECT u.UserID, c.ID, MIN(datetime(c.Created, '+' || (ABS(hex(u.UserID) % 5)) || ' hours'), CURRENT_DATE)
FROM Channels c
JOIN TempUsers u ON u.UserID != c.OwnerID
ORDER BY RANDOM()
LIMIT 50;

-- Precompute indexed users
WITH IndexedUsers AS (
  SELECT UserID, Username, Avatar, ROW_NUMBER() OVER () AS rn
  FROM TempUsers
),
Counter(x) AS (
  SELECT 1
  UNION ALL
  SELECT x+1 FROM Counter WHERE x < 1000
)
INSERT INTO Posts (Title, Content, Images, Created, IsCommentable, IsFlagged, Author, AuthorID, AuthorAvatar)
SELECT
  'Discussion: ' || c.Name || ' #' || x,
  'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque vel sem eget justo consequat convallis. Integer porta purus at egestas tincidunt.\n\nVestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Curabitur id nunc id nulla dapibus fermentum.',
  'noimage',
  MIN(datetime('2024-01-01', '+' || (x % 365) || ' days'), CURRENT_DATE),
  (x % 2),
  (ROWID % 2),
  iu.Username,
  iu.UserID,
  iu.Avatar
FROM Counter
JOIN Channels c ON c.ID = ((x - 1) % 25) + 1
JOIN IndexedUsers iu ON ((x - 1) % 5) + 1 = iu.rn;

-- Link each post to a channel
INSERT INTO PostChannels (PostID, ChannelID, Created)
SELECT p.ID, ((p.ID - 1) % 25) + 1, MIN(p.Created, CURRENT_DATE)
FROM Posts p;

COMMIT;
