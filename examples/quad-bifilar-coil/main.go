// quad-bifilar-coil creates Gerber files (and a bundled ZIP) representing
// four bifilar coils (https://en.wikipedia.org/wiki/Bifilar_coil) (one on
// each layer of a four-layer PCB) for manufacture on a printed circuit
// board (PCB).
package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	_ "github.com/gmlewis/go-fonts/fonts/freeserif"
	. "github.com/gmlewis/go-gerber/gerber"
)

var (
	step     = flag.Float64("step", 0.02, "Resolution (in radians) of the spiral")
	n        = flag.Int("n", 100, "Number of full winds in each spiral")
	gap      = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace    = flag.Float64("trace", 0.15, "Width of traces in mm")
	prefix   = flag.String("prefix", "quad-bifilar-coil", "Filename prefix for all Gerber files and zip")
	fontName = flag.String("font", "freeserif", "Name of font to use for writing source on PCB (empty to not write)")
	pts      = flag.Float64("pts", 18.0, "Font point size (72 pts = 1 inch = 25.4 mm)")
)

const (
	message = `With a trace and gap size of 0.15mm, this
quad bifilar coil should have a DC resistance
of approx. 928.8Ω. Each spiral has 100 coils.`
	message2 = `Top layer: hole5 ⇨ hole6
Bottom layer: hole6 ⇨ hole2
Top layer: hole2 ⇨ hole7
Bottom layer: hole7 ⇨ hole4
Layer 3: hole4 ⇨ hole3
Layer 2: hole3 ⇨ hole8
Layer 3: hole8 ⇨ hole1
Layer 2: hole1 ⇨ hole9`
)

func main() {
	flag.Parse()

	g := New(*prefix)

	s := newSpiral()

	topSpiralR, endR := s.genSpiral(1.0, 0, 0)
	topSpiralL, endL := s.genSpiral(1.0, math.Pi, 0)
	botSpiralR, _ := s.genSpiral(-1.0, 0, 0)
	startR := genPt(1.0, s.startAngle, 0, 0)
	// endR := genPt(1.0, s.endAngle, 0, 0)
	startL := genPt(1.0, s.startAngle, 0, math.Pi)
	// endL := genPt(1.0, s.endAngle, 0, math.Pi)

	viaPadD := 0.5
	viaDrillD := 0.25
	// viaPadOffset := 0.5 * (viaPadD - *trace)

	padD := 2.0
	drillD := 1.0
	botSpiralL, botEndL := s.genSpiral(-1.0, math.Pi, *trace+padD)

	shiftAngle := 0.5 * math.Pi
	layer2SpiralR, layer2EndR := s.genSpiral(1.0, shiftAngle, 0)
	layer2SpiralL, layer2EndL := s.genSpiral(1.0, math.Pi+shiftAngle, -(*trace + padD))
	layer3SpiralR, _ := s.genSpiral(-1.0, shiftAngle, 0)
	layer3SpiralL, layer3EndL := s.genSpiral(-1.0, math.Pi+shiftAngle, *trace+padD)
	startLayer2R := genPt(1.0, s.startAngle, 0, shiftAngle)
	//	endLayer2R := genPt(1.0, s.endAngle, 0, shiftAngle)
	startLayer2L := genPt(1.0, s.startAngle, 0, math.Pi+shiftAngle)
	//	endLayer2L := genPt(1.0, s.endAngle, 0, math.Pi+shiftAngle)

	viaOffset := math.Sqrt(0.5 * (*trace + viaPadD) * (*trace + viaPadD))
	hole2Offset := 0.5 * (*trace + viaPadD)
	hole4PadOffset := 0.5 * (viaPadD + *trace)
	hole5PadOffset := 0.5 * (padD + *trace)

	// Lower connecting trace between two spirals
	// hole1 := Point(startR.X, startR.Y+viaPadOffset)
	hole1 := Point(viaOffset, 0)
	hole2 := Point(endL.X-hole2Offset, endL.Y)
	// Upper connecting trace for left spiral
	// hole3 := Point(startL.X, startL.Y-viaPadOffset)
	hole3 := Point(-viaOffset, 0)
	hole4 := Point(endR.X+hole4PadOffset, *trace+padD)
	// Lower connecting trace for right spiral
	hole5 := Point(endR.X+hole5PadOffset, endR.Y)
	// Layer 2 and 3 inner connecting holes
	hole6 := Point(0, viaOffset)
	hole7 := Point(0, -viaOffset)
	// Layer 2 and 3 outer connecting hole
	hole8 := Point(layer2EndR.X, layer2EndR.Y+hole2Offset)
	hole9 := Point(layer2EndL.X+hole5PadOffset, -(*trace + padD))

	top := g.TopCopper()
	top.Add(
		Polygon(0, 0, true, topSpiralR, 0.0),
		Polygon(0, 0, true, topSpiralL, 0.0),
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, viaPadD),
		Circle(hole2.X, hole2.Y, viaPadD),
		Line(endL.X, endL.Y, hole2.X, hole2.Y, RectShape, *trace),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, viaPadD),
		Circle(hole4.X, hole4.Y, viaPadD),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, padD),
		Line(endR.X, endR.Y, hole5.X, hole5.Y, RectShape, *trace),
		// Layer 2 and 3 inner connecting holes
		Circle(hole6.X, hole6.Y, viaPadD),
		Circle(hole7.X, hole7.Y, viaPadD),
		Line(startR.X, startR.Y, hole6.X, hole6.Y, RectShape, *trace),
		Line(startL.X, startL.Y, hole7.X, hole7.Y, RectShape, *trace),
		// Layer 2 and 3 outer connecting hole
		Circle(hole8.X, hole8.Y, viaPadD),
		Circle(hole9.X, hole9.Y, padD),
	)

	layer2 := g.Layer2()
	layer2.Add(
		Polygon(0, 0, true, layer2SpiralR, 0.0),
		Polygon(0, 0, true, layer2SpiralL, 0.0),
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, viaPadD),
		Circle(hole2.X, hole2.Y, viaPadD),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, viaPadD),
		Circle(hole4.X, hole4.Y, viaPadD),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, padD),
		// Layer 2 and 3 inner connecting holes
		Circle(hole6.X, hole6.Y, viaPadD),
		Circle(hole7.X, hole7.Y, viaPadD),
		Circle(hole1.X, hole1.Y, viaPadD),
		Circle(hole3.X, hole3.Y, viaPadD),
		Line(startLayer2R.X, startLayer2R.Y, hole3.X, hole3.Y, RectShape, *trace),
		Line(startLayer2L.X, startLayer2L.Y, hole1.X, hole1.Y, RectShape, *trace),
		// Layer 2 and 3 outer connecting hole
		Circle(hole8.X, hole8.Y, viaPadD),
		Line(layer2EndR.X, layer2EndR.Y, hole8.X, hole8.Y, RectShape, *trace),
		Circle(hole9.X, hole9.Y, padD),
		Line(layer2EndL.X, layer2EndL.Y, hole9.X, hole9.Y, RectShape, *trace),
	)

	topMask := g.TopSolderMask()
	topMask.Add(
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, viaPadD),
		Circle(hole2.X, hole2.Y, viaPadD),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, viaPadD),
		Circle(hole4.X, hole4.Y, viaPadD),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, padD),
		// Layer 2 and 3 inner connecting holes
		Circle(hole6.X, hole6.Y, viaPadD),
		Circle(hole7.X, hole7.Y, viaPadD),
		// Layer 2 and 3 outer connecting hole
		Circle(hole8.X, hole8.Y, viaPadD),
		Circle(hole9.X, hole9.Y, padD),
	)

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(0, 0, true, botSpiralR, 0.0),
		Polygon(0, 0, true, botSpiralL, 0.0),
		Line(endL.X, endL.Y, hole2.X, hole2.Y, RectShape, *trace),
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, viaPadD),
		Circle(hole2.X, hole2.Y, viaPadD),
		Line(endL.X, endL.Y, hole2.X, hole2.Y, RectShape, *trace),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, viaPadD),
		Circle(hole4.X, hole4.Y, viaPadD),
		Line(botEndL.X, botEndL.Y, hole4.X, hole4.Y, RectShape, *trace),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, padD),
		// Layer 2 and 3 inner connecting holes
		Circle(hole6.X, hole6.Y, viaPadD),
		Circle(hole7.X, hole7.Y, viaPadD),
		Line(startR.X, startR.Y, hole6.X, hole6.Y, RectShape, *trace),
		Line(startL.X, startL.Y, hole7.X, hole7.Y, RectShape, *trace),
		// Layer 2 and 3 outer connecting hole
		Circle(hole8.X, hole8.Y, viaPadD),
		Circle(hole9.X, hole9.Y, padD),
	)

	layer3 := g.Layer3()
	layer3.Add(
		Polygon(0, 0, true, layer3SpiralR, 0.0),
		Polygon(0, 0, true, layer3SpiralL, 0.0),
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, viaPadD),
		Circle(hole2.X, hole2.Y, viaPadD),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, viaPadD),
		Circle(hole4.X, hole4.Y, viaPadD),
		Line(layer3EndL.X, layer3EndL.Y, hole4.X, hole4.Y, RectShape, *trace),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, padD),
		// Layer 2 and 3 inner connecting holes
		Circle(hole6.X, hole6.Y, viaPadD),
		Circle(hole7.X, hole7.Y, viaPadD),
		Circle(hole1.X, hole1.Y, viaPadD),
		Circle(hole3.X, hole3.Y, viaPadD),
		Line(startLayer2R.X, startLayer2R.Y, hole3.X, hole3.Y, RectShape, *trace),
		Line(startLayer2L.X, startLayer2L.Y, hole1.X, hole1.Y, RectShape, *trace),
		// Layer 2 and 3 outer connecting hole
		Circle(hole8.X, hole8.Y, viaPadD),
		Line(layer2EndR.X, layer2EndR.Y, hole8.X, hole8.Y, RectShape, *trace),
		Circle(hole9.X, hole9.Y, padD),
	)

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, viaPadD),
		Circle(hole2.X, hole2.Y, viaPadD),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, viaPadD),
		Circle(hole4.X, hole4.Y, viaPadD),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, padD),
		// Layer 2 and 3 inner connecting holes
		Circle(hole6.X, hole6.Y, viaPadD),
		Circle(hole7.X, hole7.Y, viaPadD),
		// Layer 2 and 3 outer connecting hole
		Circle(hole8.X, hole8.Y, viaPadD),
		Circle(hole9.X, hole9.Y, padD),
	)

	drill := g.Drill()
	drill.Add(
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, viaDrillD),
		Circle(hole2.X, hole2.Y, viaDrillD),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, viaDrillD),
		Circle(hole4.X, hole4.Y, viaDrillD),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, drillD),
		// Layer 2 and 3 inner connecting holes
		Circle(hole6.X, hole6.Y, viaDrillD),
		Circle(hole7.X, hole7.Y, viaDrillD),
		// Layer 2 and 3 outer connecting hole
		Circle(hole8.X, hole8.Y, viaDrillD),
		Circle(hole9.X, hole9.Y, drillD),
	)

	outline := g.Outline()
	r := 0.5*s.size + padD + *trace
	outline.Add(
		Arc(0, 0, r, CircleShape, 1, 1, 0, 360, 0.1),
	)
	fmt.Printf("n=%v: (%.2f,%.2f)\n", *n, 2*r, 2*r)

	if *fontName != "" {
		radius := -endL.X
		labelSize := 6.0

		tss := g.TopSilkscreen()
		tss.Add(
			Text(0, 0.3*radius, 1.0, message, *fontName, *pts, Center),
			Text(hole1.X+viaPadD, hole1.Y, 1.0, "hole1", *fontName, labelSize, CenterLeft),
			Text(hole2.X+viaPadD, hole2.Y, 1.0, "hole2", *fontName, labelSize, CenterLeft),
			Text(hole3.X-viaPadD, hole3.Y, 1.0, "hole3", *fontName, labelSize, CenterRight),
			Text(hole4.X-padD, hole4.Y, 1.0, "hole4", *fontName, labelSize, CenterRight),
			Text(hole5.X-padD, hole5.Y, 1.0, "hole5", *fontName, labelSize, CenterRight),
			Text(hole6.X, hole6.Y+viaPadD, 1.0, "hole6", *fontName, labelSize, BottomCenter),
			Text(hole7.X, hole7.Y-viaPadD, 1.0, "hole7", *fontName, labelSize, TopCenter),
			Text(hole8.X, hole8.Y-viaPadD, 1.0, "hole8", *fontName, labelSize, TopCenter),
			Text(hole9.X-padD, hole9.Y, 1.0, "hole9", *fontName, labelSize, CenterRight),
			Text(0, -0.5*radius, 1.0, message2, *fontName, *pts, Center),
		)
	}

	if err := g.WriteGerber(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
}

func genPt(xScale, angle, halfTW, offset float64) Pt {
	r := (angle + *trace + *gap) / (3 * math.Pi)
	x := (r + halfTW) * math.Cos(angle+offset)
	y := (r + halfTW) * math.Sin(angle+offset)
	return Point(x*xScale, y)
}

type spiral struct {
	startAngle float64
	endAngle   float64
	size       float64
}

func newSpiral() *spiral {
	startAngle := 2.5 * math.Pi
	endAngle := 2*math.Pi + float64(*n)*2.0*math.Pi
	p1 := genPt(1.0, endAngle, *trace*0.5, 0)
	size := 2 * math.Abs(p1.X)
	p2 := genPt(1.0, endAngle, *trace*0.5, math.Pi)
	xl := 2 * math.Abs(p2.X)
	if xl > size {
		size = xl
	}
	return &spiral{
		startAngle: startAngle,
		endAngle:   endAngle,
		size:       size,
	}
}

func (s *spiral) genSpiral(xScale, offset, trimY float64) (pts []Pt, endPt Pt) {
	halfTW := *trace * 0.5
	endAngle := s.endAngle
	if trimY < 0 { // Only for layer2SpiralL - extend another Pi/2
		endAngle += 0.5 * math.Pi
	}
	steps := int(0.5 + (endAngle-s.startAngle) / *step)
	for i := 0; i < steps; i++ {
		angle := s.startAngle + *step*float64(i)
		pts = append(pts, genPt(xScale, angle, halfTW, offset))
	}
	var trimYsteps int
	if trimY > 0 {
		trimYsteps++
		for {
			if pts[len(pts)-trimYsteps].Y > trimY {
				break
			}
			trimYsteps++
		}
		lastStep := len(pts) - trimYsteps
		trimYsteps--
		pts = pts[0 : lastStep+1]
		pts = append(pts, Pt{X: pts[lastStep].X, Y: trimY})
		angle := s.startAngle + *step*float64(steps-1-trimYsteps)
		eX := genPt(xScale, angle, 0, offset)
		endPt = Pt{X: eX.X, Y: trimY}
		nX := genPt(xScale, angle, -halfTW, offset)
		pts = append(pts, Pt{X: nX.X, Y: trimY})
	} else if trimY < 0 { // Only for layer2SpiralL
		trimYsteps++
		for {
			if pts[len(pts)-trimYsteps].Y < trimY {
				break
			}
			trimYsteps++
		}
		lastStep := len(pts) - trimYsteps
		trimYsteps--
		pts = pts[0 : lastStep+1]
		pts = append(pts, Pt{X: pts[lastStep].X, Y: trimY})
		angle := s.startAngle + *step*float64(steps-1-trimYsteps)
		eX := genPt(xScale, angle, 0, offset)
		endPt = Pt{X: eX.X, Y: trimY}
		nX := genPt(xScale, angle, -halfTW, offset)
		pts = append(pts, Pt{X: nX.X, Y: trimY})
	} else {
		pts = append(pts, genPt(xScale, endAngle, halfTW, offset))
		endPt = genPt(xScale, endAngle, 0, offset)
		pts = append(pts, genPt(xScale, endAngle, -halfTW, offset))
	}
	for i := steps - 1 - trimYsteps; i >= 0; i-- {
		angle := s.startAngle + *step*float64(i)
		pts = append(pts, genPt(xScale, angle, -halfTW, offset))
	}
	pts = append(pts, pts[0])
	return pts, endPt
}
