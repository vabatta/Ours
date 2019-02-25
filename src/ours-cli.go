// Package ours contains the code to the ours application.
package main

// Imports declarations
import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/cbroglie/mustache"
)

// Constants
const (
	SyntaxVersion = "2.0"
	RuleSyntax    = `^/\* ours@(\d\.\d)\ \*/$`
	RuleActivity  = `^([A-Z0-9]+)@([A-Z]+|#[a-fA-F0-9]{6}#[a-fA-F0-9]{6})@(.+)$`
	RuleSlot      = `^([A-Z0-9]+):([-a-zA-Z0-9]*):([a-zA-Z0-9\s.-]+):(MON|TUE|WED|THU|FRI|01|02|03|04|05):((?:0[089]|1[0-9]|2[0-2])[0-5][0-9]):((?:0[089]|1[0-9]|2[0-3])[0-5][0-9])$`
	RuleComment   = `^#.+$`
	BaseTemplate  = "base.mustache"
	SlotTemplate  = "slot.mustache"
	StylesFile    = "stylus.css"
)

// BuiltInColors available: name -> Color
var BuiltInColors = map[string]Color{
	"green":     Color{Name: "green", Background: "#2ecc71", Foreground: "#fefefe"},
	"turquoise": Color{Name: "turquoise", Background: "#1abc9c", Foreground: "#fefefe"},
	"navy":      Color{Name: "navy", Background: "#34495e", Foreground: "#fefefe"},
	"blue":      Color{Name: "blue", Background: "#3498db", Foreground: "#fefefe"},
	"purple":    Color{Name: "purple", Background: "#9b59b6", Foreground: "#fefefe"},
	"grey":      Color{Name: "grey", Background: "#bdc3c7", Foreground: "#202020"},
	"red":       Color{Name: "red", Background: "#e74c3c", Foreground: "#fefefe"},
	"orange":    Color{Name: "orange", Background: "#f39c12", Foreground: "#fefefe"},
	"yellow":    Color{Name: "yellow", Background: "#f1c40f", Foreground: "#303030"},
}

// BuiltInDays is used to convert the slot rule days to a storable format
var BuiltInDays = map[string]byte{
	"MON": 0, "TUE": 1, "WED": 2, "THU": 3, "FRI": 4,
	"01": 0, "02": 1, "03": 2, "04": 3, "05": 4,
}

// Color is the structure representing an activity color.
type Color struct {
	Name       string
	Background string
	Foreground string
}

// Activity is the structure representing an activity in our timetable.
type Activity struct {
	ID    string
	Name  string
	Color Color
	Slots []Slot
}

// Slot is the structure representing a time slot in the timetable.
type Slot struct {
	Activity *Activity
	ID       string
	Icon     string
	Location string
	Day      byte
	Start    string
	End      string
}

// StartPrintable returns a human readable version of the start time.
func (s *Slot) StartPrintable() string {
	return fmt.Sprintf("%s:%s", s.Start[0:2], s.Start[2:4])
}

// EndPrintable returns a human readable version of the end time.
func (s *Slot) EndPrintable() string {
	return fmt.Sprintf("%s:%s", s.End[0:2], s.End[2:4])
}

// ActivityName returns the associated activity name of the slot.
func (s *Slot) ActivityName() string {
	return s.Activity.Name
}

// StartDelay returns the number of minutes since the beginning of the day's activity.
func (s *Slot) StartDelay() int {
	// Start hours and minutes
	sh, _ := strconv.Atoi(s.Start[0:2])
	sm, _ := strconv.Atoi(s.Start[2:4])
	// Compute and convert to int
	return (sh*60 + sm) - (8*60 + 00)
}

// Duration returns the duration time in minutes.
func (s *Slot) Duration() int {
	// Start hours and minutes
	sh, _ := strconv.Atoi(s.Start[0:2])
	sm, _ := strconv.Atoi(s.Start[2:4])
	// End hours and minutes
	eh, _ := strconv.Atoi(s.End[0:2])
	em, _ := strconv.Atoi(s.End[2:4])
	// Compute and convert to int
	return (eh*60 + em) - (sh*60 + sm)
}

// Styles returns the css to apply to a single slot.
func (s *Slot) Styles() string {
	// Reference to the activity color
	c := s.Activity.Color
	// Compute the height and delay
	h := (float64(s.Duration()) / 10)
	d := (float64(s.StartDelay()) / 10)
	// Return the css styles string
	return fmt.Sprintf(
		"height: %.1fvh; top: %.1fvh; background-color: %s; color: %s;",
		h, d, c.Background, c.Foreground,
	)
}

// Main function
func main() {
	// Check for arguments passed
	if len(os.Args) < 3 {
		log.Fatalf("Usage v%s: $> ours-cli <input file> <output file> [<templates dir>]", SyntaxVersion)
	}
	// Get the arguments
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	templatesDir := "templates"
	// Check if also the templates dir is passed or use the default
	if len(os.Args) > 3 {
		templatesDir = os.Args[3]
	}

	// Parse the input file
	parseInput(inputFile, outputFile, templatesDir)
}

// Returns a corresponding Color type to the given name or an error otherwise.
func convertColor(colorString string) (Color, error) {
	// Check if the color is in hex format
	hex := strings.Split(colorString, "#")
	if len(hex) > 1 {
		return Color{
			Name:       "custom hex color",
			Background: "#" + hex[1],
			Foreground: "#" + hex[2],
		}, nil
	}
	// Check if the given color exists
	for name, color := range BuiltInColors {
		if name == strings.ToLower(colorString) {
			return color, nil
		}
	}

	// Report error in case we didn't find a suitable color
	return Color{}, errors.New("color not valid")
}

func parseInput(inputFile string, outputFile string, templatesDir string) {
	// Open file and create scanner on top of it
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// Open the scanner
	scanner := bufio.NewScanner(file)
	// Split on new lines
	scanner.Split(bufio.ScanLines)

	// Init data buffers
	activities := make(map[string]*Activity)
	// Error found, prevent to process
	preventProcess := false
	// Line counter
	line := 0
	// Iterator on the content
	for scanner.Scan() {
		// Increment the line counter
		line++
		// Check for the syntax rule on the first rule
		if line == 1 {
			// Compile the rule and find all the groups
			regSyntax := regexp.MustCompile(RuleSyntax)
			syntaxMatches := regSyntax.FindAllStringSubmatch(scanner.Text(), -1)
			// Parsing error on the rule
			if syntaxMatches == nil {
				log.Fatalf("Line %d: Unsupported syntax file.", line)
			}
			// Mismatch error on the rule
			versionFound := syntaxMatches[0][1]
			if versionFound != SyntaxVersion {
				log.Fatalf(
					"Line %d: Mismatch file syntax version => Expecting %s, found %s",
					line, SyntaxVersion, versionFound,
				)
			}
			// Go to next line
			continue
		}
		// From second line parse other rules
		content := scanner.Text()

		// Compile the activity rule and find all the groups
		regActivity := regexp.MustCompile(RuleActivity)
		activityMatches := regActivity.FindAllStringSubmatch(content, -1)
		// Parsing error on the rule
		if activityMatches != nil {
			// Parsed groups
			id, colorString, name := activityMatches[0][1], activityMatches[0][2], activityMatches[0][3]

			// Check if id exists already
			_, exists := activities[id]
			if exists {
				log.Fatalf("Line %d: Activity with ID `%s` is already registered", line, id)
			}

			// Get the color if exists or report error
			color, err := convertColor(colorString)
			if err != nil {
				log.Fatalf("Line %d: Color %s not found in registered colors", line, colorString)
			}

			// Assign the new activity struct
			activities[id] = &Activity{
				Name:  name,
				Color: color,
				ID:    id,
				Slots: make([]Slot, 0, 5),
			}

			// Go to next line
			continue
		}

		// Compile the slot's activity rule and find all the groups
		regSlot := regexp.MustCompile(RuleSlot)
		slotMatches := regSlot.FindAllStringSubmatch(content, -1)
		// Parsing error on the rule
		if slotMatches != nil {
			// Check that activity exists
			id := slotMatches[0][1]
			activity, exists := activities[id]
			// Check if the id is present
			if !exists {
				log.Fatalf("Line %d: Activity with ID `%s` not found in registered activities", line, id)
			}
			// Convert the day to a number
			day, exists := BuiltInDays[slotMatches[0][4]]
			if !exists {
				log.Fatalf("Line %d: Day is incorrect -> %s", line, err)
			}

			// Create the slot's activity
			slot := Slot{
				Activity: activity,
				ID:       id,
				Icon:     slotMatches[0][2],
				Location: slotMatches[0][3],
				Day:      day,
				Start:    slotMatches[0][5],
				End:      slotMatches[0][6],
			}
			// Reallocate and add the slot to the activity
			activity.Slots = append(activity.Slots, slot)

			// Go to next line
			continue
		}

		// Skip comments
		regComment := regexp.MustCompile(RuleComment)
		slotComments := regComment.FindAllStringSubmatch(content, -1)
		// Parsing error on the rule
		if slotComments != nil {
			continue
		}

		// Skip empty lines
		if content == "" {
			continue
		}

		// fmt.Println(scanner.Text())
		log.Fatalf("Line %d: Invalid syntax rule.", line)
		preventProcess = true
	}

	// Check if we can proceed to process the input
	if !preventProcess {
		processInput(activities, outputFile, templatesDir)
	}
}

// Process the input to print it
func processInput(activities map[string]*Activity, outputFile string, templatesDir string) {
	// Map of rendered templates by days
	var days [5]bytes.Buffer
	// The slot template
	tmplSlot, err := mustache.ParseFile(fmt.Sprintf("%s/%s", templatesDir, SlotTemplate))
	if err != nil {
		log.Fatalf("There was an error while rendering your timetable -> %s", err)
	}
	// Buffer for the slot template render
	var buf bytes.Buffer
	// Iterate over the activites and their slots
	for _, activity := range activities {
		// Slots iterator
		for _, slot := range activity.Slots {
			// Reset the buffer
			buf.Reset()
			// Render the slot template
			err := tmplSlot.FRender(&buf, &slot)
			if err != nil {
				log.Fatalf("There was an error while rendering your timetable -> %s", err)
			}
			// Add it to the final buffer
			days[slot.Day].Write(buf.Bytes())
		}
	}

	// Open the css file
	styles, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", templatesDir, StylesFile))
	if err != nil {
		log.Fatalf("There was an error while opening the styles' file -> %s", err)
	}

	// Partials for the final template
	sp := mustache.StaticProvider{
		Partials: map[string]string{
			"monday":    days[0].String(),
			"tuesday":   days[1].String(),
			"wednesday": days[2].String(),
			"thursday":  days[3].String(),
			"friday":    days[4].String(),
			"css":       string(styles),
		},
	}

	// Render final output
	tmpl, err := mustache.ParseFilePartials(fmt.Sprintf("%s/%s", templatesDir, BaseTemplate), &sp)
	if err != nil {
		log.Fatalf("There was an error while rendering your timetable -> %s", err)
	}
	// Reset the buffer and render
	buf.Reset()
	tmpl.FRender(&buf, nil)

	// Write out the result
	err = ioutil.WriteFile(outputFile, buf.Bytes(), 0644)
	// Check error and report the status
	if err != nil {
		log.Fatalf("There was an error while saving your file -> %s", err)
	}
	// Else everything was good
	log.Printf("Done! Your timetable is available in the file %s", outputFile)
}
