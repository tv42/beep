package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/oto"
)

//go:generate go build github.com/tv42/becky
//go:generate ./becky snooker.pcm_s16le

func pcm_s16le(a asset) io.WriterTo {
	return strings.NewReader(a.Content)
}

func run() error {
	const (
		// If you change the sample file, you'll probably need to change these.
		// These assume pcm_s16le, like you'd get from this command:
		//
		//     ffmpeg -i someaudiofile -f s16le -acodec pcm_s16le foo.pcm_s16le

		numChannels = 2
		sampleRate  = 48000
		depthBytes  = 2
	)

	// Let it buffer the whole sample if it wants to. If this size
	// goes below whole/2, we start losing parts of the sample. The
	// library wasn't very helpful in debugging this, it seems audio
	// is a world where "sleep roughly long enough" is used instead of
	// proper completion guarantees.
	const bufSize = numChannels * depthBytes * sampleRate

	c, err := oto.NewContext(sampleRate, numChannels, depthBytes, bufSize)
	if err != nil {
		return err
	}
	defer c.Close()
	p := c.NewPlayer()
	defer p.Close()
	if _, err := snooker.WriteTo(p); err != nil {
		return err
	}

	if err := p.Close(); err != nil {
		return err
	}

	return nil
}

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s\n", prog)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s takes no arguments.\n", prog)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
