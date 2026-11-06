package runtime

import (
	"fmt"
	"sync"

	"github.com/nguyenngocanh94/mate/internal/observability"
)

// MemoryNameRegistry is an in-process LiveNameRegistry: it lives and dies
// with one mate invocation, so it protects nothing across processes or
// workspaces. The fake runtime uses it. A durable backing needs what a map
// in one process cannot give, because every mate invocation is its own
// process: a runtime_binding unique constraint for agent names, a
// cross-workspace owner record for session names, and Herdr's live list for
// staleness.
type MemoryNameRegistry struct {
	mu    sync.Mutex
	names map[string]string // session\x00name -> rawID
}

// NewMemoryNameRegistry returns an empty registry.
func NewMemoryNameRegistry() *MemoryNameRegistry {
	return &MemoryNameRegistry{names: make(map[string]string)}
}

func registryKey(session, name string) string {
	return session + "\x00" + name
}

// Occupied implements LiveNameRegistry.
func (r *MemoryNameRegistry) Occupied(session, name string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	owner, ok := r.names[registryKey(session, name)]
	return owner, ok
}

// Reserve implements LiveNameRegistry.
func (r *MemoryNameRegistry) Reserve(session, name, rawID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := registryKey(session, name)
	if owner, ok := r.names[key]; ok && owner != rawID {
		return observability.WrapError(
			observability.CodeAlreadyExists,
			fmt.Sprintf("agent name %q is taken", name),
			ErrNameCollision,
		)
	}
	r.names[key] = rawID
	return nil
}

// Release implements LiveNameRegistry.
func (r *MemoryNameRegistry) Release(session, name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.names, registryKey(session, name))
}
