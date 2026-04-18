package local

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestProviderPutSignedURLDelete(t *testing.T) {
	basePath := t.TempDir()
	provider, err := New(basePath, "/static")
	if err != nil {
		t.Fatalf("new provider failed: %v", err)
	}

	url, err := provider.Put(context.Background(), "uploads/test.txt", strings.NewReader("hello"), 5, "text/plain")
	if err != nil {
		t.Fatalf("put failed: %v", err)
	}
	if url != "/static/uploads/test.txt" {
		t.Fatalf("unexpected url: %s", url)
	}

	path := filepath.Join(basePath, "uploads", "test.txt")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	if string(content) != "hello" {
		t.Fatalf("unexpected file content: %s", string(content))
	}

	signed, err := provider.SignedURL(context.Background(), "uploads/test.txt", time.Hour)
	if err != nil {
		t.Fatalf("signed url failed: %v", err)
	}
	if signed != "/static/uploads/test.txt" {
		t.Fatalf("unexpected signed url: %s", signed)
	}

	if err := provider.Delete(context.Background(), "uploads/test.txt"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected file removed")
	}
}

func TestProviderRejectsTraversal(t *testing.T) {
	provider, err := New(t.TempDir(), "/static")
	if err != nil {
		t.Fatalf("new provider failed: %v", err)
	}
	_, err = provider.Put(context.Background(), "../evil.txt", strings.NewReader("x"), 1, "text/plain")
	if err == nil {
		t.Fatalf("expected traversal to fail")
	}

}
