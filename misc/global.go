package misc

import "sync"

// Mut mutable shared data
var Mut = &sync.RWMutex{}

// Mutex global
var Mutex = &sync.Mutex{}
