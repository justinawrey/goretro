package main

import "context"

type WebviewAudioDriver struct {
	ctx context.Context
}

func newWebviewAudioDriver() *WebviewAudioDriver {
	return &WebviewAudioDriver{}
}
