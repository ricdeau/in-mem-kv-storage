package tests

import (
	"github.com/ricdeau/in-mem-kv-storage/storage"
	"testing"
)

func TestStorageCRUD(t *testing.T) {
	data1 := "some-data"
	data2 := "some-other-data"
	key := "some-key"
	stor := storage.New()

	actual, exists := stor.Get(key)
	if exists {
		t.Error("Initial data exists")
	}

	newVal := stor.Set(key, &storage.Data{Type: "string", Payload: []byte(data1)})
	if !newVal {
		t.Error("Created value wasn't marked as new")
	}

	actual, exists = stor.Get(key)
	if !exists {
		t.Error("Created value doesn't exist")
	}
	if string(actual.Payload) != data1 {
		t.Errorf("Retrieved and initial values does't match")
	}

	newVal = stor.Set(key, &storage.Data{Type: "string", Payload: []byte(data2)})
	if newVal {
		t.Error("Updated value marked as new")
	}

	actual, exists = stor.Get(key)
	if !exists {
		t.Error("Updated value doesn't exist")
	}
	if string(actual.Payload) != data2 {
		t.Errorf("Retrieved and initial values does't match")
	}

	stor.Delete(key)
	_, exists = stor.Get(key)
	if exists {
		t.Error("Value hasn't been deleted")
	}

	stor.Delete("  ")
}
