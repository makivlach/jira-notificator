package audio

import (
	"bytes"
	"github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/mp3"
)

type reader struct {
	*bytes.Reader
}

func (*reader) Close() error {
	return nil
}

type Player struct {
	audio audio.Audio
}

func NewPlayer(data []byte) (*Player, error) {
	reader := &reader{bytes.NewReader(data)}
	a, err := mp3.Load(reader)
	if err != nil {
		return nil, err
	}

	return &Player{
		a,
	}, nil
}

func (p Player) Play() error {
	return <-p.audio.Play()
}
