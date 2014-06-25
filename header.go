package main

const xml_config = `<?xml version="1.0" encoding="UTF-8" standalone="no" ?>
<ismrmrdHeader xmlns="http://www.ismrm.org/ISMRMRD" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.ismrm.org/ISMRMRD ismrmrd.xsd">

  <experimentalConditions>
    <H1resonanceFrequency_Hz>63500000</H1resonanceFrequency_Hz>
  </experimentalConditions>

  <encoding>
    <encodedSpace>
      <matrixSize>
        <x>256</x>
        <y>128</y>
        <z>1</z>
      </matrixSize>
      <fieldOfView_mm>
        <x>600</x>
        <y>300</y>
        <z>6</z>
      </fieldOfView_mm>
    </encodedSpace>
    <reconSpace>
      <matrixSize>
        <x>128</x>
        <y>128</y>
        <z>1</z>
      </matrixSize>
      <fieldOfView_mm>
        <x>300</x>
        <y>300</y>
        <z>6</z>
      </fieldOfView_mm>
    </reconSpace>
    <encodingLimits>
      <kspace_encoding_step_1>
        <minimum>0</minimum>
        <maximum>127</maximum>
        <center>64</center>
      </kspace_encoding_step_1>
      <slice>
        <minimum>0</minimum>
        <maximum>4</maximum>
        <center>2</center>
      </slice>
    </encodingLimits>
    <trajectory>cartesian</trajectory>
  </encoding>

</ismrmrdHeader>
`
