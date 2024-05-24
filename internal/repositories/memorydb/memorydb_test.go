package memorydb

import (
	"context"
	"service-for-working-with-posts-on-graphql/internal/models"
	"testing"
)

func TestCreatePost(t *testing.T) {
	repo := NewMemoryRepository()
	post, err := repo.CreatePost(context.Background(), "Test Post", "Test post.", true)
	if err != nil {
		t.Fatalf("CreatePost() error = %v, wantErr nil", err)
	}
	if post == nil {
		t.Fatal("CreatePost() post is nil, want non-nil")
	}
	if post.Title != "Test Post" || post.Content != "Test post." || !post.CommentsEnabled {
		t.Errorf("CreatePost() got = %v, want = %v", post, &models.Post{Title: "Test Post", Content: "This is a test post.", CommentsEnabled: true})
	}
}

func TestGetPosts(t *testing.T) {
	repo := NewMemoryRepository()
	_, _ = repo.CreatePost(context.Background(), "First Post", "Content", true)
	_, _ = repo.CreatePost(context.Background(), "Second Post", "Content", true)

	posts, err := repo.GetPosts(context.Background())
	if err != nil {
		t.Fatalf("GetPosts() error = %v, wantErr nil", err)
	}
	if len(posts) != 2 {
		t.Fatalf("GetPosts() got = %d, want = %d", len(posts), 2)
	}
}

func TestCreateComment(t *testing.T) {
	repo := NewMemoryRepository()
	post, _ := repo.CreatePost(context.Background(), "New Post", "Content", true)
	comment, err := repo.CreateComment(context.Background(), post.ID, "", "This is a comment.")
	if err != nil {
		t.Fatalf("CreateComment() error = %v, wantErr nil", err)
	}
	if comment == nil {
		t.Fatal("CreateComment() comment is nil, want non-nil")
	}
	if comment.Content != "This is a comment." {
		t.Errorf("CreateComment() got = %v, want = %v", comment.Content, "This is a comment.")
	}
}

func TestGetCommentsByPost(t *testing.T) {
	repo := NewMemoryRepository()
	post, _ := repo.CreatePost(context.Background(), "Post", "Content of the post.", true)
	_, _ = repo.CreateComment(context.Background(), post.ID, "", "First comment")
	_, _ = repo.CreateComment(context.Background(), post.ID, "", "Second comment")

	comments, err := repo.GetCommentsByPost(context.Background(), post.ID, 1, 10)
	if err != nil {
		t.Fatalf("GetCommentsByPost() error = %v, wantErr nil", err)
	}
	if len(comments) != 2 {
		t.Fatalf("GetCommentsByPost() got = %d, want = %d", len(comments), 2)
	}
}

func TestGetPostByID(t *testing.T) {
	repo := NewMemoryRepository()
	newPost, _ := repo.CreatePost(context.Background(), "Test Post", "Test Content", true)
	post, err := repo.GetPostByID(context.Background(), newPost.ID)
	if err != nil {
		t.Errorf("GetPostByID() error = %v, wantErr nil", err)
	}
	if post == nil || post.ID != newPost.ID {
		t.Errorf("GetPostByID() got = %v, want %v", post, newPost)
	}

	_, err = repo.GetPostByID(context.Background(), "nil")
	if err == nil {
		t.Errorf("GetPostByID() expected error, got nil")
	}
}

func TestGetChildrenComments(t *testing.T) {
	repo := NewMemoryRepository()
	post, _ := repo.CreatePost(context.Background(), "Title", "Content", true)
	parentComment, _ := repo.CreateComment(context.Background(), post.ID, "", "Parent comment")
	_, _ = repo.CreateComment(context.Background(), post.ID, parentComment.ID, "Child comment 1")
	_, _ = repo.CreateComment(context.Background(), post.ID, parentComment.ID, "Child comment 2")

	children, err := repo.GetChildrenComments(context.Background(), parentComment.ID)
	if err != nil {
		t.Errorf("GetChildrenComments() error = %v, wantErr nil", err)
	}

	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}

	for _, child := range children {
		if child.ParentID != parentComment.ID {
			t.Errorf("Expected parent ID %s, got %s", parentComment.ID, child.ParentID)
		}
	}
}
