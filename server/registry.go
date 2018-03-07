package server

import (
	"encoding/gob"
	"fmt"
	"mgrok/cache"
	"mgrok/log"
	"net"
	"sync"
	"time"
)

const (
	cacheSaveInterval time.Duration = 10 * time.Minute
)

type cacheURL string

func (url cacheURL) Size() int {
	return len(url)
}

// TunnelRegistry maps a tunnel URL to Tunnel structures
type TunnelRegistry struct {
	tunnels  map[string]*Tunnel
	affinity *cache.LRUCache
	log.Logger
	sync.RWMutex
	centerRegistry *TunnelCenterRegistry
}

func newTunnelRegistry(cacheSize uint64, cacheFile string, redirectData *TunnelCenterRegistry) *TunnelRegistry {
	registry := &TunnelRegistry{
		tunnels:        make(map[string]*Tunnel),
		affinity:       cache.NewLRUCache(cacheSize),
		Logger:         log.NewPrefixLogger("registry", "tun"),
		centerRegistry: redirectData,
	}

	// LRUCache uses Gob encoding. Unfortunately, Gob is fickle and will fail
	// to encode or decode any non-primitive types that haven't been "registered"
	// with it. Since we store cacheUrl objects, we need to register them here first
	// for the encoding/decoding to work
	var urlobj cacheURL
	gob.Register(urlobj)

	// try to load and then periodically save the affinity cache to file, if specified
	if cacheFile != "" {
		err := registry.affinity.LoadItemsFromFile(cacheFile)
		if err != nil {
			registry.Error("Failed to load affinity cache %s: %v", cacheFile, err)
		}

		registry.saveCacheThread(cacheFile, cacheSaveInterval)
	} else {
		registry.Info("No affinity cache specified")
	}

	return registry
}

// Spawns a goroutine the periodically saves the cache to a file.
func (r *TunnelRegistry) saveCacheThread(path string, interval time.Duration) {
	go func() {
		r.Info("Saving affinity cache to %s every %s", path, interval.String())
		for {
			time.Sleep(interval)

			r.Debug("Saving affinity cache")
			err := r.affinity.SaveItemsToFile(path)
			if err != nil {
				r.Error("Failed to save affinity cache: %v", err)
			} else {
				r.Info("Saved affinity cache")
			}
		}
	}()
}

// register a tunnel with a specific url, returns an error
// if a tunnel is already registered at that url
func (r *TunnelRegistry) register(url string, t *Tunnel) error {
	r.Lock()
	defer r.Unlock()

	if r.tunnels[url] != nil {
		return fmt.Errorf("The tunnel %s is already registered.\r", url)
	}

	r.tunnels[url] = t
	if r.centerRegistry != nil {
		r.centerRegistry.register(url, t)
	}
	return nil
}

func (r *TunnelRegistry) cacheKeys(t *Tunnel) (ip string, id string) {
	clientIP := t.ctl.conn.RemoteAddr().(*net.TCPAddr).IP.String()
	clientID := t.ctl.id

	ipKey := fmt.Sprintf("client-ip-%s:%s", t.req.Protocol, clientIP)
	idKey := fmt.Sprintf("client-id-%s:%s", t.req.Protocol, clientID)
	return ipKey, idKey
}

func (r *TunnelRegistry) getCachedRegistration(t *Tunnel) (url string) {
	ipCacheKey, idCacheKey := r.cacheKeys(t)

	// check cache for ID first, because we prefer that over IP which might
	// not be specific to a user because of NATs
	if v, ok := r.affinity.Get(idCacheKey); ok {
		url = string(v.(cacheURL))
		t.Debug("Found registry affinity %s for %s", url, idCacheKey)
	} else if v, ok := r.affinity.Get(ipCacheKey); ok {
		url = string(v.(cacheURL))
		t.Debug("Found registry affinity %s for %s", url, ipCacheKey)
	}
	return
}

func (r *TunnelRegistry) registerAndCache(url string, t *Tunnel) (err error) {
	if err = r.register(url, t); err == nil {
		// we successfully assigned a url, cache it
		ipCacheKey, idCacheKey := r.cacheKeys(t)
		r.affinity.Set(ipCacheKey, cacheURL(url))
		r.affinity.Set(idCacheKey, cacheURL(url))
	}
	return

}

// Register a tunnel with the following process:
// Consult the affinity cache to try to assign a previously used tunnel url if possible
// Generate new urls repeatedly with the urlFn and register until one is available.
func (r *TunnelRegistry) registerRepeat(urlFn func() string, t *Tunnel) (string, error) {
	url := r.getCachedRegistration(t)
	if url == "" {
		url = urlFn()
	}

	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		if err := r.registerAndCache(url, t); err != nil {
			// pick a new url and try again
			url = urlFn()
		} else {
			// we successfully assigned a url, we're done
			return url, nil
		}
	}

	return "", fmt.Errorf("Failed to assign a URL after %d attempts!\r", maxAttempts)
}

func (r *TunnelRegistry) del(url string) {
	r.Lock()
	defer r.Unlock()
	r.Debug("Delete tunnel named %s.", url)
	delete(r.tunnels, url)
	r.centerRegistry.del(url)
}

func (r *TunnelRegistry) get(url string) *Tunnel {
	r.RLock()
	defer r.RUnlock()
	return r.tunnels[url]
}

// ControlRegistry maps a client ID to Control structures
type ControlRegistry struct {
	controls map[string]*Control
	log.Logger
	sync.RWMutex
}

func newControlRegistry() *ControlRegistry {
	return &ControlRegistry{
		controls: make(map[string]*Control),
		Logger:   log.NewPrefixLogger("registry", "ctl"),
	}
}

func (r *ControlRegistry) get(clientID string) *Control {
	r.RLock()
	defer r.RUnlock()
	return r.controls[clientID]
}

func (r *ControlRegistry) add(clientID string, ctl *Control) (oldCtl *Control) {
	r.Lock()
	defer r.Unlock()

	oldCtl = r.controls[clientID]
	if oldCtl != nil {
		oldCtl.Replaced(ctl)
	}

	r.controls[clientID] = ctl
	r.Info("Registered control with id %s", clientID)
	return
}

func (r *ControlRegistry) del(clientID string) error {
	r.Lock()
	defer r.Unlock()
	if r.controls[clientID] == nil {
		return fmt.Errorf("No control found for client id: %s", clientID)
	}

	r.Info("Removed control registry id %s", clientID)
	delete(r.controls, clientID)
	return nil

}
