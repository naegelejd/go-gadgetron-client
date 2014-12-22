package main

import (
	"flag"
	ismrmrd "github.com/naegelejd/go-ismrmrd"
	"log"
)

func main() {
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

	head := ismrmrd.AcquisitionHeader{}
	head.Version = 1
	head.NumberOfSamples = 256
	head.AvailableChannels = 1
	head.ActiveChannels = 1
	head.CenterSample = 128
	head.SampleTimeUs = 5.0
	//log.Printf("%+v\n", head)

	var traj = make([]float32, 0)
	var data [][]float32 = makeKSpaceData()

	// 5 slices (corresponding to slice limits in xml_config)
	for s := 0; s < 5; s++ {
		for i, line := range data {
			if i == 0 {
				head.Flags = ismrmrd.ACQ_FIRST_IN_SLICE
			} else if i == 127 {
				head.Flags = ismrmrd.ACQ_LAST_IN_SLICE
			}

			head.Idx.KspaceEncodeStep1 = uint16(i)
			head.Idx.Slice = uint16(s)
			acq := &ismrmrd.Acquisition{head, traj, line}
			conn.sendIsmrmrdAcquisition(acq)
		}
	}

	conn.sendCloseMessage()

	conn.readImages()

	log.Println("Finished!")
}
