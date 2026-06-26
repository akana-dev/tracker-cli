package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"tracker/internal/client"
	"tracker/internal/config"
	"tracker/internal/ui"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Вход в систему",
	RunE: func(cmd *cobra.Command, args []string) error {
		method, _ := cmd.Flags().GetString("method")

		if method == "" || method == "password" && !cmd.Flags().Changed("method") {
			return loginAuto()
		}

		switch method {
		case "ad":
			return loginAD()
		case "password":
			return loginPassword()
		default:
			return fmt.Errorf("неизвестный метод: %s (допустимы: password, ad)", method)
		}
	},
}

func loginAuto() error {
	server, err := config.GetCurrentServer()
	if err != nil {
		return err
	}

	methods := server.AuthMethods
	if len(methods) == 0 {
		methods = []string{"password"}
	}

	if len(methods) == 1 {
		method := strings.ToLower(strings.TrimSpace(methods[0]))
		fmt.Println(ui.Dimf("Используется метод: %s", method))
		switch method {
		case "ad":
			return loginAD()
		case "password":
			return loginPassword()
		default:
			return fmt.Errorf("неизвестный метод в конфиге: %s", method)
		}
	}

	fmt.Println(ui.Dimf("Доступные методы: %s", strings.Join(methods, ", ")))
	fmt.Println()

	var lastErr error
	for _, method := range methods {
		method = strings.ToLower(strings.TrimSpace(method))
		fmt.Println(ui.Boldf("Пробуем: %s", method))

		switch method {
		case "ad":
			err = loginAD()
		case "password":
			err = loginPassword()
		default:
			fmt.Println(ui.Warningf("Неизвестный метод: %s, пропускаем", method))
			continue
		}

		if err == nil {
			return nil
		}

		fmt.Println(ui.Warningf("Метод %s не сработал: %v", method, err))
		fmt.Println()
		lastErr = err
	}

	return fmt.Errorf("все методы авторизации не сработали. Последняя ошибка: %w", lastErr)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Выход из системы",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.SaveToken(""); err != nil {
			return err
		}
		fmt.Println(ui.Checkmark(), ui.Success("Выход выполнен"))
		return nil
	},
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Информация о текущем пользователе",
	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := client.GetMe()
		if err != nil {
			return err
		}

		fmt.Println()
		ui.Header("Пользователь:")
		ui.Label("ID", fmt.Sprintf("%d", user.ID))

		if user.FullName != nil && *user.FullName != "" {
			ui.Label("ФИО", ui.Bold(*user.FullName))
			ui.Label("Логин", ui.Cyan(user.Username))
		} else {
			ui.Label("Логин", ui.Cyan(user.Username))
		}

		ui.Label("Email", user.Email)
		ui.Label("Роль", ui.RoleColor(user.Role))
		ui.Label("Активен", func() string {
			if user.IsActive {
				return ui.StatusOK()
			}
			return ui.StatusNo()
		}())
		if !user.CreatedAt.IsZero() {
			ui.Label("Создан", user.CreatedAt.Format("02.01.2006 15:04"))
		}
		fmt.Println()
		return nil
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Регистрация нового пользователя",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Имя пользователя: ")
		username := readLine()
		if username == "" {
			return fmt.Errorf("имя пользователя не может быть пустым")
		}

		fmt.Print("Email: ")
		email := readLine()
		if email == "" {
			return fmt.Errorf("email не может быть пустым")
		}

		fmt.Print("Пароль: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Println()

		if err := client.RegisterUser(username, email, string(passwordBytes)); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Пользователь %s зарегистрирован", ui.Bold(username)))
		fmt.Println(ui.Dim("Теперь выполните: tracker login"))
		return nil
	},
}

func init() {
	loginCmd.Flags().StringP("method", "m", "", "Метод авторизации: password, ad (по умолчанию — из конфига)")
}

func loginPassword() error {
	fmt.Print("Логин трекера: ")
	username := readLine()
	if username == "" {
		return fmt.Errorf("логин не может быть пустым")
	}

	fmt.Print("Пароль: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println()

	resp, err := client.LoginPassword(username, string(passwordBytes))
	if err != nil {
		fmt.Println(ui.Cross(), ui.Error("Ошибка авторизации"))
		return fmt.Errorf("%w", err)
	}

	if err := config.SaveToken(resp.AccessToken); err != nil {
		return err
	}

	displayName := getUserDisplayName()

	fmt.Println(ui.Checkmark(), ui.Successf("Вход выполнен через %s как %s",
		ui.Bold("password"), ui.Bold(displayName)))
	return nil
}

func loginAD() error {
	domain := config.GetADDomain()
	if domain == "" {
		return fmt.Errorf("не настроен AD домен. Выполните: tracker configure")
	}

	fmt.Printf("Логин AD (user@%s): ", domain)
	username := readLine()
	if username == "" {
		return fmt.Errorf("логин не может быть пустым")
	}

	if !strings.Contains(username, "@") && !strings.Contains(username, "\\") {
		username = fmt.Sprintf("%s@%s", username, domain)
	}

	fmt.Print("Пароль AD: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println()

	resp, err := client.LoginAD(username, string(passwordBytes))
	if err != nil {
		fmt.Println(ui.Cross(), ui.Error("Ошибка AD"))
		return fmt.Errorf("%w", err)
	}

	if err := config.SaveToken(resp.AccessToken); err != nil {
		return err
	}

	displayName := getUserDisplayName()

	fmt.Println(ui.Checkmark(), ui.Successf("Вход выполнен через %s как %s",
		ui.Bold("AD"), ui.Bold(displayName)))
	return nil
}

func getUserDisplayName() string {
	user, err := client.GetMe()
	if err != nil {
		return ""
	}

	config.SaveUserRole(user.Role)

	return user.GetFullName()
}

func formatUserName(username string, fullName *string) string {
	if fullName != nil && *fullName != "" && *fullName != username {
		return fmt.Sprintf("%s (%s)", *fullName, username)
	}
	return username
}
