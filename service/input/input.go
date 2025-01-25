package input

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
)

// GetInput获取用户终端输入
func GetInput(label string) (string, error) {
	for {
		var input string

		// 使用 survey.Input 类型
		prompt := &survey.Input{
			Message: label, // 提示信息
		}

		// 读取用户输入
		err := survey.AskOne(prompt, &input)
		if err != nil {
			return "", nil
		}
		if input != "" {
			return input, nil
		}
		color.Yellow.Println("Please do not leave the input empty.")
	}
}

// GetOptionInput获取可选填的数据
func GetOptionalInput(label string) (string, error) {
	var input string

	// 使用 survey.Input 类型
	prompt := &survey.Input{
		Message: label, // 提示信息
	}

	// 读取用户输入
	err := survey.AskOne(prompt, &input)
	if err != nil {
		return "", nil
	}
	return input, nil
}
func GetPasswordInput(label string) (string, error) {
	for {
		var password string

		// 使用 survey 的 Password 类型
		prompt := &survey.Password{
			Message: label,
		}

		// 读取密码输入
		err := survey.AskOne(prompt, &password)
		if err != nil {
			return "", err
		}
		if password != "" {
			return password, nil
		}
		color.Yellow.Println("Please do not leave the input empty.")

	}
}

func GetOptionalPassword(label string) (string, error) {
	var password string

	// 使用 survey 的 Password 类型
	prompt := &survey.Password{
		Message: label,
	}

	// 读取密码输入
	err := survey.AskOne(prompt, &password)
	if err != nil {
		return "", err
	}

	return password, nil
}
