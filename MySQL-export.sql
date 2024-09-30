-- Exported from QuickDBD: https://www.quickdatabasediagrams.com/
-- NOTE! If you have used non-SQL datatypes in your design, you will have to change these here.

-- Modify this code to update the DB schema diagram.
-- To reset the sample schema, replace everything with
-- two dots ('..' - without quotes).

CREATE TABLE `user` (
    `ID` INT AUTO_INCREMENT NOT NULL ,
    `Username` VARCHAR(30)  NOT NULL ,
    `Password` VARCHAR(20)  NOT NULL ,
    `Email_address` VARCHAR(30)  NOT NULL ,
    `Avatar` UUID  NULL ,
    `Banner` UUID  NULL ,
    `Desc` VARCHAR(255)  NULL ,
    `Type` VARCHAR(10)  NOT NULL ,
    `Created` DATETIME  NOT NULL DEFAULT getutcdate(),
    `Membership` JSON  NOT NULL ,
    `Followers` JSON  NOT NULL ,
    `Following` JSON  NOT NULL ,
    `Bookmarks_user` JSON  NOT NULL ,
    `Bookmarks_channel` JSON  NOT NULL ,
    `Bookmarks_post` JSON  NOT NULL ,
    `Bookmarks_comment` JSON  NOT NULL ,
    `Posts` JSON  NOT NULL ,
    `Comments` JSON  NOT NULL ,
    `Flagged_items` JSON  NOT NULL ,
    `Mod_of` JSON  NOT NULL ,
    `Reactions` JSON  NOT NULL ,
    `Is_flagged` JSON  NOT NULL ,
    `Is_flagged_approved` JSON  NOT NULL ,
    PRIMARY KEY (
        `ID`
    )
);

CREATE TABLE `channel` (
    `ID` INT AUTO_INCREMENT NOT NULL ,
    `Name` VARCHAR(30)  NOT NULL ,
    `Avatar` UUID  NULL ,
    `Desc` VARCHAR(555)  NULL ,
    `Created` DATETIME  NOT NULL DEFAULT getutcdate(),
    `Rules` VARCHAR  NULL ,
    `Privacy` BOOL  NOT NULL ,
    `Members` JSON  NOT NULL ,
    `Mods` JSON  NOT NULL ,
    `Posts` JSON  NOT NULL ,
    PRIMARY KEY (
        `ID`
    )
);

CREATE TABLE `post` (
    `ID` INT AUTO_INCREMENT NOT NULL ,
    `Title` VARCHAR(80)  NOT NULL ,
    `Content` VARCHAR  NOT NULL ,
    `Images` JSON  NOT NULL ,
    `Created` DATETIME  NOT NULL DEFAULT getutcdate(),
    `Commentable` BOOL  NOT NULL ,
    `Author` VARCHAR  NOT NULL ,
    `Reactions` JSON  NOT NULL ,
    `Comments` JSON  NOT NULL ,
    `Channel` INT  NOT NULL ,
    `Is_flagged` JSON  NOT NULL ,
    `Is_flagged_approved` JSON  NOT NULL ,
    PRIMARY KEY (
        `ID`
    )
);

CREATE TABLE `comment` (
    `ID` INT AUTO_INCREMENT NOT NULL ,
    `Content` VARCHAR  NOT NULL ,
    `Images` JSON  NOT NULL ,
    `Created` DATETIME  NOT NULL DEFAULT getutcdate(),
    `Author` INT  NOT NULL ,
    `Reactions` JSON  NOT NULL ,
    `Replies` JSON  NOT NULL ,
    `Channel` INT  NOT NULL ,
    `Is_flagged` JSON  NOT NULL ,
    `Is_flagged_approved` JSON  NOT NULL ,
    -- figure out how to add multiple
    `Parent_post` INT  NOT NULL ,
    `Parent_comment` INT  NOT NULL ,
    PRIMARY KEY (
        `ID`
    )
);

CREATE TABLE `reaction` (
    `ID` INT AUTO_INCREMENT NOT NULL ,
    `Liked` BOOL  NOT NULL ,
    `Disliked` BOOL  NOT NULL ,
    `Created` DATETIME  NOT NULL DEFAULT getutcdate(),
    `Author` INT  NOT NULL ,
    `Channel` INT  NOT NULL ,
    -- figure out how to add multiple
    `Parent_post` JSON  NOT NULL ,
    `Parent_comment` JSON  NOT NULL ,
    -- figure out how to add multiple
    `Parent_post_author` INT  NOT NULL ,
    `Parent_comment_author` INT  NOT NULL ,
    PRIMARY KEY (
        `ID`
    )
);

CREATE TABLE `flag` (
    `ID` INT AUTO_INCREMENT NOT NULL ,
    `Type` VARCHAR  NOT NULL ,
    `Content` VARCHAR(160)  NOT NULL ,
    `Created` DATETIME  NOT NULL DEFAULT getutcdate(),
    `Approved` BOOL  NOT NULL ,
    `Author` INT  NOT NULL ,
    -- figure out how to add multiple
    `Parent_post` JSON  NOT NULL ,
    `Parent_comment` JSON  NOT NULL ,
    -- figure out how to add multiple
    `Parent_post_author` INT  NOT NULL ,
    `Parent_comment_author` INT  NOT NULL ,
    PRIMARY KEY (
        `ID`
    )
);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Membership` FOREIGN KEY(`Membership`)
REFERENCES `channel` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Followers` FOREIGN KEY(`Followers`)
REFERENCES `user` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Following` FOREIGN KEY(`Following`)
REFERENCES `user` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Bookmarks_user` FOREIGN KEY(`Bookmarks_user`)
REFERENCES `user` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Bookmarks_channel` FOREIGN KEY(`Bookmarks_channel`)
REFERENCES `channel` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Bookmarks_post` FOREIGN KEY(`Bookmarks_post`)
REFERENCES `post` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Bookmarks_comment` FOREIGN KEY(`Bookmarks_comment`)
REFERENCES `comment` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Posts` FOREIGN KEY(`Posts`)
REFERENCES `post` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Comments` FOREIGN KEY(`Comments`)
REFERENCES `comment` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Flagged_items_Is_flagged_approved` FOREIGN KEY(`Flagged_items`, `Is_flagged_approved`)
REFERENCES `flag` (`ID`, `Approved`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Mod_of` FOREIGN KEY(`Mod_of`)
REFERENCES `channel` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Reactions` FOREIGN KEY(`Reactions`)
REFERENCES `reaction` (`ID`);

ALTER TABLE `user` ADD CONSTRAINT `fk_user_Is_flagged` FOREIGN KEY(`Is_flagged`)
REFERENCES `flag` (`ID`);

ALTER TABLE `channel` ADD CONSTRAINT `fk_channel_Members` FOREIGN KEY(`Members`)
REFERENCES `user` (`ID`);

ALTER TABLE `channel` ADD CONSTRAINT `fk_channel_Mods` FOREIGN KEY(`Mods`)
REFERENCES `user` (`ID`);

ALTER TABLE `channel` ADD CONSTRAINT `fk_channel_Posts` FOREIGN KEY(`Posts`)
REFERENCES `post` (`ID`);

ALTER TABLE `post` ADD CONSTRAINT `fk_post_Author` FOREIGN KEY(`Author`)
REFERENCES `user` (`ID`);

ALTER TABLE `post` ADD CONSTRAINT `fk_post_Reactions` FOREIGN KEY(`Reactions`)
REFERENCES `reaction` (`ID`);

ALTER TABLE `post` ADD CONSTRAINT `fk_post_Comments` FOREIGN KEY(`Comments`)
REFERENCES `comment` (`ID`);

ALTER TABLE `post` ADD CONSTRAINT `fk_post_Channel` FOREIGN KEY(`Channel`)
REFERENCES `channel` (`ID`);

ALTER TABLE `post` ADD CONSTRAINT `fk_post_Is_flagged_Is_flagged_approved` FOREIGN KEY(`Is_flagged`, `Is_flagged_approved`)
REFERENCES `flag` (`ID`, `Approved`);

ALTER TABLE `comment` ADD CONSTRAINT `fk_comment_Author` FOREIGN KEY(`Author`)
REFERENCES `user` (`ID`);

ALTER TABLE `comment` ADD CONSTRAINT `fk_comment_Reactions` FOREIGN KEY(`Reactions`)
REFERENCES `reaction` (`ID`);

ALTER TABLE `comment` ADD CONSTRAINT `fk_comment_Replies` FOREIGN KEY(`Replies`)
REFERENCES `comment` (`ID`);

ALTER TABLE `comment` ADD CONSTRAINT `fk_comment_Channel` FOREIGN KEY(`Channel`)
REFERENCES `channel` (`ID`);

ALTER TABLE `comment` ADD CONSTRAINT `fk_comment_Is_flagged_Is_flagged_approved` FOREIGN KEY(`Is_flagged`, `Is_flagged_approved`)
REFERENCES `flag` (`ID`, `Approved`);

ALTER TABLE `comment` ADD CONSTRAINT `fk_comment_Parent_post` FOREIGN KEY(`Parent_post`)
REFERENCES `post` (`ID`);

ALTER TABLE `comment` ADD CONSTRAINT `fk_comment_Parent_comment` FOREIGN KEY(`Parent_comment`)
REFERENCES `comment` (`ID`);

ALTER TABLE `reaction` ADD CONSTRAINT `fk_reaction_Author` FOREIGN KEY(`Author`)
REFERENCES `user` (`ID`);

ALTER TABLE `reaction` ADD CONSTRAINT `fk_reaction_Channel` FOREIGN KEY(`Channel`)
REFERENCES `channel` (`ID`);

ALTER TABLE `reaction` ADD CONSTRAINT `fk_reaction_Parent_post_Parent_post_author` FOREIGN KEY(`Parent_post`, `Parent_post_author`)
REFERENCES `post` (`ID`, `Author`);

ALTER TABLE `reaction` ADD CONSTRAINT `fk_reaction_Parent_comment_Parent_comment_author` FOREIGN KEY(`Parent_comment`, `Parent_comment_author`)
REFERENCES `comment` (`ID`, `Author`);

ALTER TABLE `flag` ADD CONSTRAINT `fk_flag_Author` FOREIGN KEY(`Author`)
REFERENCES `user` (`ID`);

ALTER TABLE `flag` ADD CONSTRAINT `fk_flag_Parent_post_Parent_post_author` FOREIGN KEY(`Parent_post`, `Parent_post_author`)
REFERENCES `post` (`ID`, `Author`);

ALTER TABLE `flag` ADD CONSTRAINT `fk_flag_Parent_comment_Parent_comment_author` FOREIGN KEY(`Parent_comment`, `Parent_comment_author`)
REFERENCES `comment` (`ID`, `Author`);

CREATE INDEX `idx_user_Username`
ON `user` (`Username`);

CREATE INDEX `idx_channel_Name`
ON `channel` (`Name`);

