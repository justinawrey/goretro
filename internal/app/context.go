package app

import "context"

type ContextHolder struct {
	ctx context.Context
}

func (c *ContextHolder) SetContext(ctx context.Context) {
	c.ctx = ctx
}
