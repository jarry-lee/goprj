package main

import (
	"math"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

const (
	sampleRate   = 44100     // 샘플링 속도 (Hz)
	bitDuration  = 0.1       // 각 비트의 지속 시간 (초)
	freq0        = 1000.0    // 0에 해당하는 주파수 (Hz)
	freq1        = 2000.0    // 1에 해당하는 주파수 (Hz)
	amplitude    = 0.5       // 파형의 최대 진폭 (0.0 ~ 1.0)
)

func encodeFSK(data []byte, filename string) error {
	numSamples := int(bitDuration * sampleRate)
	samples := make([]int, 0)

	for _, byteData := range data {
		for i := 0; i < 8; i++ { // 각 바이트를 비트로 변환
			bit := (byteData >> (7 - i)) & 1
			freq := freq0
			if bit == 1 {
				freq = freq1
			}

			for j := 0; j < numSamples; j++ {
				t := float64(j) / sampleRate
				sampleValue := int(amplitude * math.Sin(2*math.Pi*freq*t) * math.MaxInt16)
				samples = append(samples, sampleValue)
			}
		}
	}

	buf := &audio.IntBuffer{
		Data:           samples,
		Format:         &audio.Format{SampleRate: sampleRate, NumChannels: 1},
		SourceBitDepth: 16,
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := wav.NewEncoder(file, sampleRate, 16, 1, 1)
	err = encoder.Write(buf)
	if err != nil {
		return err
	}
	return encoder.Close()
}

func main() {
	data := []byte("Hello, FSK!") // 인코딩할 데이터
	err := encodeFSK(data, "modem_output.wav")
	if err != nil {
		panic(err)
	}
	println("Encoding completed: modem_output.wav")
}
