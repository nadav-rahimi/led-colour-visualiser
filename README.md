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
- [ ] Rewrite portaudio and ui into their own interfaces so they can be interfaced with independently
- [ ] Add logging
- [ ] Ability to choose device to receive input from 

#### Fixes
- [ ] Simulate momentum i.e. if between the new and last frequency there is a big change in magnitude show it, if its smaller then smooth it out, this wil keep the peaks high while keeping the small changes small
- [ ] Implement a median filter
- [ ] Convert colours to uint32 manually
- [x] Speed up interpolation and dampening
- [x] Reduce the processing of data to improve memory usage and speed

#### Ideas
- [ ] Graph the frequencies to visualise difference in output
- [ ] Makes the peaks stand out more and make the rest more similar to reduce crazy frequency shifting at high and mid ranges where frequency changes constantly
- [x] Add a smoothing algorithm in addition to dampening
- [x] Switch to float32 (done by converting the dsputils library to float32)
- [ ] Put emphasis on higher and lower frequencies
- [ ] Make a custom smaller colour range which will lend to smoother transitions
- [ ] Mode with frequency bands to reduce the number of colours used, e.g. for every 100 hz change colour

