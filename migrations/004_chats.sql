PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;

-- Chats table: stores both buddy (1-on-1) and group chats
CREATE TABLE IF NOT EXISTS Chats (
    ID BLOB PRIMARY KEY,
    Type TEXT NOT NULL CHECK (Type IN ('buddy', 'group')),
    Name TEXT,
    GroupID BLOB,
    BuddyID BLOB,
    Created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    LastActive DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (BuddyID) REFERENCES Users(ID) ON DELETE SET NULL,
    CHECK (
        (Type = 'buddy' AND BuddyID IS NOT NULL AND GroupID IS NULL) OR
        (Type = 'group' AND GroupID IS NOT NULL AND BuddyID IS NULL)
    )
);

-- Messages table: stores individual chat messages
CREATE TABLE IF NOT EXISTS Messages (
    ID BLOB PRIMARY KEY,
    ChatID BLOB NOT NULL,
    UserID BLOB,
    Content TEXT NOT NULL,
    Created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ChatID) REFERENCES Chats(ID) ON DELETE CASCADE,
    FOREIGN KEY (UserID) REFERENCES Users(ID) ON DELETE SET NULL
);

-- ChatUsers table: many-to-many relationship between users and chats
CREATE TABLE IF NOT EXISTS ChatUsers (
    ID INTEGER PRIMARY KEY,
    ChatID BLOB NOT NULL,
    UserID BLOB NOT NULL,
    Created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(ChatID, UserID),
    FOREIGN KEY (ChatID) REFERENCES Chats(ID) ON DELETE CASCADE,
    FOREIGN KEY (UserID) REFERENCES Users(ID) ON DELETE CASCADE
);

COMMIT;

PRAGMA foreign_keys = ON;
