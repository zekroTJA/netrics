package watcher

func (w *Watcher) RTTHandler() float64 {
	return float64(w.LastValues.Ping.Microseconds()) / 1000
}

func (w *Watcher) DLSpeedHandler() float64 {
	return w.LastValues.DLSpeed
}

func (w *Watcher) ULSpeedHandler() float64 {
	return w.LastValues.ULSpeed
}

func (w *Watcher) ErrorCountHandler() float64 {
	return float64(w.ErrorCount)
}
