
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

func passwordPrompt() string {
	oldPasswordPrompt := promptui.Prompt{
		Label: "Enter Password",
		Mask:  '*',
	}

	oldPassword, oldErr := oldPasswordPrompt.Run()
	if oldErr != nil {
		log.Fatal("Unable to read old password")
	}
	return oldPassword

}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to the ThreatPlaybook API with an email and password",
	Long:  `Authenticate to the ThreatPlaybook API with an email and password`,
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		if !IsConfigured() {
			log.Fatal("First setup configure for ThreatPlaybook client")
		}

		if password == "" {
			password = passwordPrompt()
			initiateLoginRequest(email, password)
		} else {
			initiateLoginRequest(email, password)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("email", "e", "", "User's email")
	loginCmd.Flags().StringP("password", "p", "", "User Password")
	loginCmd.MarkFlagRequired("email")

}

func initiateLoginRequest(email string, password string) {
	configValue := GetJsonConfiguration()
	if (ConfigObj{}) == configValue {
		log.Fatal("Unable to fetch value from cred file")
	}
	url := fmt.Sprintf("http://%s:%d/api/login", configValue.host, configValue.port)

	requestBody, json_err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
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
		log.Fatal("Unable to login")
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	token := result["token"].(string) // type assertion
	ChangeTokenInCred(token)
	fmt.Println("Successfully logged in")
}
