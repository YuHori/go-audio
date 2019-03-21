package main

import (
    "encoding/binary"
    "fmt"
    "math"
    "time"

    "golang.org/x/mobile/exp/audio/al"
    "golang.org/x/mobile/exp/f32"
)

const (
    Pi         = float32(math.Pi)
    Fmt        = al.FormatStereo16
    SampleRate = 64
)

type Oscillator func() float32

func GenOscillator(freq float32) Oscillator {
    dt := 1.0 / float32(SampleRate)
    k := 2.0 * Pi * freq
    T := 1.0 / freq
    t := float32(0.0)
    return func() float32 {
        res := f32.Sin(k * t)
        t += dt
        if t > T {
            t -= T
        }
        return res
    }
}

func Multiplex(fs ...Oscillator) Oscillator {
    return func() float32 {
        res := float32(0)
        for _, osc := range fs {
            res += osc()
        }
        return res
    }
}

type Piano struct {
    notes      []bool
    oscillator Oscillator
}

func NewPiano(freqs []float32) *Piano {
    p := new(Piano)
    p.notes = make([]bool, len(freqs))
    osc := []Oscillator{}
    for _, f := range freqs {
        osc = append(osc, GenOscillator(f))
    }
    p.oscillator = Multiplex(osc...)
    return p
}

func main() {
    fmt.Println("----start----")
    pianoPlayer := NewPiano([]float32{
        246.941650628,
        261.625565301,
        277.182630977,
        293.664767917,
        311.126983722,
        329.627556913,
        349.228231433,
        369.994422712,
        391.995435982,
        415.30469758,
        440.0,
        466.163761518,
        493.883301256,
        523.251130601,
    })

    al.OpenDevice()
    s := al.GenSources(1)
    b := al.GenBuffers(1)
    buf := make([]byte, 2048)
    for n := 0; n < 2048; n += 2 {
        f := pianoPlayer.oscillator()
        v := int16(float32(92767) * f)
        binary.LittleEndian.PutUint16(buf[n:n+2], uint16(v))
    }
    fmt.Println(len(buf))
    b[0].BufferData(Fmt, buf, SampleRate)
    s[0].QueueBuffers(b...)

    al.PlaySources(s[0])

    time.Sleep(1000 * time.Millisecond)
}
