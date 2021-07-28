# sonometre
`sonometre` is a small command line tool that reads data from an inexpensive sound level meter, then submits the result to Datadog.

# Usage
```
DD_METRIC_NAME=home.backyard.sound_level DD_API_KEY=xxxxxxxxx ./sonometre
```

P.S. Only works on Linux (tested on Rapsberry Pi 4)

# Device
**Brand:** GAIN EXPRESS

**Model**: SLM-25

![61MPs+gIRIL _SL1500_](https://user-images.githubusercontent.com/1836721/127361441-a7f8074a-fab1-407c-a9e4-ee73bda6c799.jpg)

# Result
<img width="919" alt="Screenshot 2021-07-28 at 18 39 39" src="https://user-images.githubusercontent.com/1836721/127362546-e44e4ad5-806d-4909-b7f6-2cb21af1a272.png">
