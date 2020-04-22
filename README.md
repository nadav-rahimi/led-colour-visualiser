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

## TODO
#### Features
- [x] Ability to enable/disable dampening
- [x] Ability to enable/disable interpolation
- [ ] Rewrite portaudio and ui into their own interfaces so they can be interfaced with independently
- [ ] Add logging

#### Fixes
- [ ] Implement a median filter
- [ ] Convert colours to uint32 manually
- [x] Speed up interpolation and dampening
- [x] Reduce the processing of data to improve memory usage and speed

#### Ideas
- [x] Switch to float32 (done by converting the dsputils library to float32)
- [ ] Put emphasis on higher and lower frequencies
- [ ] Make a custom smaller colour range which will lend to smoother transitions
- [ ] Mode with frequency bands to reduce the number of colours used, e.g. for every 100 hz change colour

