package pluginruntime

import (
	"log"
	"net/http"
	"sync"

	"github.com/snassr/blog-348OEjOG-goplugins/external/gen/plugin-proto-go/plugin/v1/pluginv1connect"
	"github.com/snassr/blog-348OEjOG-goplugins/external/plugin/v1/plugin"
	"github.com/snassr/blog-348OEjOG-goplugins/internal/pluginruntime/plugins/plugin_en"
)

type Manager struct {
	mu      sync.RWMutex
	plugins map[string]plugin.Plugin
}

func NewManager() *Manager {
	m := &Manager{plugins: make(map[string]plugin.Plugin)}

	m.Add("plugin_en", plugin_en.New())

	return m
}

func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]string, 0, len(m.plugins))

	for id := range m.plugins {
		out = append(out, id)
	}

	return out
}

func (m *Manager) Get(id string) (plugin.Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.plugins[id]

	return p, ok
}

func (m *Manager) Add(id string, p plugin.Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.plugins[id] = p
	log.Printf("[manager] added plugin %q", id)
}

func (m *Manager) RegisterRemote(id, address string) error {
	addr := "http://" + address
	client := pluginv1connect.NewPluginServiceClient(http.DefaultClient, addr)

	m.Add(id, &PluginAdapter{client: client})

	return nil
}
