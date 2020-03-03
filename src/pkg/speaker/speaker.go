package speaker

import (
	"log"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/generator"
	"github.com/go-audio/transforms"
	"github.com/gordonklaus/portaudio"
)

type Speaker struct {
	playTime uint16
	freq     float32
}

func (s *Speaker) Init() {
	s.playTime = 10
	s.freq = 440.0
}

func (s *Speaker) PlayTone(frequency float64, duration uint) {

	bufferSize := 512
	buf := &audio.FloatBuffer{
		Data:   make([]float64, bufferSize),
		Format: audio.FormatMono44100,
	}
	currentNote := frequency
	osc := generator.NewOsc(generator.WaveSine, currentNote, buf.Format.SampleRate)
	osc.Amplitude = 1

	currentVol := osc.Amplitude

	// Audio output
	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]float32, bufferSize)
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(out), &out)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()

	timeout := time.After(time.Duration(duration) * time.Second)
	tick := time.Tick(1 * time.Millisecond)
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			return
		// Got a tick, we should check on doSomething()
		case <-tick:
			// populate the out buffer
			if err := osc.Fill(buf); err != nil {
				log.Printf("error filling up the buffer")
			}

			transforms.Gain(buf, currentVol)
			s.f64ToF32Copy(out, buf.Data)

			if err := stream.Write(); err != nil {
				log.Printf("error writing to stream : %v\n", err)
			}
		}
	}
}

func (s *Speaker) f64ToF32Copy(dst []float32, src []float64) {
	for i := range src {
		dst[i] = float32(src[i])
	}
}
