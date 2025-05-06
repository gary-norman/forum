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

-- Create 25 channels with static names
INSERT INTO Channels (OwnerID, Name, Avatar, Banner, Description, Created, Privacy, IsFlagged, IsMuted)
SELECT
  (SELECT UserID FROM TempUsers ORDER BY RANDOM() LIMIT 1),
  name,
  'noimage',
  'default.png',
  'Community focused on ' || name || '. Join discussions, share ideas, and collaborate with fellow tech enthusiasts.',
  MIN(datetime('2024-01-01', '+' || ROWID || ' days'), CURRENT_DATE),
  (ROWID % 2),
  (ROWID % 2),
  ((ROWID + 1) % 2)
FROM ChannelNames;

-- Insert memberships: Each channel gets the owner + 2 random users
INSERT INTO Memberships (UserID, ChannelID, Created)
SELECT u.UserID, c.ID, MIN(datetime(c.Created, '+' || (ABS(hex(u.UserID) % 5)) || ' hours'), CURRENT_DATE)
FROM Channels c
JOIN TempUsers u ON u.UserID != c.OwnerID
ORDER BY RANDOM()
LIMIT 50;

COMMIT;
