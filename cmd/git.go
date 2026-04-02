package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git command generator",
	Long:  `Generate Git commands for common operations.`,
}

var gitCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate git commit command",
	RunE: func(cmd *cobra.Command, args []string) error {
		message, _ := cmd.Flags().GetString("message")
		all, _ := cmd.Flags().GetBool("all")
		amend, _ := cmd.Flags().GetBool("amend")

		var cmdParts []string
		cmdParts = append(cmdParts, "git commit")

		if message != "" {
			cmdParts = append(cmdParts, fmt.Sprintf("-m \"%s\"", message))
		}
		if all {
			cmdParts = append(cmdParts, "-a")
		}
		if amend {
			cmdParts = append(cmdParts, "--amend")
		}

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

var gitPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Generate git push command",
	RunE: func(cmd *cobra.Command, args []string) error {
		remote, _ := cmd.Flags().GetString("remote")
		branch, _ := cmd.Flags().GetString("branch")
		force, _ := cmd.Flags().GetBool("force")
		tags, _ := cmd.Flags().GetBool("tags")

		var cmdParts []string
		cmdParts = append(cmdParts, "git push")
		cmdParts = append(cmdParts, remote)
		if branch != "" {
			cmdParts = append(cmdParts, branch)
		}
		if force {
			cmdParts = append(cmdParts, "--force")
		}
		if tags {
			cmdParts = append(cmdParts, "--tags")
		}

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

var gitPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Generate git pull command",
	RunE: func(cmd *cobra.Command, args []string) error {
		remote, _ := cmd.Flags().GetString("remote")
		branch, _ := cmd.Flags().GetString("branch")
		rebase, _ := cmd.Flags().GetBool("rebase")

		var cmdParts []string
		cmdParts = append(cmdParts, "git pull")
		cmdParts = append(cmdParts, remote)
		if branch != "" {
			cmdParts = append(cmdParts, branch)
		}
		if rebase {
			cmdParts = append(cmdParts, "--rebase")
		}

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

var gitBranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Generate git branch commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		operation, _ := cmd.Flags().GetString("operation")
		name, _ := cmd.Flags().GetString("name")
		oldName, _ := cmd.Flags().GetString("old-name")

		switch operation {
		case "create":
			fmt.Printf("git branch %s\n", name)
		case "delete":
			fmt.Printf("git branch -d %s\n", name)
		case "rename":
			fmt.Printf("git branch -m %s %s\n", oldName, name)
		case "list":
			fmt.Println("git branch -a")
		case "checkout":
			fmt.Printf("git checkout -b %s\n", name)
		default:
			fmt.Println("git branch")
		}
		return nil
	},
}

var gitMergeCmd = &cobra.Command{
	Use:   "merge [branch]",
	Short: "Generate git merge command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]
		noFF, _ := cmd.Flags().GetBool("no-ff")

		var cmdParts []string
		cmdParts = append(cmdParts, "git merge")
		if noFF {
			cmdParts = append(cmdParts, "--no-ff")
		}
		cmdParts = append(cmdParts, branch)

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

var gitRebaseCmd = &cobra.Command{
	Use:   "rebase [branch]",
	Short: "Generate git rebase command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]
		interactive, _ := cmd.Flags().GetBool("interactive")

		var cmdParts []string
		cmdParts = append(cmdParts, "git rebase")
		if interactive {
			cmdParts = append(cmdParts, "-i")
		}
		cmdParts = append(cmdParts, branch)

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

var gitCloneCmd = &cobra.Command{
	Use:   "clone [url]",
	Short: "Generate git clone command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		directory, _ := cmd.Flags().GetString("directory")
		depth, _ := cmd.Flags().GetInt("depth")
		branch, _ := cmd.Flags().GetString("branch")

		var cmdParts []string
		cmdParts = append(cmdParts, "git clone")
		if depth > 0 {
			cmdParts = append(cmdParts, fmt.Sprintf("--depth %d", depth))
		}
		if branch != "" {
			cmdParts = append(cmdParts, fmt.Sprintf("-b %s", branch))
		}
		cmdParts = append(cmdParts, url)
		if directory != "" {
			cmdParts = append(cmdParts, directory)
		}

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

var gitStashCmd = &cobra.Command{
	Use:   "stash",
	Short: "Generate git stash commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		operation, _ := cmd.Flags().GetString("operation")
		message, _ := cmd.Flags().GetString("message")

		switch operation {
		case "save":
			if message != "" {
				fmt.Printf("git stash push -m \"%s\"\n", message)
			} else {
				fmt.Println("git stash push")
			}
		case "list":
			fmt.Println("git stash list")
		case "pop":
			fmt.Println("git stash pop")
		case "apply":
			fmt.Println("git stash apply")
		case "drop":
			fmt.Println("git stash drop")
		case "clear":
			fmt.Println("git stash clear")
		default:
			fmt.Println("git stash")
		}
		return nil
	},
}

var gitLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Generate git log command",
	RunE: func(cmd *cobra.Command, args []string) error {
		oneline, _ := cmd.Flags().GetBool("oneline")
		graph, _ := cmd.Flags().GetBool("graph")
		count, _ := cmd.Flags().GetInt("count")

		var cmdParts []string
		cmdParts = append(cmdParts, "git log")
		if oneline {
			cmdParts = append(cmdParts, "--oneline")
		}
		if graph {
			cmdParts = append(cmdParts, "--graph")
		}
		if count > 0 {
			cmdParts = append(cmdParts, fmt.Sprintf("-%d", count))
		}

		fmt.Println(strings.Join(cmdParts, " "))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
	gitCmd.AddCommand(gitCommitCmd)
	gitCmd.AddCommand(gitPushCmd)
	gitCmd.AddCommand(gitPullCmd)
	gitCmd.AddCommand(gitBranchCmd)
	gitCmd.AddCommand(gitMergeCmd)
	gitCmd.AddCommand(gitRebaseCmd)
	gitCmd.AddCommand(gitCloneCmd)
	gitCmd.AddCommand(gitStashCmd)
	gitCmd.AddCommand(gitLogCmd)

	gitCommitCmd.Flags().StringP("message", "m", "", "Commit message")
	gitCommitCmd.Flags().BoolP("all", "a", false, "Stage all modified files")
	gitCommitCmd.Flags().Bool("amend", false, "Amend previous commit")

	gitPushCmd.Flags().StringP("remote", "r", "origin", "Remote name")
	gitPushCmd.Flags().StringP("branch", "b", "", "Branch name")
	gitPushCmd.Flags().BoolP("force", "f", false, "Force push")
	gitPushCmd.Flags().Bool("tags", false, "Push tags")

	gitPullCmd.Flags().StringP("remote", "r", "origin", "Remote name")
	gitPullCmd.Flags().StringP("branch", "b", "", "Branch name")
	gitPullCmd.Flags().Bool("rebase", false, "Rebase instead of merge")

	gitBranchCmd.Flags().StringP("operation", "o", "list", "Operation (create, delete, rename, list, checkout)")
	gitBranchCmd.Flags().StringP("name", "n", "", "Branch name")
	gitBranchCmd.Flags().String("old-name", "", "Old branch name (for rename)")

	gitMergeCmd.Flags().Bool("no-ff", false, "Create merge commit even if fast-forward possible")

	gitRebaseCmd.Flags().BoolP("interactive", "i", false, "Interactive rebase")

	gitCloneCmd.Flags().StringP("directory", "d", "", "Directory name")
	gitCloneCmd.Flags().Int("depth", 0, "Create shallow clone with specified depth")
	gitCloneCmd.Flags().StringP("branch", "b", "", "Checkout specific branch")

	gitStashCmd.Flags().StringP("operation", "o", "save", "Operation (save, list, pop, apply, drop, clear)")
	gitStashCmd.Flags().StringP("message", "m", "", "Stash message")

	gitLogCmd.Flags().Bool("oneline", false, "Show one commit per line")
	gitLogCmd.Flags().Bool("graph", false, "Show ASCII graph")
	gitLogCmd.Flags().IntP("count", "n", 0, "Limit number of commits")
}
