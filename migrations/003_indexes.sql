BEGIN TRANSACTION;

-- Users
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON Users(Username);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_avatar ON Users(Avatar);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON Users(EmailAddress);
CREATE INDEX IF NOT EXISTS idx_users_sessiontoken ON Users(SessionToken);
CREATE INDEX IF NOT EXISTS idx_users_csrftoken ON Users(CsrfToken);
CREATE INDEX IF NOT EXISTS idx_users_created ON Users(Created);
CREATE INDEX IF NOT EXISTS idx_users_updated ON Users(Updated);

-- Channels
CREATE UNIQUE INDEX IF NOT EXISTS idx_channels_name ON Channels(Name);
CREATE INDEX IF NOT EXISTS idx_channels_owner ON Channels(OwnerID);
CREATE INDEX IF NOT EXISTS idx_channels_created ON Channels(Created);
CREATE INDEX IF NOT EXISTS idx_channels_updated ON Channels(Updated);

-- Posts
CREATE INDEX IF NOT EXISTS idx_posts_authorid ON Posts(AuthorID);
CREATE INDEX IF NOT EXISTS idx_posts_created ON Posts(Created);
CREATE INDEX IF NOT EXISTS idx_posts_updated ON Posts(Updated);

-- Comments
CREATE INDEX IF NOT EXISTS idx_comments_postid_author ON Comments(CommentedPostID, AuthorID);
CREATE INDEX IF NOT EXISTS idx_comments_commentid ON Comments(CommentedCommentID);
CREATE INDEX IF NOT EXISTS idx_comments_authorid ON Comments(AuthorID);
CREATE INDEX IF NOT EXISTS idx_comments_created ON Comments(Created);
CREATE INDEX IF NOT EXISTS idx_comments_updated ON Comments(Updated);

-- Reactions
CREATE UNIQUE INDEX IF NOT EXISTS idx_reactions_post ON Reactions(AuthorID, ReactedPostID) WHERE ReactedPostID IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_reactions_comment ON Reactions(AuthorID, ReactedCommentID) WHERE ReactedCommentID IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_reactions_author ON Reactions(AuthorID);

-- Flags
CREATE INDEX IF NOT EXISTS idx_flags_post ON Flags(FlaggedPostID);
CREATE INDEX IF NOT EXISTS idx_flags_comment ON Flags(FlaggedCommentID);
CREATE INDEX IF NOT EXISTS idx_flags_user ON Flags(FlaggedUserID);
CREATE INDEX IF NOT EXISTS idx_flags_author ON Flags(AuthorID);
CREATE INDEX IF NOT EXISTS idx_flags_channel ON Flags(ChannelID);

-- PostChannels
CREATE INDEX IF NOT EXISTS idx_postchannels_postid ON PostChannels(PostID);
CREATE INDEX IF NOT EXISTS idx_postchannels_channelid ON PostChannels(ChannelID);

-- Memberships & Mods
CREATE INDEX IF NOT EXISTS idx_memberships_userid ON Memberships(UserID);
CREATE INDEX IF NOT EXISTS idx_memberships_channelid ON Memberships(ChannelID);
CREATE INDEX IF NOT EXISTS idx_mods_userid ON Mods(UserID);
CREATE INDEX IF NOT EXISTS idx_mods_channelid ON Mods(ChannelID);

-- Followers & Following
CREATE INDEX IF NOT EXISTS idx_followers_followerid ON Followers(FollowerUserID);
CREATE INDEX IF NOT EXISTS idx_following_userid ON Following(UserID);
CREATE INDEX IF NOT EXISTS idx_following_followingid ON Following(FollowingUserID);

-- MutedChannels
CREATE INDEX IF NOT EXISTS idx_mutedchannels_userid ON MutedChannels(UserID);
CREATE INDEX IF NOT EXISTS idx_mutedchannels_channelid ON MutedChannels(ChannelID);

-- Images
CREATE INDEX IF NOT EXISTS idx_images_postid ON Images(PostID);
CREATE INDEX IF NOT EXISTS idx_images_authorid ON Images(AuthorID);

-- Notifications & NotificationsUsers
CREATE INDEX IF NOT EXISTS idx_notifications_read ON Notifications(Read);
CREATE INDEX IF NOT EXISTS idx_notifications_archived ON Notifications(Archived);
CREATE INDEX IF NOT EXISTS idx_notificationsusers_userid ON NotificationsUsers(UserID);
CREATE INDEX IF NOT EXISTS idx_notificationsusers_notifid ON NotificationsUsers(NotificationID);

-- ChannelsRules
CREATE UNIQUE INDEX IF NOT EXISTS idx_channelsrules_unique ON ChannelsRules(ChannelID, RuleID);

-- Chats
CREATE INDEX IF NOT EXISTS idx_chats_type ON Chats(Type);
CREATE INDEX IF NOT EXISTS idx_chats_buddyid ON Chats(BuddyID);
CREATE INDEX IF NOT EXISTS idx_chats_groupid ON Chats(GroupID);
CREATE INDEX IF NOT EXISTS idx_chats_lastactive ON Chats(LastActive);

-- Messages
CREATE INDEX IF NOT EXISTS idx_messages_chatid ON Messages(ChatID);
CREATE INDEX IF NOT EXISTS idx_messages_userid ON Messages(UserID);
CREATE INDEX IF NOT EXISTS idx_messages_created ON Messages(Created);
CREATE INDEX IF NOT EXISTS idx_messages_chatid_created ON Messages(ChatID, Created);

-- ChatUsers
CREATE UNIQUE INDEX IF NOT EXISTS idx_chatusers_unique ON ChatUsers(ChatID, UserID);
CREATE INDEX IF NOT EXISTS idx_chatusers_userid ON ChatUsers(UserID);
CREATE INDEX IF NOT EXISTS idx_chatusers_chatid ON ChatUsers(ChatID);

COMMIT;

