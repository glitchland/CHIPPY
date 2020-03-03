package main

import (
	spk "pkg/speaker"
)

func main() {
	s := new(spk.Speaker)
	s.Init()
	s.PlayTone(440, 1) // frequency, duration (in seconds)
}
