package memorydb

import (
	"context"
	"errors"
	"service-for-working-with-posts-on-graphql/internal/models"
	"strconv"
	"sync"
)

type MemoryRepository struct {
	posts          map[string]*models.Post
	comments       map[string]*models.Comment
	postCounter    int
	commentCounter int
	mu             sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		posts:    make(map[string]*models.Post),
		comments: make(map[string]*models.Comment),
	}
}

func (r MemoryRepository) GetPosts(ctx context.Context) ([]*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var posts []*models.Post
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r MemoryRepository) GetPostByID(ctx context.Context, postID string) (*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	post, ok := r.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (r MemoryRepository) CreatePost(ctx context.Context, title, content string, commentsEnabled bool) (*models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.postCounter++
	postID := strconv.Itoa(r.postCounter)
	post := &models.Post{
		ID:              postID,
		Title:           title,
		Content:         content,
		CommentsEnabled: commentsEnabled,
	}
	r.posts[post.ID] = post
	return post, nil
}

func (r MemoryRepository) GetCommentsByPost(ctx context.Context, postID string, page int, pageSize int) ([]*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filteredComments []*models.Comment

	for _, comment := range r.comments {
		if comment.PostID == postID {
			filteredComments = append(filteredComments, comment)
		}
	}

	if page == -1 && pageSize == -1 {
		return filteredComments, nil
	}

	if page < 1 || pageSize < 1 {
		return nil, errors.New("invalid values for page or pageSize")
	}

	offset := (page - 1) * pageSize
	if offset >= len(filteredComments) {
		return nil, errors.New("page number out of range")
	}

	end := offset + pageSize
	if end > len(filteredComments) {
		end = len(filteredComments)
	}

	return filteredComments[offset:end], nil
}

func (r MemoryRepository) CreateComment(ctx context.Context, postID, parentID, content string) (*models.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commentCounter++
	commentID := strconv.Itoa(r.commentCounter)
	comment := &models.Comment{
		ID:       commentID,
		Content:  content,
		PostID:   postID,
		ParentID: parentID,
	}
	r.comments[comment.ID] = comment
	return comment, nil
}

func (r MemoryRepository) GetChildrenComments(ctx context.Context, parentID string) ([]*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var comments []*models.Comment
	for _, comment := range r.comments {
		if comment.ParentID == parentID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}
