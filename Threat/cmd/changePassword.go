
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type ChangePassword struct {
	email           string
	currentPassword string
	newPassword     string
	confirmPassword string
}

// changePasswordCmd represents the changePassword command
var changePasswordCmd = &cobra.Command{
	Use:   "change-password",
	Short: "Change User Password",
	Long:  `Changes the user password from the commandline`,
	Run: func(cmd *cobra.Command, args []string) {
		if IsConfigured() {
			email, _ := cmd.Flags().GetString("email")
			currentPassword, _ := cmd.Flags().GetString("current-password")
			newPassword, _ := cmd.Flags().GetString("new-password")
			confirmPassword, _ := cmd.Flags().GetString("confirm-password")

			if currentPassword == "" || newPassword == "" || confirmPassword == "" {
				old, new, confirm := getChangePasswordPrompt()
				initiateChangePassword(ChangePassword{
					email:           email,
					currentPassword: old,
					newPassword:     new,
					confirmPassword: confirm,
				})
			} else {
				if newPassword == confirmPassword {
					initiateChangePassword(ChangePassword{
						email:           email,
						currentPassword: currentPassword,
						newPassword:     newPassword,
						confirmPassword: confirmPassword,
					})
				} else {
					log.Fatal("New password and Confirm Password values don't match")
				}
			}

		} else {
			log.Fatal("You need to configure ThreatPlaybook first. Run `configure` first")
		}

	},
}

func init() {
	rootCmd.AddCommand(changePasswordCmd)

	changePasswordCmd.Flags().StringP("email", "e", "", "User's email")
	changePasswordCmd.Flags().StringP("current-password", "c", "", "Current User Password")
	changePasswordCmd.Flags().StringP("new-password", "n", "", "New Password")
	changePasswordCmd.Flags().StringP("confirm-password", "r", "", "Confirm New Password")
	changePasswordCmd.MarkFlagRequired("email")
}

func getChangePasswordPrompt() (string, string, string) {
	oldPasswordPrompt := promptui.Prompt{
		Label: "Current Password",
		Mask:  '*',
	}
	newPasswordPrompt := promptui.Prompt{
		Label: "New Password",
		Mask:  '*',
	}
	confirmPasswordPrompt := promptui.Prompt{
		Label: "Confirm Password",
		Mask:  '*',
	}

	oldPassword, oldErr := oldPasswordPrompt.Run()
	if oldErr != nil {
		log.Fatal("Unable to read old password")
	}
	newPassword, newErr := newPasswordPrompt.Run()
	if newErr != nil {
		log.Fatal("Unable to read new password")
	}
	confirmPassword, confirmErr := confirmPasswordPrompt.Run()
	if confirmErr != nil {
		log.Fatal("Unable to read confirm password")
	}

	if newPassword != confirmPassword {
		log.Fatal("New and Confirm Password values dont match")
	}

	return oldPassword, newPassword, confirmPassword

}

func initiateChangePassword(passChange ChangePassword) {
	configValue := GetJsonConfiguration()
	if (ConfigObj{}) == configValue {
		log.Fatal("Unable to fetch value from cred file")
	}
	url := fmt.Sprintf("http://%s:%d/api/change-password", configValue.host, configValue.port)

	requestBody, json_err := json.Marshal(map[string]string{
		"email":           passChange.email,
		"old_password":    passChange.currentPassword,
		"new_password":    passChange.newPassword,
		"verify_password": passChange.confirmPassword,
	})

	if json_err != nil {
		log.Fatal("Unable to generate HTTP request")
	}

	resp, httpError := http.Post(url, "application/json", bytes.NewBuffer(requestBody))

	if httpError != nil {
		log.Fatal("Unable to generate HTTP request to change password")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(*resp)
		log.Fatal("Unable to change password")
	}

	fmt.Println("Successfully changed user password")

}
