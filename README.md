# `gax`

[![Apache License](https://img.shields.io/github/license/octu0/gax)](https://github.com/octu0/gax/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/octu0/gax?status.svg)](https://godoc.org/github.com/octu0/gax)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/gax)](https://goreportcard.com/report/github.com/octu0/gax)
[![Releases](https://img.shields.io/github/v/release/octu0/gax)](https://github.com/octu0/gax/releases)

`gax` is an image processing framework.  
`gax` provides a [Halide](https://halide-lang.org) like interface and reduces size of the image processing code.

## Overview

### Halide Code

```c++
using namespace Halide;

Func blur_3x3(Func in) {
	Var x, y, c;

	Func blur_x;
	blur_x(x, y, c) = (in(x-1, y, c) + in(x, y, c) + in(x+1, y, c))/3;

	Func blur_y;
	blur_y(x, y, c) = (blur_x(x, y-1, c) + blur_x(x, y, c) + blur_x(x, y+1, c))/3;

	return blur_y;
}

void main() {
	Buffer<uint8_t> input = load_image("images/rgb.png");

	Func blur = blur_3x3(input);

	Buffer<uint8_t> output = blur.realize({input.width(), input.height(), input.channels()});

	save_image(output, "brighter.png");
}
```

### gax Code

```go
package main

import (
	"github.com/octu0/gax"
)

func main() {
	img, err := readImage("./src.png")
	if err != nil {
		panic(err)
	}

	in := gax.ImageRGBA[uint8](img)

	blurX := gax.Func[uint8]()
	blurX.SetXYC(func(x, y, ch gax.Var) uint8 {
		val1 := in.XYC(x-1, y, ch)
		val2 := in.XYC(x, y, ch)
		val3 := in.XYC(x+1, y, ch)
		return (val1 + val2 + val3) / 3
	})

	blurY := gax.Func[uint8]()
	blurY.SetXYC(func(x, y, ch gax.Var) uint8 {
		val1 := blurX.XYC(x, y-1, ch)
		val2 := blurX.XYC(x, y, ch)
		val3 := blurX.XYC(x, y+1, ch)
		return (val1 + val2 + val3) / 3
	})

	out := blurY.Realize(img.Rect.Dx(), img.Rect.Dy(), 4) // 4 = r,g,b,a
	outImg := &image.RGBA{
		Pix:    out,
		Stride: 4 * img.Rect.Dx(),
		Rect:   image.Rect(0, 0, img.Rect.Dx(), img.Rect.Dy()),
	}

	saveImage(outImg)
}
```

Running this code will output:

| img           |                                     |
| :-----------: | :---------------------------------: |
| input         | ![img](_example/src.png)            |
| output        | ![img](_example/blur.png)           |

# License

MIT, see LICENSE file for details.
