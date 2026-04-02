package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var coordinateCmd = &cobra.Command{
	Use:   "coordinate",
	Short: "Coordinate converter",
	Long:  `Convert between different coordinate formats (decimal degrees, DMS).`,
}

var coordinateToDMSCmd = &cobra.Command{
	Use:   "to-dms [decimal]",
	Short: "Convert decimal degrees to DMS",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		decimal, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return fmt.Errorf("invalid decimal value: %s", args[0])
		}

		dms := decimalToDMS(decimal)
		fmt.Printf("Decimal: %.6f°\n", decimal)
		fmt.Printf("DMS: %s\n", dms)

		return nil
	},
}

var coordinateToDecimalCmd = &cobra.Command{
	Use:   "to-decimal [dms]",
	Short: "Convert DMS to decimal degrees",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		decimal, err := dmsToDecimal(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("DMS: %s\n", args[0])
		fmt.Printf("Decimal: %.6f°\n", decimal)

		return nil
	},
}

var coordinateDistanceCmd = &cobra.Command{
	Use:   "distance [lat1] [lon1] [lat2] [lon2]",
	Short: "Calculate distance between two coordinates",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		lat1, _ := strconv.ParseFloat(args[0], 64)
		lon1, _ := strconv.ParseFloat(args[1], 64)
		lat2, _ := strconv.ParseFloat(args[2], 64)
		lon2, _ := strconv.ParseFloat(args[3], 64)

		distance := haversine(lat1, lon1, lat2, lon2)

		fmt.Printf("Point 1: %.6f°, %.6f°\n", lat1, lon1)
		fmt.Printf("Point 2: %.6f°, %.6f°\n", lat2, lon2)
		fmt.Printf("Distance: %.2f km\n", distance)
		fmt.Printf("Distance: %.2f miles\n", distance*0.621371)

		return nil
	},
}

func decimalToDMS(decimal float64) string {
	isNegative := decimal < 0
	if isNegative {
		decimal = -decimal
	}

	degrees := math.Floor(decimal)
	minutesFloat := (decimal - degrees) * 60
	minutes := math.Floor(minutesFloat)
	seconds := (minutesFloat - minutes) * 60

	sign := ""
	if isNegative {
		sign = "-"
	}

	return fmt.Sprintf("%s%.0f° %.0f' %.4f\"", sign, degrees, minutes, seconds)
}

func dmsToDecimal(dms string) (float64, error) {
	// Parse DMS format like "40° 26' 46.12\" N" or "40:26:46.12"
	dms = strings.TrimSpace(dms)

	// Try to extract numbers
	var degrees, minutes, seconds float64
	var direction string

	// Check for direction at end
	if strings.HasSuffix(dms, "N") || strings.HasSuffix(dms, "E") {
		direction = "+"
		dms = dms[:len(dms)-1]
	} else if strings.HasSuffix(dms, "S") || strings.HasSuffix(dms, "W") {
		direction = "-"
		dms = dms[:len(dms)-1]
	}

	// Replace common separators with spaces
	dms = strings.ReplaceAll(dms, "°", " ")
	dms = strings.ReplaceAll(dms, "'", " ")
	dms = strings.ReplaceAll(dms, "\"", " ")
	dms = strings.ReplaceAll(dms, ":", " ")
	dms = strings.TrimSpace(dms)

	parts := strings.Fields(dms)
	if len(parts) < 1 {
		return 0, fmt.Errorf("invalid DMS format")
	}

	degrees, _ = strconv.ParseFloat(parts[0], 64)
	if len(parts) > 1 {
		minutes, _ = strconv.ParseFloat(parts[1], 64)
	}
	if len(parts) > 2 {
		seconds, _ = strconv.ParseFloat(parts[2], 64)
	}

	decimal := degrees + minutes/60 + seconds/3600

	if direction == "-" {
		decimal = -decimal
	}

	return decimal, nil
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers

	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func init() {
	rootCmd.AddCommand(coordinateCmd)
	coordinateCmd.AddCommand(coordinateToDMSCmd)
	coordinateCmd.AddCommand(coordinateToDecimalCmd)
	coordinateCmd.AddCommand(coordinateDistanceCmd)
}
