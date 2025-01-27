package parser

import (
	"os/exec"
	"strings"
)

type CTagsOptions struct {
	LangExtension string
	ListKinds     []string
	LangMap       map[string]string
}

func RunCtags(filename string, options *CTagsOptions) ([]string, error) {
	cmdArgs := []string{"--output-format=json", "--fields=+l+n"}
	if options != nil {
		if options.LangExtension != "" {
			kinds := strings.Join(options.ListKinds, "")
			if kinds == "" {
				kinds = "f"
			}
			cmdArgs = append(cmdArgs, "--"+options.LangExtension+"-kinds="+kinds)
		}
	}
	cmdArgs = append(cmdArgs, filename)
	cmd, err := exec.Command("ctags", cmdArgs...).Output()
	if err != nil {
		return nil, err
	}

	sp := strings.Split(string(cmd), "\n")
	return sp, nil
}
