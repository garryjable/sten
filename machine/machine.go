package machine

type Machine interface {
	StartCapture() error
	StopCapture()
	SetCallback(cb StrokeCallback)
}
