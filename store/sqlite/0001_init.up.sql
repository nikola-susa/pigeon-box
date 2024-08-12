
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY,
    slack_id TEXT NOT NULL,
    name TEXT DEFAULT NULL,
    username TEXT DEFAULT NULL,
    avatar TEXT DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(slack_id) ON CONFLICT REPLACE
);

CREATE TABLE IF NOT EXISTS thread (
    id INTEGER PRIMARY KEY,
    slack_id TEXT DEFAULT NULL,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT DEFAULT NULL,
    slack_timestamp TEXT DEFAULT NULL,
    key TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME DEFAULT NULL,
    messages_expire_at DATETIME DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS message (
    "id" INTEGER PRIMARY KEY,
    "thread_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "file_id" INTEGER DEFAULT NULL,
    "text" TEXT DEFAULT NULL,
    "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "expires_at" DATETIME DEFAULT NULL,
    FOREIGN KEY (thread_id) REFERENCES thread(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES file(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS file (
    "id" INTEGER PRIMARY KEY,
    "thread_id" INTEGER NOT NULL,
    "user_id" INTEGER DEFAULT NULL,
    "path" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "content_type" TEXT NOT NULL,
    "size" INTEGER NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (thread_id) REFERENCES thread(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS session (
    "id"         INTEGER PRIMARY KEY,
    "user_id"    INTEGER  NOT NULL,
    "thread_id"  INTEGER  NOT NULL,
    "expires_at" DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (thread_id) REFERENCES thread(id) ON DELETE CASCADE
);
