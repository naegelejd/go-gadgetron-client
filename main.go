package main

import (
    /* "fmt" */
    "log"
    "flag"
    "github.com/mjibson/go-dsp/fft"
    "github.com/mjibson/go-dsp/dsputils"
    "os"
    "image"
    "image/color"
    "image/png"
)

func fftshift2(data [][]complex128) [][]complex128 {
    // copy data
    tmp := make([][]complex128, len(data))
    for y := 0; y < len(data); y++ {
        tmp[y] = make([]complex128, len(data[y]))
        copy(tmp[y], data[y])
    }

    for y := 0; y < len(data) / 2 - 1; y++ {
        for x := 0; x < len(data[y]); x++ {
            data[y][x] = tmp[2*y][x]
            data[2*y][x] = tmp[y][x]
        }
    }

    for y := 0; y < len(data) / 2 - 1; y++ {
        for x := 0; x < len(data[y]) / 2; x++ {
            data[y][x] = data[y][2*x]
            data[y][2*x] = tmp[y][x]
        }
    }

    return data
}

func saveTmpPNG(data [][]float32) {
    f, _ := os.Create("tmp.png")
    defer f.Close()

    r := image.Rect(0, 0, len(data[0]), len(data))
    img := image.NewGray(r)
    for y := 0; y < len(data); y++ {
        for x := 0; x < len(data[y]); x++ {
            img.Set(x, y, color.Gray{uint8(data[y][x] * 255)})
        }
    }

    png.Encode(f, img)
}

func makeData() [][]float32 {
    nX, nY := 256, 128
    square := make([][]float64, nY)
    for y := 0; y < nY; y++ {
        square[y] = make([]float64, nX)
        for x := 0; x < nX; x++ {
            if (x > nX / 4) && (x < 3 * nX / 4) && (y > nY / 8) && (y < 7 * nY / 8) {
                square[y][x] = 1.0
            } else {
                square[y][x] = 0.0
            }
        }
    }

    squarec := dsputils.ToComplex2(square)
    fft_square := fftshift2(fft.FFT2(fftshift2(squarec)))
    /* fft_square := fft.FFT2(squarec) */
    /* fft_square := fft.FFT2Real(square) */

    data := make([][]float32, len(fft_square))
    for a := 0; a < len(fft_square); a++ {
        data[a] = make([]float32, 2 * len(fft_square[a]))
        for b := 0; b < len(fft_square[a]); b++ {
            data[a][b * 2] = float32(real(fft_square[a][b]))
            data[a][b * 2 + 1] = float32(imag(fft_square[a][b]))
        }
    }

    return data
}

func main () {
    var host string
    var port int
    flag.StringVar(&host, "host", "", "hostname")
    flag.IntVar(&port, "port", 9002, "port number")
    flag.Parse()

    conn := newGadgetronConnector(host, port)
    defer conn.Close()

    conn.registerReader(GADGET_MESSAGE_ISMRMRD_IMAGE_REAL_FLOAT, &IsmrmrdImagePNGReader{})

    conn.sendGadgetronConfigurationFile("default.xml")
    conn.sendGadgetronParameters(xml_config)

    head := AcquisitionHeader {}
    head.Version = 1
    head.Number_of_samples = 256
    head.Available_channels = 1
    head.Active_channels = 1
    head.Center_sample = 128
    head.Sample_time_us = 5.0
    //log.Printf("%+v\n", head)

    // 5 slices (corresponding to slice limits in xml_config)
    for s := 0; s < 5; s++ {
        for i,line := range data {
            if (i == 0) {
                head.Flags = ACQ_FIRST_IN_SLICE
            } else if (i == 127) {
                head.Flags = ACQ_LAST_IN_SLICE
            }

            head.Idx.Kspace_encode_step_1 = uint16(i)
            head.Idx.Slice = uint16(s)
            acq := &Acquisition {head, traj, line}
            conn.sendIsmrmrdAcquisition(acq)
        }
    }

    conn.sendCloseMessage()

    conn.readImages()

    log.Println("Finished!")
}
