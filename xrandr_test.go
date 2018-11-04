package xrandr_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vcraescu/go-xrandr"
	"testing"
)

func TestParseSize(t *testing.T) {
	r, err := xrandr.ParseSize("3000x1200")
	assert.NoError(t, err)

	assert.Equal(t, float32(3000), r.Width)
	assert.Equal(t, float32(1200), r.Height)

	r, err = xrandr.ParseSize("\n3000 x 1200   \t")
	assert.NoError(t, err)
	assert.Equal(t, float32(3000), r.Width)
	assert.Equal(t, float32(1200), r.Height)

	r, err = xrandr.ParseSize("\n3000 1200   \t")
	assert.Error(t, err)
}

func TestParseSizeWithPosition(t *testing.T) {
	s, p, err := xrandr.ParseSizeWithPosition("3000x1200+1200+400")
	assert.NoError(t, err)

	assert.Equal(t, float32(3000), s.Width)
	assert.Equal(t, float32(1200), s.Height)
	assert.Equal(t, 1200, p.X)
	assert.Equal(t, 400, p.Y)

	_, _, err = xrandr.ParseSizeWithPosition("3000x1200 + 1200 +   400")
	assert.NoError(t, err)
}

func TestParseRefreshRate(t *testing.T) {
	r, err := xrandr.ParseRate("  30.12 \t \n")
	assert.NoError(t, err)
	assert.Equal(t, xrandr.RefreshRateValue(30.12), r.Value)
	assert.False(t, r.Current)
	assert.False(t, r.Preferred)

	r, err = xrandr.ParseRate("  30.12+ ")
	assert.NoError(t, err)
	assert.Equal(t, xrandr.RefreshRateValue(30.12), r.Value)
	assert.False(t, r.Current)
	assert.True(t, r.Preferred)

	r, err = xrandr.ParseRate("  30.12* ")
	assert.NoError(t, err)
	assert.Equal(t, xrandr.RefreshRateValue(30.12), r.Value)
	assert.True(t, r.Current)
	assert.False(t, r.Preferred)

	r, err = xrandr.ParseRate("  30.12+* ")
	assert.NoError(t, err)
	assert.Equal(t, xrandr.RefreshRateValue(30.12), r.Value)
	assert.True(t, r.Current)
	assert.True(t, r.Preferred)

	r, err = xrandr.ParseRate("  30.12 + * ")
	assert.NoError(t, err)
	assert.Equal(t, xrandr.RefreshRateValue(30.12), r.Value)
	assert.True(t, r.Current)
	assert.True(t, r.Preferred)

	r, err = xrandr.ParseRate("60.05")
	assert.NoError(t, err)
	assert.Equal(t, xrandr.RefreshRateValue(60.05), r.Value)
	assert.False(t, r.Current)
	assert.False(t, r.Preferred)

	r, err = xrandr.ParseRate("23.98+")
	assert.NoError(t, err)
	assert.Equal(t, xrandr.RefreshRateValue(23.98), r.Value)
	assert.False(t, r.Current)
	assert.True(t, r.Preferred)
}

func TestParseModeLine(t *testing.T) {
	m, err := xrandr.ParseModeLine(" \t 1920x1080     60.00    59.94    50.00    23.98    60.05    60.00    50.04 \t")
	assert.NoError(t, err)
	assert.Equal(t, float32(1920), m.Resolution.Width)
	assert.Equal(t, float32(1080), m.Resolution.Height)
	assert.Len(t, m.RefreshRates, 7)
	assert.Equal(t, xrandr.RefreshRateValue(60), m.RefreshRates[0].Value)
	assert.Equal(t, xrandr.RefreshRateValue(59.94), m.RefreshRates[1].Value)
	assert.Equal(t, xrandr.RefreshRateValue(50.00), m.RefreshRates[2].Value)
	assert.Equal(t, xrandr.RefreshRateValue(23.98), m.RefreshRates[3].Value)
	assert.Equal(t, xrandr.RefreshRateValue(60.05), m.RefreshRates[4].Value)
	assert.Equal(t, xrandr.RefreshRateValue(60.00), m.RefreshRates[5].Value)
	assert.Equal(t, xrandr.RefreshRateValue(50.04), m.RefreshRates[6].Value)

	m, err = xrandr.ParseModeLine(" \t 1920x1080     60.00    59.94*    50.00    23.98+    60.05    60.00    50.04 \t")
	assert.NoError(t, err)
	assert.Equal(t, float32(1920), m.Resolution.Width)
	assert.Equal(t, float32(1080), m.Resolution.Height)
	assert.Equal(t, xrandr.RefreshRateValue(59.94), m.RefreshRates[1].Value)
	assert.True(t, m.RefreshRates[1].Current)

	assert.Equal(t, xrandr.RefreshRateValue(23.98), m.RefreshRates[3].Value)
	assert.True(t, m.RefreshRates[3].Preferred)
}

func TestParseScreenLine(t *testing.T) {
	s, err := xrandr.ParseScreenLine(" Screen 3: minimum 8 x 8, current 9600 x 3240, maximum 32767 x 32767 ")
	assert.NoError(t, err)
	assert.Equal(t, 3, s.No)
	assert.Equal(t, float32(8), s.MinResolution.Width)
	assert.Equal(t, float32(8), s.MinResolution.Height)

	assert.Equal(t, float32(9600), s.CurrentResolution.Width)
	assert.Equal(t, float32(3240), s.CurrentResolution.Height)

	assert.Equal(t, float32(32767), s.MaxResolution.Width)
	assert.Equal(t, float32(32767), s.MaxResolution.Height)
}

func TestParseMonitorLine(t *testing.T) {
	m, err := xrandr.ParseMonitorLine(" HDMI-0 connected primary 5760x3240+1000+500 (normal left inverted right x axis y axis) 597mm x 336mm")
	assert.NoError(t, err)
	assert.Equal(t, "HDMI-0", m.ID)
	assert.True(t, m.Connected)
	assert.True(t, m.Primary)
	assert.Equal(t, float32(597), m.Size.Width)
	assert.Equal(t, float32(336), m.Size.Height)
	assert.Equal(t, float32(5760), m.Resolution.Width)
	assert.Equal(t, float32(3240), m.Resolution.Height)
	assert.Equal(t, 1000, m.Position.X)
	assert.Equal(t, 500, m.Position.Y)

	m, err = xrandr.ParseMonitorLine("DP-1-1 disconnected (normal left inverted right x axis y axis)")
	assert.NoError(t, err)
	assert.Equal(t, "DP-1-1", m.ID)
}

func TestIsScreenLine(t *testing.T) {
	b := xrandr.IsScreenLine("Screen 0: minimum 8 x 8, current 9408 x 3132, maximum 32767 x 32767")
	assert.True(t, b)

	b = xrandr.IsScreenLine(" minimum 8 x 8, Screen 0, current 9408 x 3132, maximum 32767 x 32767")
	assert.False(t, b)

	b = xrandr.IsScreenLine(" minimum 8 x 8, current 9408 x 3132, maximum 32767 x 32767")
	assert.False(t, b)
}

func TestIsMonitorLine(t *testing.T) {
	b := xrandr.IsMonitorLine("HDMI-0 connected primary 5568x3132+0+0 (normal left inverted right x axis y axis) 597mm x 336mm")
	assert.True(t, b)

	b = xrandr.IsMonitorLine("HDMI-0 disconnected primary 5568x3132+0+0 (normal left inverted right x axis y axis) 597mm x 336mm")
	assert.True(t, b)

	b = xrandr.IsMonitorLine(" minimum 8 x 8, current 9408 x 3132, maximum 32767 x 32767")
	assert.False(t, b)
}

func TestParseScreens(t *testing.T) {
	screens, err := xrandr.ParseScreens(`Screen 0: minimum 8 x 8, current 9408 x 3132, maximum 32767 x 32767
HDMI-0 connected primary 5568x3132+0+0 (normal left inverted right x axis y axis) 597mm x 336mm
   3840x2160     30.00*+  29.97    25.00    23.98  
   1920x1200     59.88  
   1920x1080     60.00    59.94    50.00    23.98    60.05    60.00    50.04  
   1680x1050     59.95  
   1600x1200     60.00  
   1280x1024     75.02    60.02  
   1280x800      59.81  
   1280x720      60.00    59.94    50.00  
   1152x864      75.00  
   1024x768      75.03    60.00  
   800x600       75.00    60.32  
   720x576       50.00  
   720x480       59.94  
   640x480       75.00    59.94    59.93  
HDMI-1 disconnected primary 5568x3132+0+0 (normal left inverted right x axis y axis) 597mm x 336mm
eDP-1-1 connected 3840x2160+5568+0 (normal left inverted right x axis y axis) 346mm x 194mm
   3840x2160     60.00*+  59.98    48.02    59.97  
   3200x1800     59.96    59.94  
   2880x1620     59.96    59.97  
   2560x1600     59.99    59.97  
   2560x1440     59.99    59.99    59.96    59.95  
   2048x1536     60.00  
   1920x1440     60.00  
   1856x1392     60.01  
   1792x1344     60.01  
   2048x1152     59.99    59.98    59.90    59.91  
   1920x1200     59.88    59.95  
   1920x1080     60.01    59.97    59.96    59.93  
DP-1-1 disconnected (normal left inverted right x axis y axis)
  1920x1200 (0x5b) 193.250MHz -HSync +VSync
        h: width  1920 start 2056 end 2256 total 2592 skew    0 clock  74.56KHz
        v: height 1200 start 1203 end 1209 total 1245           clock  59.88Hz
  1600x1200 (0x61) 162.000MHz +HSync +VSync
        h: width  1600 start 1664 end 1856 total 2160 skew    0 clock  75.00KHz
        v: height 1200 start 1201 end 1204 total 1250           clock  60.00Hz
  1680x1050 (0x62) 146.250MHz -HSync +VSync
        h: width  1680 start 1784 end 1960 total 2240 skew    0 clock  65.29KHz
        v: height 1050 start 1053 end 1059 total 1089           clock  59.95Hz
  1280x1024 (0x69) 108.000MHz +HSync +VSync
        h: width  1280 start 1328 end 1440 total 1688 skew    0 clock  63.98KHz
        v: height 1024 start 1025 end 1028 total 1066           clock  60.02Hz
  1280x800 (0x73) 83.500MHz -HSync +VSync
        h: width  1280 start 1352 end 1480 total 1680 skew    0 clock  49.70KHz
        v: height  800 start  803 end  809 total  831           clock  59.81Hz
  1024x768 (0x7a) 65.000MHz -HSync -VSync
        h: width  1024 start 1048 end 1184 total 1344 skew    0 clock  48.36KHz
        v: height  768 start  771 end  777 total  806           clock  60.00Hz
  800x600 (0x89) 40.000MHz +HSync +VSync
        h: width   800 start  840 end  968 total 1056 skew    0 clock  37.88KHz
        v: height  600 start  601 end  605 total  628           clock  60.32Hz
  640x480 (0x96) 25.175MHz -HSync -VSync
        h: width   640 start  656 end  752 total  800 skew    0 clock  31.47KHz
        v: height  480 start  490 end  492 total  525           clock  59.94Hz
eDP-2-2 connected 3840x2160+5568+0 (normal left inverted right x axis y axis) 346mm x 194mm
   2048x1152     59.99    59.98    59.90    59.91  
   1920x1200     59.88    59.95  
   1920x1080     60.01    59.97    59.96    59.93  
`)

	assert.NoError(t, err)
	assert.Len(t, screens, 1)
	screen := screens[0]
	assert.Equal(t, 0, screen.No)
	assert.Equal(t, xrandr.Size{
		Width:  8,
		Height: 8,
	}, screens[0].MinResolution)
	assert.Equal(t, xrandr.Size{
		Width:  9408,
		Height: 3132,
	}, screen.CurrentResolution)
	assert.Equal(t, xrandr.Size{
		Width:  32767,
		Height: 32767,
	}, screen.MaxResolution)
	assert.Len(t, screen.Monitors, 5)

	monitor := screen.Monitors[0]
	assert.Equal(t, "HDMI-0", monitor.ID)
	assert.True(t, monitor.Connected)
	assert.True(t, monitor.Primary)

	assert.Equal(t, xrandr.Size{
		Width:  5568,
		Height: 3132,
	}, monitor.Resolution)
	assert.Equal(t, xrandr.Position{
		X: 0,
		Y: 0,
	}, monitor.Position)
	assert.Equal(t, xrandr.Size{
		Width:  597,
		Height: 336,
	}, monitor.Size)
	assert.Len(t, monitor.Modes, 14)
	assert.Equal(t, xrandr.Size{
		Width:  3840,
		Height: 2160,
	}, monitor.Modes[0].Resolution)
	assert.Equal(t, xrandr.RefreshRate{
		Value:     30,
		Current:   true,
		Preferred: true,
	}, monitor.Modes[0].RefreshRates[0])
	assert.Equal(t, xrandr.RefreshRate{
		Value:     23.98,
		Current:   false,
		Preferred: false,
	}, monitor.Modes[0].RefreshRates[3])

	assert.Equal(t, xrandr.Size{
		Width:  640,
		Height: 480,
	}, monitor.Modes[13].Resolution)
	assert.Equal(t, xrandr.RefreshRate{
		Value:     75,
		Current:   false,
		Preferred: false,
	}, monitor.Modes[13].RefreshRates[0])
	assert.Equal(t, xrandr.RefreshRate{
		Value:     59.93,
		Current:   false,
		Preferred: false,
	}, monitor.Modes[13].RefreshRates[2])

	monitor = screen.Monitors[1]
	assert.Equal(t, "HDMI-1", monitor.ID)
	assert.False(t, monitor.Connected)

	monitor = screen.Monitors[2]
	assert.Equal(t, "eDP-1-1", monitor.ID)
	assert.True(t, monitor.Connected)
	assert.False(t, monitor.Primary)
	assert.Equal(t, xrandr.Size{
		Width:  3840,
		Height: 2160,
	}, monitor.Resolution)
	assert.Equal(t, xrandr.Position{
		X: 5568,
		Y: 0,
	}, monitor.Position)
	assert.Equal(t, xrandr.Size{
		Width:  346,
		Height: 194,
	}, monitor.Size)
	assert.Len(t, monitor.Modes, 12)

	assert.Equal(t, xrandr.Size{
		Width:  3840,
		Height: 2160,
	}, monitor.Modes[0].Resolution)
	assert.Equal(t, xrandr.RefreshRate{
		Value:     60,
		Current:   true,
		Preferred: true,
	}, monitor.Modes[0].RefreshRates[0])

	monitor = screen.Monitors[3]
	assert.Equal(t, "DP-1-1", monitor.ID)
	assert.False(t, monitor.Connected)
	assert.False(t, monitor.Primary)
	assert.Len(t, monitor.Modes, 0)

	monitor = screen.Monitors[4]
	assert.Equal(t, "eDP-2-2", monitor.ID)
	assert.True(t, monitor.Connected)
	assert.False(t, monitor.Primary)
	assert.Len(t, monitor.Modes, 3)
}

func TestScreen_GeMonitorByID(t *testing.T) {
	monitors := []xrandr.Monitor{
		{ID: "Monitor1"},
		{ID: "Monitor2"},
		{ID: "Monitor3"},
	}

	screen := xrandr.Screen{Monitors: monitors}
	m, ok := screen.MonitorByID("Monitor2")
	assert.True(t, ok)
	assert.Equal(t, "Monitor2", m.ID)
}

func TestScreens_GeMonitorByID(t *testing.T) {
	monitors := []xrandr.Monitor{
		{ID: "Monitor1"},
		{ID: "Monitor2"},
		{ID: "Monitor3"},
	}

	screen := xrandr.Screen{Monitors: monitors}
	screens := xrandr.Screens{screen}
	m, ok := screens.MonitorByID("Monitor3")
	assert.True(t, ok)
	assert.Equal(t, "Monitor3", m.ID)
}

func TestMode_CurrentRefreshRate(t *testing.T) {
	refreshRates := []xrandr.RefreshRate{
		{Current: false, Value: xrandr.RefreshRateValue(60.0)},
		{Current: false, Value: xrandr.RefreshRateValue(30.0)},
		{Current: true, Value: xrandr.RefreshRateValue(144)},
	}

	mode := xrandr.Mode{RefreshRates: refreshRates}

	crr, ok := mode.CurrentRefreshRate()
	assert.True(t, ok)
	assert.Equal(t, xrandr.RefreshRateValue(144.0), crr.Value)
}

func TestMonitor_CurrentMode(t *testing.T) {
	refreshRates1 := []xrandr.RefreshRate{
		{Current: false, Value: xrandr.RefreshRateValue(60)},
		{Current: false, Value: xrandr.RefreshRateValue(30)},
		{Current: false, Value: xrandr.RefreshRateValue(144)},
	}
	refreshRates2 := []xrandr.RefreshRate{
		{Current: true, Value: xrandr.RefreshRateValue(30)},
		{Current: false, Value: xrandr.RefreshRateValue(54)},
	}

	monitor := xrandr.Monitor{
		Modes: []xrandr.Mode{
			{RefreshRates: refreshRates1},
			{RefreshRates: refreshRates2},
		},
	}

	_, ok := monitor.CurrentMode()
	assert.True(t, ok)
}

func TestMonitor_DPI(t *testing.T) {
	monitor := xrandr.Monitor{
		Modes: []xrandr.Mode{
			{
				RefreshRates: []xrandr.RefreshRate{
					{Current: true, Value: xrandr.RefreshRateValue(30)},
					{Current: false, Value: xrandr.RefreshRateValue(54)},
				},
				Resolution: xrandr.Size{
					Width: 3840,
					Height: 2160,
				},
			},
		},

		Size: xrandr.Size{
			Width: 487,
			Height: 247,
		},
	}

	dpi, err := monitor.DPI()
	assert.Nil(t, err)
	assert.Equal(t, 200, int(dpi))
}

func TestSize_Rescale(t *testing.T) {
	s := xrandr.Size{Width: 1000, Height: 2000}
	s = s.Rescale(1.333333)

	assert.Equal(t, float32(1334), s.Width)
	assert.Equal(t, float32(2667), s.Height)
}
