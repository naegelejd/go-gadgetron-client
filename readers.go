package main

import (
    "os"
    "fmt"
    "net"
    "image"
    "image/png"
    "image/color"
    "encoding/binary"
)

type GadgetMessageReader interface {
    Read (net.Conn) error
}

type IsmrmrdImagePNGReader struct {
    count int
}

func (r *IsmrmrdImagePNGReader) Read(conn net.Conn) error {
    var head ImageHeader
    err := binary.Read(conn, binary.LittleEndian, &head)
    if err != nil { return err }
    //fmt.Printf("%+v\n", head)

    nelem := head.Matrix_size[0] * head.Matrix_size[1] * head.Matrix_size[2]
    if (head.Channels > 1) {
        nelem *= head.Channels
    }

    data := make([]float32, nelem)
    err = binary.Read(conn, binary.LittleEndian, &data)
    if err != nil { return err }

    img := image.NewRGBA(image.Rect(0, 0, int(head.Matrix_size[0]), int(head.Matrix_size[1])))
    b := img.Bounds()
    for y := b.Min.Y; y < b.Max.Y; y++ {
        for x := b.Min.X; x < b.Max.X; x++ {
            v := data[(y * b.Max.X) + x]
            // Gadgetron AutoScale Gadget uses 2048 as maximum
            s := uint8(v / 2048 * 256)
            g := color.RGBA {s, s, s, 255}
            img.Set(x, y, g)
        }
    }

    fname := fmt.Sprintf("%d.png", r.count)
    err = savePNG(fname, img)
    if err != nil { return err }

    r.count++
    return nil
}

func savePNG(filename string, img image.Image) error {
    f, err := os.Create(filename)
    defer f.Close()
    if err != nil { return err }

    return png.Encode(f, img)
}
