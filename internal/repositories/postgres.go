package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"service-for-working-with-posts-on-graphql/internal/models"
)

type PgRepository struct {
	client *pgxpool.Pool
}

func NewPgRepository(client *pgxpool.Pool) PgRepository {
	return PgRepository{client: client}
}

func (r PgRepository) GetPosts(ctx context.Context) ([]*models.Post, error) {
	rows, err := r.client.Query(ctx, getPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CommentsEnabled); err != nil {
			return nil, fmt.Errorf("Error when receiving posts: %w", err)
		}
		posts = append(posts, &p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r PgRepository) GetPostByID(ctx context.Context, postID string) (*models.Post, error) {
	var post models.Post
	err := r.client.QueryRow(ctx, getPostById, postID).Scan(&post.ID, &post.Title, &post.Content, &post.CommentsEnabled)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("No post found with ID %s", postID)
		}
		return nil, fmt.Errorf("Error when receiving post %s: %w", postID, err)
	}
	return &post, nil
}

func (r PgRepository) GetCommentsByPost(ctx context.Context, postID string, page int, pageSize int) ([]*models.Comment, error) {
	var rows pgx.Rows
	var err error
	if page == -1 && pageSize == -1 {
		rows, err = r.client.Query(ctx, "SELECT id, content, post_id, parent_id FROM comments WHERE post_id = $1", postID)
	} else {
		offset := (page - 1) * pageSize
		rows, err = r.client.Query(ctx, getCommentsByPost, postID, pageSize, offset)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("No comments found for post ID %s", postID)
		}
		return nil, fmt.Errorf("Error when receiving post %s comments: %w", postID, err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.PostID, &c.ParentID); err != nil {
			return nil, fmt.Errorf("Error scanning the comment post %s: %w", postID, err)
		}
		comments = append(comments, &c)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Failed to fetch comments for post %s: %w", postID, err)
	}
	return comments, nil
}

func (r PgRepository) CreatePost(ctx context.Context, title, content string, commentsEnabled bool) (*models.Post, error) {
	post := &models.Post{}

	err := r.client.QueryRow(ctx, createPost, title, content, commentsEnabled).Scan(&post.ID, &post.Title, &post.Content, &post.CommentsEnabled)
	if err != nil {
		return nil, fmt.Errorf("Error creating post %s: %w", title, err)
	}

	return post, nil
}

func (r PgRepository) CreateComment(ctx context.Context, postID, parentID, content string) (*models.Comment, error) {
	if len(content) > 2000 {
		return nil, errors.New("The content of the comment exceeds 2000 characters")
	}

	var commentsEnabled bool
	err := r.client.QueryRow(ctx, "SELECT comments_enabled FROM posts WHERE id = $1", postID).Scan(&commentsEnabled)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("No post found with ID %s, unable to add comment", postID)
		}
		return nil, fmt.Errorf("Database error while checking if comments are enabled for post %s: %w", postID, err)
	}
	if !commentsEnabled {
		return nil, fmt.Errorf("comments are disabled for this post with ID %s", postID)
	}

	var comment models.Comment
	err = r.client.QueryRow(ctx, createComment, content, postID, parentID).Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.ParentID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert comment for post %s: %w", postID, err)
	}

	return &comment, nil
}

func (r PgRepository) GetChildrenComments(ctx context.Context, parentID string) ([]*models.Comment, error) {
	var children []*models.Comment
	rows, err := r.client.Query(ctx, getChildren, parentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("No child comments found for parent ID %s", parentID)
		}
		return nil, fmt.Errorf("Error when fetching child comments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.PostID, &c.ParentID); err != nil {
			return nil, fmt.Errorf("error scanning child comment: %w", err)
		}
		children = append(children, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating child comments: %w", err)
	}

	return children, nil
}
