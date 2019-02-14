package gerber

import (
	"math"
	"testing"

	_ "github.com/gmlewis/go-fonts/fonts/freeserif"
)

func TestTextT_Primitive(t *testing.T) {
	var p Primitive = &TextT{}
	if p == nil {
		// In actuality, this test won't compile if it isn't a Primitive.
		t.Errorf("TextT does not implement the Primitive interface")
	}
}

func TestText(t *testing.T) {
	const (
		message    = "012"
		fontName   = "freeserif"
		wantWidth  = 36.8554
		wantHeight = 17.526
		pts        = 72
		eps        = 1e-6
		sf         = 1e6
	)

	tests := []struct {
		name     string
		x, y     float64
		opts     TextOpts
		wantXmin float64
		wantYmin float64
		wantXmax float64
		wantYmax float64
	}{
		{
			name:     "XLeft,YBottom",
			wantXmin: 0,
			wantYmin: 0,
			wantXmax: wantWidth * sf,
			wantYmax: wantHeight * sf,
		},
		{
			name:     "XLeft,YBottom w/ offset",
			x:        10,
			y:        20,
			wantXmin: 10 * sf,
			wantYmin: 20 * sf,
			wantXmax: (10 + wantWidth) * sf,
			wantYmax: (20 + wantHeight) * sf,
		},
		{
			name:     "XCenter,YBottom",
			opts:     BottomCenter,
			wantXmin: -0.5 * wantWidth * sf,
			wantYmin: 0,
			wantXmax: 0.5 * wantWidth * sf,
			wantYmax: wantHeight * sf,
		},
		{
			name:     "XCenter,YBottom w/ offset",
			x:        10,
			y:        20,
			opts:     BottomCenter,
			wantXmin: (10 - 0.5*wantWidth) * sf,
			wantYmin: 20 * sf,
			wantXmax: (10 + 0.5*wantWidth) * sf,
			wantYmax: (20 + wantHeight) * sf,
		},
		{
			name:     "XRight,YBottom",
			opts:     BottomRight,
			wantXmin: -wantWidth * sf,
			wantYmin: 0,
			wantXmax: 0,
			wantYmax: wantHeight * sf,
		},
		{
			name:     "XRight,YBottom w/ offset",
			x:        10,
			y:        20,
			opts:     BottomRight,
			wantXmin: (10 - wantWidth) * sf,
			wantYmin: 20 * sf,
			wantXmax: 10 * sf,
			wantYmax: (20 + wantHeight) * sf,
		},
		{
			name:     "XLeft,YCenter",
			opts:     CenterLeft,
			wantXmin: 0,
			wantYmin: -0.5 * wantHeight * sf,
			wantXmax: wantWidth * sf,
			wantYmax: 0.5 * wantHeight * sf,
		},
		{
			name:     "XLeft,YCenter w/ offset",
			x:        10,
			y:        20,
			opts:     CenterLeft,
			wantXmin: 10 * sf,
			wantYmin: (20 - 0.5*wantHeight) * sf,
			wantXmax: (10 + wantWidth) * sf,
			wantYmax: (20 + 0.5*wantHeight) * sf,
		},
		{
			name:     "XCenter,YCenter",
			opts:     Center,
			wantXmin: -0.5 * wantWidth * sf,
			wantYmin: -0.5 * wantHeight * sf,
			wantXmax: 0.5 * wantWidth * sf,
			wantYmax: 0.5 * wantHeight * sf,
		},
		{
			name:     "XCenter,YCenter w/ offset",
			x:        10,
			y:        20,
			opts:     Center,
			wantXmin: (10 - 0.5*wantWidth) * sf,
			wantYmin: (20 - 0.5*wantHeight) * sf,
			wantXmax: (10 + 0.5*wantWidth) * sf,
			wantYmax: (20 + 0.5*wantHeight) * sf,
		},
		{
			name:     "XRight,YCenter",
			opts:     CenterRight,
			wantXmin: -wantWidth * sf,
			wantYmin: -0.5 * wantHeight * sf,
			wantXmax: 0,
			wantYmax: 0.5 * wantHeight * sf,
		},
		{
			name:     "XRight,YCenter w/ offset",
			x:        10,
			y:        20,
			opts:     CenterRight,
			wantXmin: (10 - wantWidth) * sf,
			wantYmin: (20 - 0.5*wantHeight) * sf,
			wantXmax: 10 * sf,
			wantYmax: (20 + 0.5*wantHeight) * sf,
		},
		{
			name:     "XLeft,YTop",
			opts:     TopLeft,
			wantXmin: 0,
			wantYmin: -wantHeight * sf,
			wantXmax: wantWidth * sf,
			wantYmax: 0,
		},
		{
			name:     "XLeft,YTop w/ offset",
			x:        10,
			y:        20,
			opts:     TopLeft,
			wantXmin: 10 * sf,
			wantYmin: (20 - wantHeight) * sf,
			wantXmax: (10 + wantWidth) * sf,
			wantYmax: 20 * sf,
		},
		{
			name:     "XCenter,YTop",
			opts:     TopCenter,
			wantXmin: -0.5 * wantWidth * sf,
			wantYmin: -wantHeight * sf,
			wantXmax: 0.5 * wantWidth * sf,
			wantYmax: 0,
		},
		{
			name:     "XCenter,YTop w/ offset",
			x:        10,
			y:        20,
			opts:     TopCenter,
			wantXmin: (10 - 0.5*wantWidth) * sf,
			wantYmin: (20 - wantHeight) * sf,
			wantXmax: (10 + 0.5*wantWidth) * sf,
			wantYmax: 20 * sf,
		},
		{
			name:     "XRight,YTop",
			opts:     TopRight,
			wantXmin: -wantWidth * sf,
			wantYmin: -wantHeight * sf,
			wantXmax: 0,
			wantYmax: 0,
		},
		{
			name:     "XRight,YTop w/ offset",
			x:        10,
			y:        20,
			opts:     TopRight,
			wantXmin: (10 - wantWidth) * sf,
			wantYmin: (20 - wantHeight) * sf,
			wantXmax: 10 * sf,
			wantYmax: 20 * sf,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := Text(tt.x, tt.y, 1, message, fontName, pts, &tt.opts)
			gotWidth := text.Width()
			if math.Abs(gotWidth-wantWidth) > eps {
				t.Errorf("width = %v, want %v", gotWidth, wantWidth)
			}
			gotHeight := text.Height()
			if math.Abs(gotHeight-wantHeight) > eps {
				t.Errorf("height = %v, want %v", gotHeight, wantHeight)
			}
			mbb := text.MBB()
			if math.Abs(mbb.Min[0]-tt.wantXmin) > eps {
				t.Errorf("Xmin = %v, want %v", mbb.Min[0], tt.wantXmin)
			}
			if math.Abs(mbb.Min[1]-tt.wantYmin) > eps {
				t.Errorf("Ymin = %v, want %v", mbb.Min[1], tt.wantYmin)
			}
			if math.Abs(mbb.Max[0]-tt.wantXmax) > eps {
				t.Errorf("Xmax = %v, want %v", mbb.Max[0], tt.wantXmax)
			}
			if math.Abs(mbb.Max[1]-tt.wantYmax) > eps {
				t.Errorf("Ymax = %v, want %v", mbb.Max[1], tt.wantYmax)
			}
		})
	}
}
