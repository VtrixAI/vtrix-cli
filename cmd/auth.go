package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/VtrixAI/vtrix-cli/internal/auth"
	"github.com/VtrixAI/vtrix-cli/internal/clierrors"
	"github.com/VtrixAI/vtrix-cli/internal/config"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to your vtrix account",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err == nil && cfg.AuthToken != "" {
			me, err := auth.VerifyToken(cfg.AuthToken)
			if err == nil {
				fmt.Printf("Already logged in as %s. Run vtrix auth logout to switch accounts.\n", me.Email)
				return nil
			}
		}

		token, refreshToken, apiKey, err := auth.Login(openBrowser)
		if err != nil {
			return err
		}

		me, err := auth.VerifyToken(token)
		if err != nil {
			return clierrors.ErrTokenVerification(err)
		}

		if err := config.Save(&config.Config{AuthToken: token, RefreshToken: refreshToken, APIKey: apiKey}); err != nil {
			return clierrors.ErrSaveConfig(err)
		}

		email := me.Email
		if email == "" {
			email = me.Account
		}
		fmt.Printf("\nLogged in as %s\n", email)
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}

		if cfg.AuthToken == "" {
			fmt.Println("Not logged in.")
			fmt.Println("  Hint: Run: vtrix auth login")
			return nil
		}

		me, err := auth.VerifyToken(cfg.AuthToken)
		if err != nil {
			fmt.Println("Token expired or invalid.")
			fmt.Println("  Hint: Run: vtrix auth login")
			return nil
		}

		email := me.Email
		if email == "" {
			email = me.Account
		}
		fmt.Printf("Logged in as %s\n", email)
		if me.Name != "" {
			fmt.Printf("Name: %s\n", me.Name)
		}
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of your vtrix account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Clear(); err != nil {
			return clierrors.ErrLogout(err)
		}
		fmt.Println("Logged out.")
		return nil
	},
}

var authSetKeyCmd = &cobra.Command{
	Use:   "set-key <api-key>",
	Short: "Set or replace the vtrix API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		newKey := args[0]
		if newKey == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		cfg, err := config.Load()
		if err != nil {
			cfg = &config.Config{}
		}

		cfg.APIKey = newKey
		if err := config.Save(cfg); err != nil {
			return clierrors.ErrSaveConfig(err)
		}

		fmt.Println("API key updated.")
		return nil
	},
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		return fmt.Errorf("unsupported operating system")
	}

	return exec.Command(cmd, args...).Start()
}

func init() {
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authSetKeyCmd)
}
