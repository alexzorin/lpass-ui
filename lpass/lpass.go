// lpass provides the interface to the lpass CLI

package lpass

import (
	"os/exec"
	"regexp"

	"github.com/pkg/errors"
)

var (
	ErrNotLoggedIn      = errors.New("Not logged in")
	regexStatusLoggedIn = regexp.MustCompile(`^Logged in as (.*)\.`)
)

func Exec(arguments ...string) (string, error) {
	bin, err := exec.LookPath("lpass")
	if err != nil {
		return "", errors.Wrap(err, "Unable to find lpass bin")
	}

	arguments = append(arguments, "--color=never")

	cmd := exec.Command(bin, arguments...)
	buf, err := cmd.CombinedOutput()

	return string(buf), err
}

func CheckLoggedIn() (string, error) {
	res, err := Exec("status")
	if err != nil {
		return "", ErrNotLoggedIn
	}

	matches := regexStatusLoggedIn.FindAllStringSubmatch(res, -1)

	if len(matches) != 1 {
		return "", ErrNotLoggedIn
	}

	return matches[0][1], nil
}
