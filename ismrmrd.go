package main

const (
    ISMRMRD_VERSION = 1
    ISMRMRD_POSITION_LENGTH = 3
    ISMRMRD_DIRECTION_LENGTH = 3
    ISMRMRD_USER_INTS = 8
    ISMRMRD_USER_FLOATS = 8
    ISMRMRD_PHYS_STAMPS = 8  //TODO: This should be changed to 3 (Major impact)
    ISMRMRD_CHANNEL_MASKS = 16

    ACQ_FIRST_IN_SLICE = 1 << 6
    ACQ_LAST_IN_SLICE = 1 << 7
)

type EncodingCounters struct {
    Kspace_encode_step_1 uint16
    Kspace_encode_step_2 uint16
    Average uint16
    Slice uint16
    Contrast uint16
    Phase uint16
    Repetition uint16
    Set uint16
    Segment uint16
    User [ISMRMRD_USER_INTS]uint16
}

type AcquisitionHeader struct {
    Version uint16
    Flags uint64
    Measurement_uid uint32
    Scan_counter uint32
    Acquisition_time_stamp uint32
    Physiology_time_stamp [ISMRMRD_PHYS_STAMPS]uint32
    Number_of_samples uint16
    Available_channels uint16
    Active_channels uint16
    Channel_mask [ISMRMRD_CHANNEL_MASKS]uint64
    Discard_pre uint16
    Discard_post uint16
    Center_sample uint16
    Encoding_space_ref uint16
    Trajectory_dimensions uint16
    Sample_time_us float32
    Position [ISMRMRD_POSITION_LENGTH]float32
    Read_dir [ISMRMRD_DIRECTION_LENGTH]float32
    Phase_dir [ISMRMRD_DIRECTION_LENGTH]float32
    Slice_dir [ISMRMRD_DIRECTION_LENGTH]float32
    Patientable_position [ISMRMRD_POSITION_LENGTH]float32
    Idx EncodingCounters
    User_int [ISMRMRD_USER_INTS]int32
    User_float32 [ISMRMRD_USER_FLOATS]float32
}

type Acquisition struct {
    head AcquisitionHeader
    traj []float32
    data []float32
}

type ImageHeader struct {
    Version uint16
    Flags uint64
    Measurement_uid uint32
    Matrix_size [3]uint16
    Field_of_view [3]float32
    Channels uint16
    Position [ISMRMRD_POSITION_LENGTH]float32
    Read_dir [ISMRMRD_DIRECTION_LENGTH]float32
    Phase_dir [ISMRMRD_DIRECTION_LENGTH]float32
    Slice_dir [ISMRMRD_DIRECTION_LENGTH]float32
    Patient_table_position [ISMRMRD_POSITION_LENGTH]float32
    Average uint16
    Slice uint16
    Contrast uint16
    Phase uint16
    Repetition uint16
    Set uint16
    Acquisition_time_stamp uint32
    Physiology_time_stamp [ISMRMRD_PHYS_STAMPS]uint32
    Image_dataype uint16
    Image_type uint16
    Image_index uint16
    Image_series_index uint16
    User_int [ISMRMRD_USER_INTS]int32
    User_float [ISMRMRD_USER_FLOATS]float32
}
