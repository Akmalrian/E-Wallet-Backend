package repositories

import "sync"

// TokenBlacklist — menyimpan token yang sudah logout
type TokenBlacklist struct {
	blacklist sync.Map
}

func NewTokenBlacklist() *TokenBlacklist {
	return &TokenBlacklist{}
}

// Add — tambahkan token ke blacklist saat logout
func (t *TokenBlacklist) Add(token string) {
	t.blacklist.Store(token, true)
}

func (t *TokenBlacklist) IsBlacklisted(token string) bool {
	_, exists := t.blacklist.Load(token)
	return exists
}
