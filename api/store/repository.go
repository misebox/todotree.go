package store

import "todotree/clock"

type Repository struct {
	Clocker clock.Clocker
}
