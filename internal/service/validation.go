package service

import (
	"fmt"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("название задачи не может быть пустым")
	}
	if len(title) > MaxTitleLength {
		return fmt.Errorf("название задачи слишком длинное (максимум %d символов, у вас %d)",
			MaxTitleLength, len(title))
	}
	return nil
}

func ValidateComment(comment string) error {
	if len(comment) > MaxCommentLength {
		return fmt.Errorf("комментарий слишком длинный (максимум %d символов, у вас %d)",
			MaxCommentLength, len(comment))
	}
	return nil
}

func ValidateSolution(solution string) error {
	if len(solution) > MaxSolutionLength {
		return fmt.Errorf("статус решения слишком длинный (максимум %d символов, у вас %d)",
			MaxSolutionLength, len(solution))
	}
	return nil
}

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("имя пользователя не может быть пустым")
	}
	if len(username) > MaxUsernameLength {
		return fmt.Errorf("имя пользователя слишком длинное (максимум %d символов, у вас %d)",
			MaxUsernameLength, len(username))
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9._\-]+$`).MatchString(username) {
		return fmt.Errorf("имя пользователя может содержать только буквы, цифры, точки, дефисы и подчёркивания")
	}
	return nil
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email не может быть пустым")
	}
	if len(email) > MaxEmailLength {
		return fmt.Errorf("email слишком длинный (максимум %d символов, у вас %d)",
			MaxEmailLength, len(email))
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("некорректный формат email: %s", email)
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("пароль должен содержать минимум 6 символов")
	}
	if len(password) > 128 {
		return fmt.Errorf("пароль слишком длинный (максимум 128 символов)")
	}
	return nil
}

func ValidateCompanyName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("название компании не может быть пустым")
	}
	if len(name) > MaxCompanyNameLength {
		return fmt.Errorf("название компании слишком длинное (максимум %d символов, у вас %d)",
			MaxCompanyNameLength, len(name))
	}
	return nil
}

func ValidateCompanyDescription(description string) error {
	if len(description) > MaxCompanyDescriptionLength {
		return fmt.Errorf("описание компании слишком длинное (максимум %d символов, у вас %d)",
			MaxCompanyDescriptionLength, len(description))
	}
	return nil
}
