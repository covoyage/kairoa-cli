package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker command generator",
	Long:  `Generate Docker commands for common operations.`,
}

var dockerRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Generate docker run command",
	RunE: func(cmd *cobra.Command, args []string) error {
		image, _ := cmd.Flags().GetString("image")
		name, _ := cmd.Flags().GetString("name")
		detached, _ := cmd.Flags().GetBool("detached")
		interactive, _ := cmd.Flags().GetBool("interactive")
		tty, _ := cmd.Flags().GetBool("tty")
		ports, _ := cmd.Flags().GetStringSlice("port")
		volumes, _ := cmd.Flags().GetStringSlice("volume")
		envs, _ := cmd.Flags().GetStringSlice("env")
		remove, _ := cmd.Flags().GetBool("remove")
		restart, _ := cmd.Flags().GetString("restart")

		var cmdParts []string
		cmdParts = append(cmdParts, "docker run")

		if detached {
			cmdParts = append(cmdParts, "-d")
		}
		if interactive {
			cmdParts = append(cmdParts, "-i")
		}
		if tty {
			cmdParts = append(cmdParts, "-t")
		}
		if remove {
			cmdParts = append(cmdParts, "--rm")
		}
		if name != "" {
			cmdParts = append(cmdParts, fmt.Sprintf("--name %s", name))
		}
		if restart != "" {
			cmdParts = append(cmdParts, fmt.Sprintf("--restart %s", restart))
		}
		for _, p := range ports {
			cmdParts = append(cmdParts, fmt.Sprintf("-p %s", p))
		}
		for _, v := range volumes {
			cmdParts = append(cmdParts, fmt.Sprintf("-v %s", v))
		}
		for _, e := range envs {
			cmdParts = append(cmdParts, fmt.Sprintf("-e %s", e))
		}
		cmdParts = append(cmdParts, image)

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

var dockerBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Generate docker build command",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		tag, _ := cmd.Flags().GetString("tag")
		file, _ := cmd.Flags().GetString("file")
		noCache, _ := cmd.Flags().GetBool("no-cache")
		pull, _ := cmd.Flags().GetBool("pull")

		var cmdParts []string
		cmdParts = append(cmdParts, "docker build")

		if tag != "" {
			cmdParts = append(cmdParts, fmt.Sprintf("-t %s", tag))
		}
		if file != "" {
			cmdParts = append(cmdParts, fmt.Sprintf("-f %s", file))
		}
		if noCache {
			cmdParts = append(cmdParts, "--no-cache")
		}
		if pull {
			cmdParts = append(cmdParts, "--pull")
		}
		cmdParts = append(cmdParts, path)

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

var dockerComposeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Generate docker compose commands",
}

var dockerComposeUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Generate docker compose up command",
	RunE: func(cmd *cobra.Command, args []string) error {
		detached, _ := cmd.Flags().GetBool("detached")
		build, _ := cmd.Flags().GetBool("build")

		var cmdParts []string
		cmdParts = append(cmdParts, "docker compose up")

		if detached {
			cmdParts = append(cmdParts, "-d")
		}
		if build {
			cmdParts = append(cmdParts, "--build")
		}

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

var dockerPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "Generate docker ps command",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		if all {
			fmt.Println("docker ps -a")
		} else {
			fmt.Println("docker ps")
		}
		return nil
	},
}

var dockerLogsCmd = &cobra.Command{
	Use:   "logs [container]",
	Short: "Generate docker logs command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		container := args[0]
		follow, _ := cmd.Flags().GetBool("follow")
		tail, _ := cmd.Flags().GetString("tail")

		var cmdParts []string
		cmdParts = append(cmdParts, "docker logs")

		if follow {
			cmdParts = append(cmdParts, "-f")
		}
		if tail != "" {
			cmdParts = append(cmdParts, fmt.Sprintf("--tail %s", tail))
		}
		cmdParts = append(cmdParts, container)

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
	dockerCmd.AddCommand(dockerRunCmd)
	dockerCmd.AddCommand(dockerBuildCmd)
	dockerCmd.AddCommand(dockerComposeCmd)
	dockerCmd.AddCommand(dockerPsCmd)
	dockerCmd.AddCommand(dockerLogsCmd)

	dockerRunCmd.Flags().StringP("image", "i", "", "Image name (required)")
	dockerRunCmd.Flags().StringP("name", "n", "", "Container name")
	dockerRunCmd.Flags().BoolP("detached", "d", false, "Run in detached mode")
	dockerRunCmd.Flags().Bool("interactive", false, "Keep STDIN open")
	dockerRunCmd.Flags().BoolP("tty", "t", false, "Allocate a pseudo-TTY")
	dockerRunCmd.Flags().StringSliceP("port", "p", []string{}, "Port mappings (host:container)")
	dockerRunCmd.Flags().StringSliceP("volume", "v", []string{}, "Volume mounts (host:container)")
	dockerRunCmd.Flags().StringSliceP("env", "e", []string{}, "Environment variables")
	dockerRunCmd.Flags().BoolP("remove", "r", false, "Auto-remove container when it exits")
	dockerRunCmd.Flags().String("restart", "", "Restart policy (no, on-failure, always, unless-stopped)")

	dockerBuildCmd.Flags().StringP("path", "p", ".", "Build path")
	dockerBuildCmd.Flags().StringP("tag", "t", "", "Image tag")
	dockerBuildCmd.Flags().StringP("file", "f", "", "Dockerfile path")
	dockerBuildCmd.Flags().Bool("no-cache", false, "Do not use cache")
	dockerBuildCmd.Flags().Bool("pull", false, "Always pull newer images")

	dockerComposeCmd.AddCommand(dockerComposeUpCmd)
	dockerComposeUpCmd.Flags().BoolP("detached", "d", false, "Run in detached mode")
	dockerComposeUpCmd.Flags().BoolP("build", "b", false, "Build images before starting")

	dockerPsCmd.Flags().BoolP("all", "a", false, "Show all containers")

	dockerLogsCmd.Flags().BoolP("follow", "f", false, "Follow log output")
	dockerLogsCmd.Flags().String("tail", "", "Number of lines to show from the end")
}
