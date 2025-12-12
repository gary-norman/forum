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

-- Trigger: Update LastActive when a new message is inserted
CREATE TRIGGER IF NOT EXISTS chats_lastactive_trigger
AFTER INSERT ON Messages
FOR EACH ROW
BEGIN
    UPDATE Chats SET LastActive = CURRENT_TIMESTAMP WHERE ID = NEW.ChatID;
END;

-- Indexes for Chats
CREATE INDEX IF NOT EXISTS idx_chats_type ON Chats(Type);
CREATE INDEX IF NOT EXISTS idx_chats_buddyid ON Chats(BuddyID);
CREATE INDEX IF NOT EXISTS idx_chats_groupid ON Chats(GroupID);
CREATE INDEX IF NOT EXISTS idx_chats_lastactive ON Chats(LastActive);

-- Indexes for Messages
CREATE INDEX IF NOT EXISTS idx_messages_chatid ON Messages(ChatID);
CREATE INDEX IF NOT EXISTS idx_messages_userid ON Messages(UserID);
CREATE INDEX IF NOT EXISTS idx_messages_created ON Messages(Created);
CREATE INDEX IF NOT EXISTS idx_messages_chatid_created ON Messages(ChatID, Created);

-- Indexes for ChatUsers
CREATE INDEX IF NOT EXISTS idx_chatusers_userid ON ChatUsers(UserID);
CREATE INDEX IF NOT EXISTS idx_chatusers_chatid ON ChatUsers(ChatID);

PRAGMA foreign_keys = ON;
