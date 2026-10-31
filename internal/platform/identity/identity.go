package identity

import "crypto/rand"

func New() string { return rand.Text() }
