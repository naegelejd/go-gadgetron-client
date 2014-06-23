package main

import (
    "os"
    "fmt"
    "net"
    "log"
    "bytes"
    "image"
    "image/png"
    "image/color"
    "encoding/binary"
)

const (
    GADGET_MESSAGE_CONFIG_FILE uint16 = 1
    GADGET_MESSAGE_PARAMETER_SCRIPT uint16 = 3
    GADGET_MESSAGE_CLOSE uint16 = 4

    GADGET_MESSAGE_ISMRMRD_ACQUISITION uint16 = 1008
    GADGET_MESSAGE_ISMRMRD_IMAGE_REAL_FLOAT uint16 = 1010
)

type GadgetronConnector struct {
    conn net.Conn
    /* readers map[uint16]GadgetMessageReader */
}

func (c *GadgetronConnector) Open(host string, port int) {
    addr := fmt.Sprintf("%s:%d", host, port)
    tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil {
        log.Fatal("Error: Could not resolve address", addr)
    }

    netconn, err := net.Dial(tcpAddr.Network(), tcpAddr.String())
    if err != nil {
        log.Fatal(err)
    }

    c.conn = netconn
}

func (c *GadgetronConnector) Close() {
    c.conn.Close()
}

func (c *GadgetronConnector) sendCloseMessage() {
    // Tell the Gadgetron we're finished sending
    err := binary.Write(c.conn, binary.LittleEndian, GADGET_MESSAGE_CLOSE)
    if err != nil {
        log.Fatal("Error writing GADGET_MESSAGE_CLOSE")
    }
}

func (c *GadgetronConnector) sendGadgetronConfigurationFile(config_xml_name string) {
    // Tell Gadgetron we're sending the configuration
    err := binary.Write(c.conn, binary.LittleEndian, GADGET_MESSAGE_CONFIG_FILE)
    if err != nil {
        log.Fatal("Error writing GADGET_MESSAGE_CONFIG_FILE")
    }

    var buffer bytes.Buffer
    buffer.WriteString(config_xml_name)
    buffer.Write(make([]byte, 1024 - len(config_xml_name)))
    config_file := buffer.Bytes()
    log.Println("config_file:", string(config_file))
    // Send the Gadgetron config filename
    c.conn.Write(config_file)
}

func (c *GadgetronConnector) sendGadgetronParameters(xml_string string) {
    // Tell Gadgetron we're sending an XML config
    err := binary.Write(c.conn, binary.LittleEndian, GADGET_MESSAGE_PARAMETER_SCRIPT)
    if err != nil {
        log.Fatal("Error writing GADGET_MESSAGE_PARAMETER_SCRIPT")
    }

    var buffer bytes.Buffer
    buffer.WriteString(xml_config)
    buffer.WriteByte(0)
    xml_config_bytes := buffer.Bytes()
    // Send the length of the XML config
    err = binary.Write(c.conn, binary.LittleEndian, uint32(len(xml_config_bytes)))
    if err != nil {
        log.Fatal("Error writing XML script length")
    }
    // Send the XML config
    c.conn.Write(xml_config_bytes)
}

func (c *GadgetronConnector) sendIsmrmrdAcquisition(a *Acquisition) {
    err := binary.Write(c.conn, binary.LittleEndian, GADGET_MESSAGE_ISMRMRD_ACQUISITION)
    if err != nil {
        log.Fatal("Error writing GADGET_MESSAGE_ISMRMRD_ACQUISITION")
    }

    err = binary.Write(c.conn, binary.LittleEndian, a.head)
    if err != nil {
        log.Fatal("Error writing AcquisitionHeader")
    }

    err = binary.Write(c.conn, binary.LittleEndian, a.traj)
    if err != nil {
        log.Fatal("Error writing Acquisition trajectory")
    }

    err = binary.Write(c.conn, binary.LittleEndian, a.data)
    if err != nil {
        log.Fatal("Error writing Acquisition data")
    }
}

func (g *GadgetronConnector) readImages() {
    for i := 0; ; i++ {
        var msg uint16
        binary.Read(g.conn, binary.LittleEndian, &msg)

        if msg == GADGET_MESSAGE_CLOSE {
            log.Println("Received GADGET_MESSAGE_CLOSE")
            break
        } else if msg != GADGET_MESSAGE_ISMRMRD_IMAGE_REAL_FLOAT {
            log.Printf("Read Gadget Message (%d), expecting GADGET_MESSAGE_IMAGE_REAL_FLOAT", msg)
            break
        }

        var head ImageHeader
        err := binary.Read(g.conn, binary.LittleEndian, &head)
        if err != nil {
            log.Fatal("Failed to read ImageHeader")
        }
        //fmt.Printf("%+v\n", head)

        nelem := head.Matrix_size[0] * head.Matrix_size[1] * head.Matrix_size[2]
        if (head.Channels > 1) {
            nelem *= head.Channels
        }

        data := make([]float32, nelem)
        err = binary.Read(g.conn, binary.LittleEndian, &data)
        if err != nil {
            log.Fatal("Failed to read image data")
        }

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

        fname := fmt.Sprintf("%d.png", i)
        f, err := os.Create(fname)
        defer f.Close()
        if err != nil {
            log.Fatal(err)
        }

        png.Encode(f, img)
    }
}
