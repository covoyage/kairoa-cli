package cmd

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var colorCmd = &cobra.Command{
	Use:   "color",
	Short: "Color conversion utilities",
	Long:  `Convert between different color formats (HEX, RGB, HSL, CMYK, etc.).`,
}

var colorConvertCmd = &cobra.Command{
	Use:   "convert [color]",
	Short: "Convert color to all formats",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.TrimSpace(args[0])

		var r, g, b int
		var err error

		// Try to parse as HEX
		if strings.HasPrefix(input, "#") || len(input) == 6 || len(input) == 3 {
			r, g, b, err = parseHex(input)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(strings.ToLower(input), "rgb") {
			// Parse RGB
			r, g, b, err = parseRGB(input)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(strings.ToLower(input), "hsl") {
			// Parse HSL and convert to RGB
			h, s, l, err := parseHSL(input)
			if err != nil {
				return err
			}
			r, g, b = hslToRGB(h, s, l)
		} else {
			return fmt.Errorf("unsupported color format: %s", input)
		}

		// Print color preview
		fmt.Println("Color Preview:")
		fmt.Printf("  \033[48;2;%d;%d;%dm        \033[0m\n\n", r, g, b)

		// Print all formats
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Format", "Value"})

		hex := rgbToHex(r, g, b)
		table.Append([]string{"HEX", hex})
		table.Append([]string{"HEX (short)", rgbToHexShort(r, g, b)})
		table.Append([]string{"RGB", fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)})
		table.Append([]string{"RGB (percent)", fmt.Sprintf("rgb(%d%%, %d%%, %d%%)", r*100/255, g*100/255, b*100/255)})

		h, s, l := rgbToHSL(r, g, b)
		table.Append([]string{"HSL", fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l)})

		c, m, y, k := rgbToCMYK(r, g, b)
		table.Append([]string{"CMYK", fmt.Sprintf("cmyk(%d%%, %d%%, %d%%, %d%%)", c, m, y, k)})

		table.Render()

		return nil
	},
}

var colorRandomCmd = &cobra.Command{
	Use:   "random",
	Short: "Generate a random color",
	RunE: func(cmd *cobra.Command, args []string) error {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(256)
		g := rand.Intn(256)
		b := rand.Intn(256)

		fmt.Printf("Random Color: %s\n", rgbToHex(r, g, b))
		fmt.Printf("\033[48;2;%d;%d;%dm        \033[0m\n", r, g, b)
		fmt.Printf("RGB: rgb(%d, %d, %d)\n", r, g, b)

		return nil
	},
}

func parseHex(hex string) (r, g, b int, err error) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) == 3 {
		// Short form: RGB -> RRGGBB
		hex = string(hex[0]) + string(hex[0]) + string(hex[1]) + string(hex[1]) + string(hex[2]) + string(hex[2])
	}

	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid hex color: %s", hex)
	}

	values, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid hex color: %s", hex)
	}

	r = int(values >> 16)
	g = int((values >> 8) & 0xFF)
	b = int(values & 0xFF)

	return r, g, b, nil
}

func parseRGB(rgb string) (r, g, b int, err error) {
	// Remove "rgb(" and ")"
	rgb = strings.TrimPrefix(strings.ToLower(rgb), "rgb(")
	rgb = strings.TrimSuffix(rgb, ")")

	parts := strings.Split(rgb, ",")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid RGB format")
	}

	r, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
	g, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
	b, _ = strconv.Atoi(strings.TrimSpace(parts[2]))

	return r, g, b, nil
}

func parseHSL(hsl string) (h, s, l int, err error) {
	// Remove "hsl(" and ")"
	hsl = strings.TrimPrefix(strings.ToLower(hsl), "hsl(")
	hsl = strings.TrimSuffix(hsl, ")")

	parts := strings.Split(hsl, ",")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid HSL format")
	}

	h, _ = strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(parts[0], "°")))
	s, _ = strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(parts[1], "%")))
	l, _ = strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(parts[2], "%")))

	return h, s, l, nil
}

func rgbToHex(r, g, b int) string {
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

func rgbToHexShort(r, g, b int) string {
	// Check if can be shortened
	if r%17 == 0 && g%17 == 0 && b%17 == 0 {
		return fmt.Sprintf("#%X%X%X", r/17, g/17, b/17)
	}
	return rgbToHex(r, g, b)
}

func rgbToHSL(r, g, b int) (h, s, l int) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	// Lightness
	lf := (max + min) / 2.0
	l = int(math.Round(lf * 100))

	// Saturation
	if delta == 0 {
		s = 0
	} else {
		s = int(math.Round((delta / (1 - math.Abs(2*lf-1))) * 100))
	}

	// Hue
	if delta == 0 {
		h = 0
	} else if max == rf {
		h = int(math.Round(60 * math.Mod((gf-bf)/delta, 6)))
	} else if max == gf {
		h = int(math.Round(60 * ((bf-rf)/delta + 2)))
	} else {
		h = int(math.Round(60 * ((rf-gf)/delta + 4)))
	}

	if h < 0 {
		h += 360
	}

	return h, s, l
}

func hslToRGB(h, s, l int) (r, g, b int) {
	hf := float64(h) / 360.0
	sf := float64(s) / 100.0
	lf := float64(l) / 100.0

	c := (1 - math.Abs(2*lf-1)) * sf
	x := c * (1 - math.Abs(math.Mod(hf*6, 2)-1))
	m := lf - c/2

	var rf, gf, bf float64

	if hf < 1.0/6.0 {
		rf, gf, bf = c, x, 0
	} else if hf < 2.0/6.0 {
		rf, gf, bf = x, c, 0
	} else if hf < 3.0/6.0 {
		rf, gf, bf = 0, c, x
	} else if hf < 4.0/6.0 {
		rf, gf, bf = 0, x, c
	} else if hf < 5.0/6.0 {
		rf, gf, bf = x, 0, c
	} else {
		rf, gf, bf = c, 0, x
	}

	r = int(math.Round((rf + m) * 255))
	g = int(math.Round((gf + m) * 255))
	b = int(math.Round((bf + m) * 255))

	return r, g, b
}

func rgbToCMYK(r, g, b int) (c, m, y, k int) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	kf := 1 - math.Max(rf, math.Max(gf, bf))

	if kf == 1 {
		return 0, 0, 0, 100
	}

	cf := (1 - rf - kf) / (1 - kf)
	mf := (1 - gf - kf) / (1 - kf)
	yf := (1 - bf - kf) / (1 - kf)

	c = int(math.Round(cf * 100))
	m = int(math.Round(mf * 100))
	y = int(math.Round(yf * 100))
	k = int(math.Round(kf * 100))

	return c, m, y, k
}

func init() {
	rootCmd.AddCommand(colorCmd)
	colorCmd.AddCommand(colorConvertCmd)
	colorCmd.AddCommand(colorRandomCmd)
}
