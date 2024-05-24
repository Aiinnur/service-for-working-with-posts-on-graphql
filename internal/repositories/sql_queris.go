package repositories

const (
	CreatePosts = `
    CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        comments_enabled BOOLEAN NOT NULL DEFAULT true
    );`

	CreateComments = `
    CREATE TABLE IF NOT EXISTS comments (
        id SERIAL PRIMARY KEY,
        content TEXT NOT NULL CHECK (length(content) <= 2000),
        post_id INTEGER NOT NULL,
        parent_id INTEGER,
        FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
        FOREIGN KEY (parent_id) REFERENCES comments (id) ON DELETE CASCADE
    );`

	getPosts = "SELECT id, title, content, comments_enabled FROM posts;"

	getPostById = "SELECT id, title, content, comments_enabled FROM posts WHERE id = $1;"

	getCommentsByPost = `SELECT id, content, post_id, parent_id
						 FROM comments
					     WHERE post_id = $1
						 LIMIT $2 OFFSET $3;`

	createPost = `
			INSERT INTO posts (title, content, comments_enabled)
			VALUES ($1, $2, $3)
			RETURNING id, title, content, comments_enabled;`

	createComment = `INSERT INTO comments (content, post_id, parent_id) 
					 VALUES ($1, $2, $3) 
					 RETURNING id, content, post_id, parent_id`

	getChildren = "SELECT id, content, post_id, parent_id FROM comments WHERE parent_id = $1"
)
