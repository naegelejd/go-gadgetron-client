import numpy as np
import matplotlib.pyplot as plt

# Generates a white rectangle (ismrmrd test_create_dataset.cpp)
nX, nY = 256, 128
rho = np.zeros((nY, nX))

x0, x1 = nX / 4, 3 * nX / 4
y0, y1 = nY / 8, 7 * nY / 8

rho[y0:y1, x0:x1] = 1

K = np.fft.fftshift(np.fft.fft2(np.fft.fftshift(rho)))

# plt.imshow(rho)
# plt.show()
# print(K.shape)

with open('data.go', 'w') as f:
    f.write('package main\n\n')
    f.write("""
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

""")
    f.write('var traj = make([]float32, 0)\n\n')
    f.write('var data [][]float32 = makeData()\n')

    # f.write('var data [][]float32 = [][]float32 {\n')
    # for y, line in enumerate(K):
    #     s = '[]float32 { '
    #     for x, sample in enumerate(line):
    #         s += '%f, %f, ' % (sample.real, sample.imag)
    #         if x % 7 == 0:
    #             f.write(s)
    #             f.write('\n')
    #             s = ''
    #     s = s[:-2]
    #     s += '},\n'
    #     f.write(s)
    # f.write('}\n')
