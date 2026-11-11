package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/goleak"
)

func TestNameAccumulation(t *testing.T) {
	goleak.VerifyNone(t)

	fmt.Println("aaa")
	a := NameAccumulation("/mnt/test_1_1", "/")
	fmt.Println(a)
}

func TestSoftDelete(t *testing.T) {
	goleak.VerifyNone(t)

	// Create temp test directory
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Soft delete it
	if err := SoftDelete(testFile); err != nil {
		t.Fatalf("SoftDelete failed: %v", err)
	}

	// Verify original is gone
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("original file should be removed after SoftDelete")
	}

	// Verify trash directory was created
	trashDir := filepath.Join(tmpDir, ".casaos-trash")
	if _, err := os.Stat(trashDir); os.IsNotExist(err) {
		t.Error("trash directory should be created")
	}

	// Verify file exists in trash
	entries, _ := os.ReadDir(trashDir)
	if len(entries) != 1 {
		t.Errorf("expected 1 file in trash, got %d", len(entries))
	}
}

func TestSoftDeleteDirectory(t *testing.T) {
	goleak.VerifyNone(t)

	tmpDir := t.TempDir()

	// Create a test directory with files
	testDir := filepath.Join(tmpDir, "testdir")
	os.MkdirAll(testDir, 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("content2"), 0644)

	// Soft delete the directory
	if err := SoftDelete(testDir); err != nil {
		t.Fatalf("SoftDelete failed: %v", err)
	}

	// Verify original is gone
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Error("original directory should be removed after SoftDelete")
	}

	// Verify trash contains the directory
	trashDir := filepath.Join(tmpDir, ".casaos-trash")
	entries, _ := os.ReadDir(trashDir)
	if len(entries) != 1 {
		t.Errorf("expected 1 entry in trash, got %d", len(entries))
	}
}
