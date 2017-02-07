# gort
[![Go Report Card](https://goreportcard.com/badge/github.com/Erdk/gort)](https://goreportcard.com/report/github.com/Erdk/gort)

Simple ray tracer written in Go.

![alt tag](https://raw.githubusercontent.com/Erdk/gort/master/static/output.png)

## Instalation

```
go get github.com/Erdk/gort
```

## Rationale

Gort started as my "pet project" I'm working on in my spare time, 
and I belive it'll stay this way for now, so please don't expect 
regular updates ;) The idea came to my  mind after seeing Peter 
Shirley's books on Amazon about ray-tracing: "Ray Tracing in One Weekend", 
"Ray Tracing: The Next Week" and 
 "Ray Tracing: The Rest of Your Life". I thought that this would be a 
 good way to improve my programming skills and learn more about Golang
 ecosystem and best practises. Soo.. here we are :) As for books I 
 highly recommend them, they're well written and quite easy to understand.

## Options

```
    -w <width>
        width of generated image, by default 640.
    -h <height>
        height of generated image, by default 480.
    -s <samples-per-pixel>
        number of rays per pixel, by default 200.
    -t <threads>
        number of parallel rendering jobs, by default 2.
    -o <output-file>
        filename (without extension) of png with output, by default "output".
    -i <input-file>
        filename of input file, if this flag is used then scene won't be randomly generated. By default empty.
    -j <save-scene>
        save genereated scene before render to <output>.json, by default false.
    -prof <profile>
        generate profiling profile to use with 'go tool pprof'. By default none (no profiling). Available profiles:
        - cpu 
        - mem
        - block
```
