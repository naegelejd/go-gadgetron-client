package main

import (
    "fmt"
    "net"
    "log"
    "bytes"
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
    readers map[uint16]GadgetMessageReader
}

func newGadgetronConnector(host string, port int) *GadgetronConnector {
    c := new(GadgetronConnector)

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

    c.readers = make(map[uint16]GadgetMessageReader, 8)

    return c
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

func (g *GadgetronConnector) registerReader(msg uint16, reader GadgetMessageReader) {
    g.readers[msg] = reader
}

func (g *GadgetronConnector) readImages() {
    for {
        var msg uint16
        binary.Read(g.conn, binary.LittleEndian, &msg)

        if msg == GADGET_MESSAGE_CLOSE {
            log.Println("Received GADGET_MESSAGE_CLOSE")
            return
        } else {
            reader := g.readers[msg]
            if reader != nil {
                err := reader.Read(g.conn)
                if err != nil {
                    log.Println(err)
                }
            } else {
                log.Printf("No reader registered for message (%d)", msg)
            }
        }
    }
}
