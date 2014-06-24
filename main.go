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
