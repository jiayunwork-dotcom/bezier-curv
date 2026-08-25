package arclen

import "context"

func abortArcContext() error {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
