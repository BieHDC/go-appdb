package cache

import (
	"sync"
	"sync/atomic"
)

type TemplateBytes map[string][]byte
type TemplateCache struct {
	syncer     atomic.Value
	writeMutex sync.Mutex // only used by writers
}

// This thing is made for many reads, little writes.
func NewTemplateCache() *TemplateCache {
	tcache := TemplateCache{}
	tcache.syncer.Store(make(TemplateBytes))
	return &tcache
}

func (tcache *TemplateCache) GetEntry(tk string) ([]byte, bool) {
	copyofsyncer := tcache.syncer.Load().(TemplateBytes)
	value, exists := copyofsyncer[tk]
	return value, exists
}

func (tcache *TemplateCache) SetEntry(tk string, tv []byte) {
	tcache.writeMutex.Lock()
	defer tcache.writeMutex.Unlock()

	copyofsyncer := tcache.syncer.Load().(TemplateBytes)
	tempcopy := make(TemplateBytes)
	for k, v := range copyofsyncer {
		tempcopy[k] = v
	}
	tempcopy[tk] = tv
	tcache.syncer.Store(tempcopy)
}
