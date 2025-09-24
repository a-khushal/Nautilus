package jobs

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/a-khushal/Nautilus/worker/models"
)

var langMap = map[string]struct {
	Ext   string
	Image string
}{
	"python3": {Ext: "py", Image: "python-runner:latest"},
	"cpp":     {Ext: "cpp", Image: "cpp-runner:latest"},
	"java":    {Ext: "java", Image: "java-runner:latest"},
}

func RunCodeJob(job models.Job) {
	lang, ok := job.Payload["lang"].(string)
	if !ok {
		fmt.Println("Invalid language in payload")
		return
	}

	data, ok := job.Payload["data"].(string)
	if !ok {
		fmt.Println("Invalid code data in payload")
		return
	}

	langInfo, ok := langMap[lang]
	if !ok {
		fmt.Println("Unsupported language:", lang)
		return
	}

	filename := fmt.Sprintf("/tmp/code_%s.%s", job.ID, langInfo.Ext)
	if err := os.WriteFile(filename, []byte(data), 0644); err != nil {
		fmt.Println("Failed to write code file:", err)
		return
	}

	var innerCmd string
	if lang == "cpp" {
		binName := fmt.Sprintf("/tmp/code_%s.out", job.ID)
		innerCmd = fmt.Sprintf("g++ %s -o %s && %s", filename, binName, binName)
	} else {
		innerCmd = fmt.Sprintf("%s %s", lang, filename)
	}

	dockerArgs := []string{
		"run", "--rm",
		"-v", "/tmp:/tmp",
		"--cpus=1",
		"--memory=128m",
		"--pids-limit=50",
		langInfo.Image,
		"bash", "-c", innerCmd,
	}

	cmd := exec.Command("docker", dockerArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Execution failed:", err)
	}

	fmt.Println("Execution output:", string(output))
}
