package server

import "context"

type Server interface {
	Start(ctx context.Context) error
}
