package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"godo/src/utils"
)

func TestSingleInstance_TryLock(t *testing.T) {
	// Create a single instance manager
	si := utils.NewSingleInstance("test-app-lock")

	// Try to acquire the lock
	locked, err := si.TryLock()
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	if !locked {
		t.Fatal("Expected to acquire lock, but failed")
	}

	// Verify lock file exists
	lockPath := si.GetLockPath()
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Fatalf("Lock file does not exist at %s", lockPath)
	}

	// Clean up
	defer si.Unlock()
}

func TestSingleInstance_MultipleInstances(t *testing.T) {
	// Create first instance
	si1 := utils.NewSingleInstance("test-multi-app")
	locked1, err := si1.TryLock()
	if err != nil {
		t.Fatalf("Failed to acquire first lock: %v", err)
	}
	if !locked1 {
		t.Fatal("Expected to acquire first lock, but failed")
	}
	defer si1.Unlock()

	// Try to create second instance
	si2 := utils.NewSingleInstance("test-multi-app")
	locked2, err := si2.TryLock()
	if err != nil {
		t.Fatalf("Unexpected error on second lock attempt: %v", err)
	}
	if locked2 {
		t.Fatal("Expected second lock to fail, but it succeeded")
	}
}

func TestSingleInstance_Unlock(t *testing.T) {
	// Create instance and acquire lock
	si := utils.NewSingleInstance("test-unlock-app")
	locked, err := si.TryLock()
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	if !locked {
		t.Fatal("Expected to acquire lock, but failed")
	}

	lockPath := si.GetLockPath()

	// Unlock
	err = si.Unlock()
	if err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}

	// Verify lock file is removed
	if _, err := os.Stat(lockPath); !os.IsNotExist(err) {
		t.Fatalf("Lock file still exists at %s after unlock", lockPath)
	}

	// Try to acquire lock again (should succeed)
	si2 := utils.NewSingleInstance("test-unlock-app")
	locked2, err := si2.TryLock()
	if err != nil {
		t.Fatalf("Failed to acquire lock after unlock: %v", err)
	}
	if !locked2 {
		t.Fatal("Expected to acquire lock after unlock, but failed")
	}
	defer si2.Unlock()
}

func TestSingleInstance_StaleLock(t *testing.T) {
	// Create a fake stale lock file with invalid PID
	lockPath := filepath.Join(os.TempDir(), "test-stale-app.lock")

	// Write a very high PID that likely doesn't exist
	err := os.WriteFile(lockPath, []byte("999999999"), 0600)
	if err != nil {
		t.Fatalf("Failed to create fake lock file: %v", err)
	}

	// Try to acquire lock - should detect stale lock and succeed
	si := utils.NewSingleInstance("test-stale-app")
	locked, err := si.TryLock()
	if err != nil {
		t.Fatalf("Failed to acquire lock with stale lock present: %v", err)
	}
	if !locked {
		t.Fatal("Expected to acquire lock (stale lock should be removed), but failed")
	}

	defer si.Unlock()
}

func TestSingleInstance_LockPath(t *testing.T) {
	si := utils.NewSingleInstance("test-path-app")
	lockPath := si.GetLockPath()

	// Verify lock path contains temp directory and app name
	if lockPath == "" {
		t.Fatal("Lock path is empty")
	}

	expectedDir := os.TempDir()
	if filepath.Dir(lockPath) != expectedDir {
		t.Fatalf("Expected lock path to be in %s, got %s", expectedDir, filepath.Dir(lockPath))
	}

	expectedName := "test-path-app.lock"
	if filepath.Base(lockPath) != expectedName {
		t.Fatalf("Expected lock file name to be %s, got %s", expectedName, filepath.Base(lockPath))
	}
}
