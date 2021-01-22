// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"github.com/spf13/viper"
	"time"
	"wraith/core"
	"wraith/version"

	"fmt"
	"github.com/spf13/cobra"
)

var viperScanGitlab *viper.Viper

// scanGitlabCmd represents the scanGitlab command
var scanGitlabCmd = &cobra.Command{
	Use:   "scanGitlab",
	Short: "Scan one or more gitlab groups or users for secrets",
	Long:  `Scan one or more gitlab groups or users for secrets`,
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "gitlab"
		sess := core.NewSession(viperScanGitlab, scanType)

		sess.Out.Warn("%s\n\n", core.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", sess.BindAddress, sess.BindPort)

		sess.GitlabAccessToken = viperScanGitlab.GetString("gitlab-api-token")

		sess.InitGitClient()

		core.GatherTargets(sess)
		core.GatherGitlabRepositories(sess)
		core.AnalyzeRepositories(sess)
		sess.Finish()

		core.PrintSessionStats(sess)

		if !sess.Silent {
			sess.Out.Important("%s", core.GitLabTanuki)
			sess.Out.Important("Press Ctrl+C to stop web server and exit.\n")
			select {}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanGitlabCmd)

	viperScanGitlab = core.SetConfig()

	scanGitlabCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanGitlabCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanGitlabCmd.Flags().Int("confidence-level", 3, "The confidence level level of the expressions used to find matches")
	scanGitlabCmd.Flags().Float64("commit-depth", -1, "Set the commit depth to scan")
	scanGitlabCmd.Flags().Bool("debug", false, "Print debugging information")
	scanGitlabCmd.Flags().Bool("gather-org-members", false, "Add members to targets when processing organizations")
	scanGitlabCmd.Flags().String("gitlab-api-token", "", "API token for access to gitlab, see doc for necessary scope")
	scanGitlabCmd.Flags().StringSlice("github-projects", nil, "List of Gitlab projects or users to scan")
	scanGitlabCmd.Flags().Bool("hide-secrets", false, "Hide secrets from any supported output")
	scanGitlabCmd.Flags().StringSlice("ignore-extension", nil, "List of extensions to ignore")
	scanGitlabCmd.Flags().StringSlice("ignore-path", nil, "List of paths to ignore")
	//scanGitlabCmd.Flags().Bool("in-mem-clone", false, "Clone repos in memory")
	scanGitlabCmd.Flags().Int("max-file-size", 10, "Max file size to scan in MB")
	scanGitlabCmd.Flags().Int("num-threads", -1, "Number of threads to execute with")
	scanGitlabCmd.Flags().Bool("scan-forks", false, "Scan repositories forked by users or orgs")
	scanGitlabCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanGitlabCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yml", "file(s) containing detection signatures.")
	scanGitlabCmd.Flags().String("signature-path", "$HOME/.wraith/signatures", "path containing detection signatures.")
	scanGitlabCmd.Flags().Bool("silent", false, "Suppress all output except for errors")

	err := viperScanGitlab.BindPFlag("bind-address", scanGitlabCmd.Flags().Lookup("bind-address"))
	err = viperScanGitlab.BindPFlag("bind-port", scanGitlabCmd.Flags().Lookup("bind-port"))
	err = viperScanGitlab.BindPFlag("commit-depth", scanGitlabCmd.Flags().Lookup("commit-depth"))
	err = viperScanGitlab.BindPFlag("debug", scanGitlabCmd.Flags().Lookup("debug"))
	err = viperScanGitlab.BindPFlag("gather-org-members", scanGitlabCmd.Flags().Lookup("gather-org-members"))
	err = viperScanGitlab.BindPFlag("gitlab-api-token", scanGitlabCmd.Flags().Lookup("gitlab-api-token"))
	err = viperScanGitlab.BindPFlag("hide-secrets", scanGitlabCmd.Flags().Lookup("hide-secrets"))
	err = viperScanGitlab.BindPFlag("ignore-extension", scanGitlabCmd.Flags().Lookup("ignore-extension"))
	err = viperScanGitlab.BindPFlag("ignore-path", scanGitlabCmd.Flags().Lookup("ignore-extension"))
	//err = viperScanGitlab.BindPFlag("in-mem-clone", scanGitlabCmd.Flags().Lookup("in-mem-clone"))
	err = viperScanGitlab.BindPFlag("confidence-level", scanGitlabCmd.Flags().Lookup("confidence-level"))
	err = viperScanGitlab.BindPFlag("max-file-size", scanGitlabCmd.Flags().Lookup("max-file-size"))
	err = viperScanGitlab.BindPFlag("num-threads", scanGitlabCmd.Flags().Lookup("num-threads"))
	err = viperScanGitlab.BindPFlag("scan-forks", scanGitlabCmd.Flags().Lookup("scan-forks"))
	err = viperScanGitlab.BindPFlag("scan-tests", scanGitlabCmd.Flags().Lookup("scan-tests"))
	err = viperScanGitlab.BindPFlag("signature-file", scanGitlabCmd.Flags().Lookup("signature-file"))
	err = viperScanGitlab.BindPFlag("signature-path", scanGitlabCmd.Flags().Lookup("signature-path"))
	err = viperScanGitlab.BindPFlag("silent", scanGitlabCmd.Flags().Lookup("silent"))
	err = viperScanGitlab.BindPFlag("gitlab-projects", scanGitlabCmd.Flags().Lookup("gitlab-projects"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
