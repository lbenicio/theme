package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func TestParseFrontMatter(t *testing.T) {
	got := parseFrontMatter("---\ntitle: \"A title: with punctuation\"\nauthor: 'Ada Lovelace'\n# comment\ninvalid line\n---\nbody")
	want := map[string]string{
		"title":  "A title: with punctuation",
		"author": "Ada Lovelace",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseFrontMatter() = %#v, want %#v", got, want)
	}

	if got := parseFrontMatter("body without front matter"); len(got) != 0 {
		t.Fatalf("front matter missing: got %#v, want empty map", got)
	}
}

func TestHumanizeSlug(t *testing.T) {
	if got, want := humanizeSlug("building-a-static-site"), "Building A Static Site"; got != want {
		t.Fatalf("humanizeSlug() = %q, want %q", got, want)
	}
}

func TestWrapLinesPreservesWordsAndFitsAvailableWidth(t *testing.T) {
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		t.Fatal(err)
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{Size: 16, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		t.Fatal(err)
	}
	defer face.Close()

	text := "A small title that wraps safely"
	lines := wrapLines(face, text, 120)
	if len(lines) < 2 {
		t.Fatalf("wrapLines() returned %v, want multiple lines", lines)
	}
	if got := strings.Join(lines, " "); got != text {
		t.Fatalf("wrapped lines lost or reordered words: %q", got)
	}
	for _, line := range lines {
		if width := textWidth(face, line); width > 120 {
			t.Errorf("line %q is %d pixels wide, want at most 120", line, width)
		}
	}
}

func TestLerpInterpolatesRGBAChannels(t *testing.T) {
	start := color.RGBA{R: 0, G: 20, B: 100, A: 255}
	end := color.RGBA{R: 100, G: 60, B: 0, A: 255}
	want := color.RGBA{R: 50, G: 40, B: 50, A: 255}
	if got := lerp(start, end, 0.5); got != want {
		t.Fatalf("lerp() = %#v, want %#v", got, want)
	}
}

func TestSavePNGWritesDecodableImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	path := filepath.Join(t.TempDir(), "preview.png")
	if err := savePNG(img, path); err != nil {
		t.Fatalf("savePNG() error = %v", err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoded, err := png.Decode(file)
	if err != nil {
		t.Fatalf("written PNG cannot be decoded: %v", err)
	}
	if got, want := decoded.Bounds(), image.Rect(0, 0, 2, 2); got != want {
		t.Fatalf("decoded bounds = %v, want %v", got, want)
	}
}
