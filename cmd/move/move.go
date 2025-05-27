package move

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/structs"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strconv"
	"strings"
)

func MoveCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about a specific move.\n\n",
			styling.StyleBold.Render("USAGE:"),
			"\n\t"+"poke-cli"+" "+styling.StyleBold.Render("move")+" <move-name>",
			"\n\n"+styling.StyleItalic.Render("Use a hyphen when typing a name with a space."),
		)
		output.WriteString(helpMessage)
	}

	flag.Parse()

	if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		flag.Usage()
		return output.String(), nil
	}

	if err := utils.ValidateMoveArgs(os.Args); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	args := flag.Args()
	endpoint := strings.ToLower(args[0])
	moveName := strings.ToLower(args[1])

	moveStruct, moveName, _ := connections.MoveApiCall(endpoint, moveName, connections.APIURL)

	moveInfoContainer(&output, moveStruct, moveName)
	moveEffectContainer(&output, moveStruct)

	return output.String(), nil
}

func moveInfoContainer(output *strings.Builder, moveStruct structs.MoveJSONStruct, moveName string) {
	capitalizedMove := cases.Title(language.English).String(strings.ReplaceAll(moveName, "-", " "))

	docStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(styling.GetTypeColor(moveStruct.Type.Name))).
		Width(32)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(styling.GetTypeColor(moveStruct.Type.Name))).
		PaddingBottom(1)

	labelStyle := lipgloss.NewStyle().Bold(true).Width(15)
	valueStyle := lipgloss.NewStyle().Faint(true)

	header := headerStyle.Render(capitalizedMove)

	infoRows := []string{
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Type"), "|", valueStyle.Render(cases.Title(language.English).String(moveStruct.Type.Name))),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Power"), "|", valueStyle.Render(strconv.Itoa(moveStruct.Power))),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("PP"), "|", valueStyle.Render(strconv.Itoa(moveStruct.PowerPoints))),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Accuracy"), "|", valueStyle.Render(strconv.Itoa(moveStruct.Accuracy))),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Category"), "|", valueStyle.Render(cases.Title(language.English).String(moveStruct.DamageClass.Name))),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Effect Chance"), "|", valueStyle.Render(fmt.Sprintf("%d%%", moveStruct.EffectChance))),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Priority"), "|", valueStyle.Render(strconv.Itoa(moveStruct.Priority))),
	}

	infoBlock := lipgloss.JoinVertical(lipgloss.Left, infoRows...)
	fullDoc := lipgloss.JoinVertical(lipgloss.Top, header, infoBlock)

	output.WriteString(docStyle.Render(fullDoc))
	output.WriteString("\n")
}

func moveEffectContainer(output *strings.Builder, moveStruct structs.MoveJSONStruct) {
	docStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(styling.GetTypeColor(moveStruct.Type.Name))).
		Width(32)

	var flavorTextEntry string
	for _, entry := range moveStruct.FlavorTextEntries {
		if entry.Language.Name == "en" && entry.VersionGroup.Name == "scarlet-violet" {
			flavorTextEntry = entry.FlavorText
			break
		}
	}

	effectBold := styling.StyleBold.Render("Effect:")
	fullDoc := lipgloss.JoinVertical(lipgloss.Top, effectBold, flavorTextEntry)

	output.WriteString(docStyle.Render(fullDoc))
	output.WriteString("\n")
}
