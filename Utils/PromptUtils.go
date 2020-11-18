package Utils

import (
	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
	"os"
)

func YesNoPrompt(label string) bool {
	prompt := promptui.Select{
		Label:        label,
		Items:        []string{"Yes", "No"},
		HideSelected: true,
	}
	_, result, err := prompt.Run()
	if err != nil {
		pterm.Error.Printf("\n something went wrong: %v\n", err)
		os.Exit(0)
	}
	if result == "Yes" {
		return true
	} else {
		return false
	}
}

func SingleChoicePrompt(label string, items []string, showSelected bool) (int, string) {
	prompt := promptui.Select{
		Label:        label,
		Items:        items,
		HideSelected: showSelected,
	}
	index, result, err := prompt.Run()
	if err != nil {
		pterm.Error.Printf("\n something went wrong: %v\n", err)
		os.Exit(0)
	}
	return index, result
}

func InputPrompt(title string, validator func(input string) error) string {
	//validate := func(input string) error {
	//	if net.ParseIP(input) == nil {
	//		return errors.New("invalid number IP address format")
	//	}
	//	return nil
	//}

	prompt := promptui.Prompt{
		Label:    title,
		Validate: validator,
	}

	result, err := prompt.Run()

	if err != nil {
		pterm.Error.Printf("\n something went wrong: %v\n", err)
		os.Exit(0)
	}
	return result
}
