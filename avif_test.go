package avif

import (
	"bytes"
	_ "embed"
	"image"
	"image/jpeg"
	"io"
	"os"
	"testing"
)

//go:embed testdata/test8.avif
var testAvif8 []byte

//go:embed testdata/test10.avif
var testAvif10 []byte

//go:embed testdata/test.avifs
var testAvifAnim []byte

func TestLoadLibrary(t *testing.T) {
	_, err := loadLibrary()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecodeDynamic(t *testing.T) {
	img, _, err := decodeDynamic(bytes.NewReader(testAvif8), false, false)
	if err != nil {
		t.Fatal(err)
	}

	w, err := writeCloser()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = jpeg.Encode(w, img.Image[0], nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDecode10Dynamic(t *testing.T) {
	img, _, err := decodeDynamic(bytes.NewReader(testAvif10), false, false)
	if err != nil {
		t.Fatal(err)
	}

	w, err := writeCloser()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = jpeg.Encode(w, img.Image[0], nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDecodeAnimDynamic(t *testing.T) {
	ret, _, err := decodeDynamic(bytes.NewReader(testAvifAnim), false, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(ret.Image) != len(ret.Delay) {
		t.Errorf("got %d, want %d", len(ret.Delay), len(ret.Image))
	}

	if len(ret.Image) != 17 {
		t.Errorf("got %d, want %d", len(ret.Image), 17)
	}

	for _, img := range ret.Image {
		err = jpeg.Encode(io.Discard, img, nil)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestImageDecode(t *testing.T) {
	img, _, err := image.Decode(bytes.NewReader(testAvif8))
	if err != nil {
		t.Fatal(err)
	}

	w, err := writeCloser()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = jpeg.Encode(w, img, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestImageDecodeAnim(t *testing.T) {
	img, _, err := image.Decode(bytes.NewReader(testAvifAnim))
	if err != nil {
		t.Fatal(err)
	}

	w, err := writeCloser()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = jpeg.Encode(w, img, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDecodeConfigDynamic(t *testing.T) {
	_, cfg, err := decodeDynamic(bytes.NewReader(testAvif8), true, false)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Width != 512 {
		t.Errorf("width: got %d, want %d", cfg.Width, 512)
	}

	if cfg.Height != 512 {
		t.Errorf("height: got %d, want %d", cfg.Height, 512)
	}
}

func TestEncodeDynamic(t *testing.T) {
	img, err := Decode(bytes.NewReader(testAvif8))
	if err != nil {
		t.Fatal(err)
	}

	w, err := writeCloser()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = encodeDynamic(w, img, DefaultQuality, DefaultQuality, DefaultSpeed, image.YCbCrSubsampleRatio420)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkDecodeDynamic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := decodeDynamic(bytes.NewReader(testAvif8), false, false)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkDecodeConfigDynamic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := decodeDynamic(bytes.NewReader(testAvif8), true, false)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkEncodeDynamic(b *testing.B) {
	img, err := Decode(bytes.NewReader(testAvif8))
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		err := encodeDynamic(io.Discard, img, DefaultQuality, DefaultQuality, DefaultSpeed, image.YCbCrSubsampleRatio420)
		if err != nil {
			b.Error(err)
		}
	}
}

type discard struct{}

func (d discard) Close() error {
	return nil
}

func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}

var discardCloser io.WriteCloser = discard{}

func writeCloser(s ...string) (io.WriteCloser, error) {
	if len(s) > 0 {
		f, err := os.Create(s[0])
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	return discardCloser, nil
}
