package main

import (
	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

func makeKSpaceData() [][]float32 {
	nX, nY := 256, 128
	square := make([][]float64, nY)
	for y := 0; y < nY; y++ {
		square[y] = make([]float64, nX)
		for x := 0; x < nX; x++ {
			if (x > nX/4) && (x < 3*nX/4) && (y > nY/8) && (y < 7*nY/8) {
				square[y][x] = 1.0
			} else {
				square[y][x] = 0.0
			}
		}
	}

	/* saveFloat64Image(square, "square.png") */

	squarec := dsputils.ToComplex2(square)

	/* shift := fftshift2(squarec) */
	/* saveCmplx128Image(shift, "shifted.png") */

	fft_square := fftshift2(fft.FFT2(fftshift2(squarec)))

	data := make([][]float32, nY)
	for a := 0; a < nY; a++ {
		data[a] = make([]float32, 2*nX)
		for b := 0; b < nX; b++ {
			data[a][b*2] = float32(real(fft_square[a][b]))
			data[a][b*2+1] = float32(imag(fft_square[a][b]))
		}
	}

	return data
}

func fftshift2(data [][]complex128) [][]complex128 {
	w, h := len(data[0]), len(data)
	hw, hh := w/2, h/2

	// copy data so we don't modify it
	ret := make([][]complex128, len(data))
	for y := 0; y < h; y++ {
		ret[y] = make([]complex128, w)
		copy(ret[y], data[y])
	}

	// swap top and bottom halves
	for y := 0; y < hh; y++ {
		for x := 0; x < w; x++ {
			tmp := ret[y][x]
			ret[y][x] = ret[hh+y][x]
			ret[hh+y][x] = tmp
		}
	}

	// swap left and right halves
	for y := 0; y < h; y++ {
		for x := 0; x < hw; x++ {
			tmp := ret[y][x]
			ret[y][x] = ret[y][hw+x]
			ret[y][hw+x] = tmp
		}
	}

	return ret
}

func saveFloat64Image(data [][]float64, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	r := image.Rect(0, 0, len(data[0]), len(data))
	img := image.NewGray(r)
	for y := 0; y < len(data); y++ {
		for x := 0; x < len(data[y]); x++ {
			img.Set(x, y, color.Gray{uint8(data[y][x] * 255)})
		}
	}

	return png.Encode(f, img)
}

func saveCmplx128Image(data [][]complex128, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	r := image.Rect(0, 0, len(data[0]), len(data))
	img := image.NewGray(r)
	for y := 0; y < len(data); y++ {
		for x := 0; x < len(data[y]); x++ {
			img.Set(x, y, color.Gray{uint8(cmplx.Abs(data[y][x]) * 255)})
		}
	}

	return png.Encode(f, img)
}
