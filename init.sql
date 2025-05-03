-- Return to this code later
CREATE TABLE sessions (
  session_id UUID PRIMARY KEY, -- cookie/session ID
  avatar_url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NOT NULL
);

CREATE TABLE posts (
  post_id UUID PRIMARY KEY, -- no gen_random_uuid(), our Go app is the one generating uuid
  session_id TEXT NOT NULL REFERENCES sessions(session_id),
  user_name TEXT NOT NULL DEFAULT 'Anonymous',
  post_title TEXT NOT NULL,
  post_content TEXT NULLABLE, -- Allow users to post only title or title and image
  image_urls TEXT[] DEFAULT '{}',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_archived BOOLEAN DEFAULT FALSE
);


CREATE TABLE comments (
  comment_id UUID PRIMARY,
  post_id UUID NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
  session_id TEXT NOT NULL REFERENCES sessions(session_id),
  user_name TEXT NOT NULL DEFAULT 'Anonymous',
  comment_content TEXT,
  parent_comment_id UUID REFERENCES comments(comment_id), -- for replies to comments
  image_urls TEXT [] DEFAULT '{}', -- turns out images can be added to comments as well
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_archived BOOLEAN DEFAULT FALSE
);

-- Indexes
CREATE INDEX idx_posts_session_id ON posts(session_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);
CREATE INDEX idx_comments_post_id ON comments (post_id); -- fetching comments by id
CREATE INDEX idx_comments_parent ON comments (parent_comment_id); -- fetching replies to comments
