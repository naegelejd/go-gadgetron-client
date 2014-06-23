package main

import (
    /* "fmt" */
    "log"
    "flag"
)

func main () {
    var host string
    var port int
    flag.StringVar(&host, "host", "", "hostname")
    flag.IntVar(&port, "port", 9002, "port number")
    flag.Parse()

    conn := new(GadgetronConnector)
    conn.Open(host, port)
    defer conn.Close()

    conn.sendGadgetronConfigurationFile("default.xml")

    conn.sendGadgetronParameters(xml_config)

    head := AcquisitionHeader {}
    head.version = 1
    head.number_of_samples = 256
    head.available_channels = 1
    head.active_channels = 1
    head.center_sample = 128
    head.sample_time_us = 5.0
    //log.Printf("%+v\n", head)

    // 5 slices (corresponding to slice limits in xml_config)
    for s := 0; s < 5; s++ {
        for i,line := range data {
            if (i == 0) {
                head.flags = ACQ_FIRST_IN_SLICE
            } else if (i == 127) {
                head.flags = ACQ_LAST_IN_SLICE
            }

            head.idx.kspace_encode_step_1 = uint16(i)
            head.idx.slice = uint16(s)
            acq := &Acquisition {head, traj, line}
            conn.sendIsmrmrdAcquisition(acq)
        }
    }

    conn.sendCloseMessage()

    conn.readImages()

    log.Println("Finished!")
}
