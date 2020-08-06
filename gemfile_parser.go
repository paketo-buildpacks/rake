package rake

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type GemfileParser struct{}

func NewGemfileParser() GemfileParser {
	return GemfileParser{}
}

func (p GemfileParser) Parse(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("failed to parse Gemfile: %w", err)
	}
	defer file.Close()

	quotes := `["']`
	rakeRe := regexp.MustCompile(fmt.Sprintf(`^gem %srake%s`, quotes, quotes))

	hasRake := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := []byte(scanner.Text())

		if !hasRake {
			hasRake = rakeRe.Match(line)
		}

		if hasRake {
			return true, nil
		}
	}

	return false, nil
}
