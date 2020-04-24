# LED Visualiser

## Program Structure
1. Read in music from phone bluetooth/computer
2. Parse music and convert it to colour
3. Update the LEDs based on the colour calculated

## Dependencies
- https://github.com/andlabs/ui (the gui)
- https://github.com/gordonklaus/portaudio (portaudio bindings for go)
- https://github.com/lucasb-eyer/go-colorful (colour library)
- https://github.com/mjibson/go-dsp/fft (fft)
- https://github.com/wcharczuk/go-chart (the graphing tool)

## TODO
#### Features
- [x] Ability to enable/disable dampening
- [x] Abiity for custom colour range
- [ ] Add logging
- [ ] Ability to choose device to receive input from 
- [ ] Send the data via bluetooth to arduino
- [x] Ability to disable graphing
- [ ] Create gui for program to decide which option to enable/disable and ability to choose gradients and ability to choose input device

#### Fixes
- [x] Reduce the processing of data to improve memory usage and speed- [ ] Document the whole codebase
- [x] Speed up interpolation and dampening
- [x] Modularise code into functions
- [ ] Finish documentation
- [ ] Implement a median/kalman filter
- [ ] Remove the use of too many pointers
- [ ] Reduce type conversions as much as possible
- [ ] Set icon for the .exe - https://stackoverflow.com/questions/25602600/how-do-you-set-the-application-icon-in-golang


#### Ideas
- [x] Graph the frequencies to visualise difference in output
- [x] Add a smoothing algorithm in addition to dampening
- [x] Switch to float32 (done by converting the dsputils library to float32)
- [ ] Damp small changes in frequency but dont damp large changes in frequency, this will stop the bass visualisation lagging in Savage, Nights etc.
