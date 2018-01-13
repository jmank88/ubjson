package ubjson

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// Based on 'encoding/json/encode.go'.
var fieldCache struct {
	value atomic.Value // map[reflect.Type]fields
	mu    sync.Mutex   // used only by writers
}

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
// Based on 'encoding/json/encode.go'.
func cachedTypeFields(t reflect.Type) fields {
	m, _ := fieldCache.value.Load().(map[reflect.Type]fields)
	f, ok := m[t]
	if ok {
		return f
	}

	// Compute names without lock.
	// Might duplicate effort but won't hold other computations back.
	f = typeFields(t)

	fieldCache.mu.Lock()
	m, _ = fieldCache.value.Load().(map[reflect.Type]fields)
	newM := make(map[reflect.Type]fields, len(m)+1)
	for k, v := range m {
		newM[k] = v
	}
	newM[t] = f
	fieldCache.value.Store(newM)
	fieldCache.mu.Unlock()
	return f
}

// Indexes fields by 'ubjson' struct tag if present, otherwise name.
func typeFields(t reflect.Type) fields {
	fs := fields{
		indexByName: make(map[string]int),
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath == "" {
			name := f.Name
			// Check for 'ubjson' struct tag.
			if v, ok := f.Tag.Lookup("ubjson"); ok {
				name = v
			}
			fs.names = append(fs.names, name)
			fs.indexByName[name] = i
		}
	}
	return fs
}

type fields struct {
	names       []string
	indexByName map[string]int
}
