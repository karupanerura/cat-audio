package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/go-mp3"
	"github.com/k0kubun/pp"
	"github.com/karupanerura/riffbin"
	"github.com/karupanerura/wavebin"
)

func main() {
	var (
		output string
	)
	flag.StringVar(&output, "out", "concat.wav", "output file path")
	flag.Parse()

	sources := flag.Args()
	if output == "" || len(sources) == 0 {
		log.Fatalf("Usage: cat-audio -o output.wav part1.wav part2.wav ...")
	}

	run(output, sources)
}

func run(output string, sources []string) {
	var readers []io.Reader
	for _, filePath := range sources {
		source, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("%s: %v", filePath, err)
		}
		defer source.Close()

		switch filepath.Ext(filePath) {
		case ".wav":
			chunk, err := riffbin.ReadSections(source)
			if err != nil {
				log.Fatalf("%s: riffbin.ReadSections: %v", filePath, err)
			}

			format, _, _, samples, err := wavebin.ParseWaveRIFF(chunk, true)
			if err != nil {
				log.Fatalf("%s: wavebin.ParseWaveRIFF: %v", filePath, err)
			}
			if !(format.Channels() == 2 && format.SamplesPerSecond() == 44100 && format.SignificantBitsPerSample() == 16) {
				log.Fatalf("%s: unsupported format %+v", filePath, pp.Sprint(format))
			}
			readers = append(readers, samples)

		case ".mp3":
			decoder, err := mp3.NewDecoder(source)
			if err != nil {
				log.Fatalf("%s: mp3.NewDecoder: %v", filePath, err)
			}

			readers = append(readers, decoder)

		default:
			log.Fatalf("%s: unsupported format", filePath)
		}
	}

	var postProcess func(string)
	switch filepath.Ext(output) {
	case ".wav":
		postProcess = func(string) {} // nop

	case ".mp3":
		output = strings.TrimSuffix(output, ".mp3") + ".wav"
		postProcess = func(output string) {
			log.Printf("%s: mp3 encode", output)
			cmd := exec.Command("lame", "-m", "j", "-h", "-v", "-b", "192", "-V2", output)
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Fatalf("%s: lame: %+v", output, err)
			}

			os.Remove(output)
			log.Printf("%s: created", strings.TrimSuffix(output, ".wav")+".mp3")
		}

	default:
		log.Fatalf("%s: unsupported format", output)
	}
	defer postProcess(output)

	f, err := os.Create(output)
	if err != nil {
		log.Fatalf("%s: %v", output, err)
	}
	defer f.Close()

	w, err := wavebin.CreateSampleWriter(f, &wavebin.ExtendedFormatChunk{
		MetaFormat: wavebin.NewPCMMetaFormat(wavebin.StereoChannels, 44100, 16),
	})
	if err != nil {
		log.Fatalf("%s: %v", output, err)
	}
	defer w.Close()

	var buf [8192]byte
	for i, r := range readers {
		log.Printf("read from %s", sources[i])
		_, err = io.CopyBuffer(w, r, buf[:])
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("%s: created", output)
}
