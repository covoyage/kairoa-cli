package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var chmodCmd = &cobra.Command{
	Use:   "chmod",
	Short: "File permission calculator",
	Long:  `Calculate and convert file permissions between octal and symbolic notation.`,
}

var chmodCalcCmd = &cobra.Command{
	Use:   "calc",
	Short: "Calculate permissions from options",
	RunE: func(cmd *cobra.Command, args []string) error {
		ownerRead, _ := cmd.Flags().GetBool("owner-read")
		ownerWrite, _ := cmd.Flags().GetBool("owner-write")
		ownerExec, _ := cmd.Flags().GetBool("owner-exec")
		groupRead, _ := cmd.Flags().GetBool("group-read")
		groupWrite, _ := cmd.Flags().GetBool("group-write")
		groupExec, _ := cmd.Flags().GetBool("group-exec")
		otherRead, _ := cmd.Flags().GetBool("other-read")
		otherWrite, _ := cmd.Flags().GetBool("other-write")
		otherExec, _ := cmd.Flags().GetBool("other-exec")

		owner := boolToInt(ownerRead)*4 + boolToInt(ownerWrite)*2 + boolToInt(ownerExec)
		group := boolToInt(groupRead)*4 + boolToInt(groupWrite)*2 + boolToInt(groupExec)
		other := boolToInt(otherRead)*4 + boolToInt(otherWrite)*2 + boolToInt(otherExec)

		octal := fmt.Sprintf("%d%d%d", owner, group, other)

		// Build symbolic notation
		var symbolic strings.Builder
		symbolic.WriteString(getPermissionString(ownerRead, ownerWrite, ownerExec))
		symbolic.WriteString(getPermissionString(groupRead, groupWrite, groupExec))
		symbolic.WriteString(getPermissionString(otherRead, otherWrite, otherExec))

		fmt.Printf("Octal: %s\n", octal)
		fmt.Printf("Symbolic: %s\n", symbolic.String())
		fmt.Printf("Command: chmod %s file\n", octal)

		return nil
	},
}

var chmodOctalCmd = &cobra.Command{
	Use:   "octal [mode]",
	Short: "Convert octal mode to symbolic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := args[0]

		if len(mode) != 3 {
			return fmt.Errorf("octal mode must be 3 digits (e.g., 755)")
		}

		octal, err := strconv.ParseInt(mode, 8, 64)
		if err != nil {
			return fmt.Errorf("invalid octal mode: %s", mode)
		}

		owner := (octal >> 6) & 7
		group := (octal >> 3) & 7
		other := octal & 7

		fmt.Printf("Octal: %s\n", mode)
		fmt.Printf("Symbolic: %s%s%s\n",
			octalToSymbolic(int(owner)),
			octalToSymbolic(int(group)),
			octalToSymbolic(int(other)))

		fmt.Println("\nBreakdown:")
		fmt.Printf("  Owner (%d): %s\n", owner, octalToSymbolic(int(owner)))
		fmt.Printf("  Group (%d): %s\n", group, octalToSymbolic(int(group)))
		fmt.Printf("  Other (%d): %s\n", other, octalToSymbolic(int(other)))

		return nil
	},
}

var chmodSymbolicCmd = &cobra.Command{
	Use:   "symbolic [mode]",
	Short: "Convert symbolic mode to octal",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := args[0]

		if len(mode) != 9 && len(mode) != 3 {
			return fmt.Errorf("symbolic mode must be 9 characters (e.g., rwxr-xr-x) or 3 (e.g., rwx)")
		}

		if len(mode) == 3 {
			mode = mode + "---" + "---"
		}

		owner := symbolicToOctal(mode[0:3])
		group := symbolicToOctal(mode[3:6])
		other := symbolicToOctal(mode[6:9])

		octal := fmt.Sprintf("%d%d%d", owner, group, other)

		fmt.Printf("Symbolic: %s\n", mode)
		fmt.Printf("Octal: %s\n", octal)

		return nil
	},
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func getPermissionString(read, write, exec bool) string {
	var result strings.Builder
	if read {
		result.WriteString("r")
	} else {
		result.WriteString("-")
	}
	if write {
		result.WriteString("w")
	} else {
		result.WriteString("-")
	}
	if exec {
		result.WriteString("x")
	} else {
		result.WriteString("-")
	}
	return result.String()
}

func octalToSymbolic(octal int) string {
	read := octal&4 != 0
	write := octal&2 != 0
	exec := octal&1 != 0
	return getPermissionString(read, write, exec)
}

func symbolicToOctal(sym string) int {
	result := 0
	if len(sym) >= 1 && sym[0] == 'r' {
		result += 4
	}
	if len(sym) >= 2 && sym[1] == 'w' {
		result += 2
	}
	if len(sym) >= 3 && sym[2] == 'x' {
		result += 1
	}
	return result
}

func init() {
	rootCmd.AddCommand(chmodCmd)
	chmodCmd.AddCommand(chmodCalcCmd)
	chmodCmd.AddCommand(chmodOctalCmd)
	chmodCmd.AddCommand(chmodSymbolicCmd)

	chmodCalcCmd.Flags().BoolP("owner-read", "r", false, "Owner read permission")
	chmodCalcCmd.Flags().BoolP("owner-write", "w", false, "Owner write permission")
	chmodCalcCmd.Flags().BoolP("owner-exec", "x", false, "Owner execute permission")
	chmodCalcCmd.Flags().Bool("group-read", false, "Group read permission")
	chmodCalcCmd.Flags().Bool("group-write", false, "Group write permission")
	chmodCalcCmd.Flags().Bool("group-exec", false, "Group execute permission")
	chmodCalcCmd.Flags().Bool("other-read", false, "Other read permission")
	chmodCalcCmd.Flags().Bool("other-write", false, "Other write permission")
	chmodCalcCmd.Flags().Bool("other-exec", false, "Other execute permission")
}
