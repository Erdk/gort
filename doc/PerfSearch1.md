## Baseline

BenchmarkGort-4   	       1	422333738242 ns/op
PASS
ok  	github.com/Erdk/gort	422.339s

## Separate rand for each goroutine

BenchmarkGort-4   	       1	385376617431 ns/op
PASS
ok  	github.com/Erdk/gort	385.379s

## Materials passed by pointer, not value

BenchmarkGort-4   	       1	385974802808 ns/op
PASS
ok  	github.com/Erdk/gort	385.977s

## Smarter job division (16x16 chunk)

BenchmarkGort-4   	       1	344897325182 ns/op
PASS
ok  	github.com/Erdk/gort	344.899s

## Smarter job division (32x32 chunk)

BenchmarkGort-4   	       1	346755674935 ns/op
PASS
ok  	github.com/Erdk/gort	346.759s
