package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"godo/src/utils"
)

// TestSingleInstance_TryLock checks first instance lock acquisition
func TestSingleInstance_TryLock(t *testing.T) {
	// Берем отдельное имя, чтобы не зацепить рабочий lock
	si := utils.NewSingleInstance("test-app-lock")

	locked, err := si.TryLock()

	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	if !locked {
		t.Fatal("Expected to acquire lock, but failed")
	}

	defer si.Unlock()

	// После успешного захвата файл должен лежать в temp
	lockPath := si.GetLockPath()
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Fatalf("Lock file does not exist at %s", lockPath)
	}
}

// TestSingleInstance_MultipleInstances checks second instance rejection
func TestSingleInstance_MultipleInstances(t *testing.T) {
	// Первый экземпляр должен занять lock
	si1 := utils.NewSingleInstance("test-multi-app")
	locked1, err := si1.TryLock()

	if err != nil {
		t.Fatalf("Failed to acquire first lock: %v", err)
	}

	if !locked1 {
		t.Fatal("Expected to acquire first lock, but failed")
	}

	defer si1.Unlock()

	// Второй экземпляр с тем же именем должен получить отказ
	si2 := utils.NewSingleInstance("test-multi-app")
	locked2, err := si2.TryLock()

	if err != nil {
		t.Fatalf("Unexpected error on second lock attempt: %v", err)
	}

	if locked2 {
		t.Fatal("Expected second lock to fail, but it succeeded")
	}
}

// TestSingleInstance_Unlock checks lock file release
func TestSingleInstance_Unlock(t *testing.T) {
	// Сначала захватываем lock обычным путем
	si := utils.NewSingleInstance("test-unlock-app")
	locked, err := si.TryLock()

	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	if !locked {
		t.Fatal("Expected to acquire lock, but failed")
	}

	lockPath := si.GetLockPath()

	// Unlock должен закрыть файл и удалить его с диска
	err = si.Unlock()
	if err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}

	if _, err := os.Stat(lockPath); !os.IsNotExist(err) {
		t.Fatalf("Lock file still exists at %s after unlock", lockPath)
	}

	// После Unlock новый экземпляр снова должен получить lock
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

// TestSingleInstance_StaleLock checks stale lock recovery
func TestSingleInstance_StaleLock(t *testing.T) {
	// Подкладываем битый lock, как будто прошлый процесс уже умер
	lockPath := filepath.Join(os.TempDir(), "test-stale-app.lock")
	_ = os.Remove(lockPath)

	t.Cleanup(func() {
		_ = os.Remove(lockPath)
	})

	// PID специально нереальный, чтобы lock считался протухшим
	err := os.WriteFile(lockPath, []byte("999999999"), 0600)
	if err != nil {
		t.Fatalf("Failed to create fake lock file: %v", err)
	}

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

// TestSingleInstance_LockPath checks generated lock path
func TestSingleInstance_LockPath(t *testing.T) {
	si := utils.NewSingleInstance("test-path-app")
	lockPath := si.GetLockPath()

	// Путь должен указывать в temp и сохранять имя приложения
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
