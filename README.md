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
- [ ] Ability to enable/disable dampening
- [ ] Ability to enable/disable interpolation
- [ ] Rewrite portaudio and ui into their own interfaces so they can be interfaced with independently

#### Fixes
- [ ] Convert colours to uint32 manually
- [ ] Speed up interpolation and dampening
- [ ] Reduce the processing of data to improve memory usage and speed

#### Ideas
- [ ] Switch to float32?
- [ ] Put emphasis on higher and lower frequencies
- [ ] Make a custom smaller colour range which will lend to smoother transitions
- [ ] Mode with frequency bands to reduce the number of colours used, e.g. for every 100 hz change colour

