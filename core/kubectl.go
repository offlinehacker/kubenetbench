package core

import (
	"fmt"
	"log"
	"time"

	"../utils"
)

// KubeGetPodIP returns the IP address of a pod using a proper selector
func (c *RunCtx) KubeGetPodIP(
	selector string,
	retries uint,
	st time.Duration,
) (string, error) {

	retries_orig := retries
	cmd := fmt.Sprintf(
		"kubectl get pod -l \"%s\" -o custom-columns=IP:.status.podIP --no-headers",
		selector,
	)
	for {
		if !c.quiet {
			log.Printf("$ %s # (remaining retries: %d)", cmd, retries)
		}
		lines, err := utils.ExecCmdLines(cmd)
		if err == nil && len(lines) == 1 && lines[0] != "<none>" {
			return lines[0], nil
		}

		if retries == 0 {
			return lines[0], fmt.Errorf("Error exeucting %s after %d retries (last error:%w)", cmd, retries_orig, err)
		}

		retries--
		time.Sleep(st)
	}
}

// KubeGetPodPhase returns the phase of a pod
func (c *RunCtx) KubeGetPodPhase(selector string) (string, error) {
	cmd := fmt.Sprintf(
		"kubectl get pod -l \"%s\" -o custom-columns=Status:.status.phase --no-headers",
		selector,
	)

	lines, err := utils.ExecCmdLines(cmd)
	if err != nil {
		return "", fmt.Errorf("command %s failed: %w", cmd, err)
	}

	if len(lines) != 1 {
		log.Fatal("Selector did not provide single result cmd:", cmd, "result:", lines)
		return "", fmt.Errorf("selector %s did not provide single result: command: %s; result: %s", selector, cmd, lines)
	}

	return lines[0], nil
}

// KubeSaveLogs saves logs using a selector
func (c *RunCtx) KubeSaveLogs(selector string, logfile string) error {
	argcmd := fmt.Sprintf("kubectl logs -l \"%s\" > %s", selector, logfile)
	if !c.quiet {
		log.Printf("$ %s ", argcmd)
	}
	return utils.ExecCmd(argcmd)
}

// KubeGetServiceIP returns the ip of a service
// NB: probably a better option to use DNS
func (c *RunCtx) KubeGetServiceIP(selector string) (string, error) {
	cmd := fmt.Sprintf(
		"kubectl get service -l '%s' -o custom-columns=IP:.spec.clusterIP --no-headers",
		selector,
	)

	lines, err := utils.ExecCmdLines(cmd)
	if err != nil {
		return "", fmt.Errorf("Error executing: %s: %w", cmd, err)
	}

	if len(lines) != 1 {
		return "", fmt.Errorf("Error executing: %s: selector did not return a singel line (result: %s)", cmd, lines)
	}

	return lines[0], nil
}

// Kubectl apply
func (c *RunCtx) KubeApply(fname string) error {
	cmd := fmt.Sprintf("kubectl apply -f %s", fname)
	if !c.quiet {
		log.Printf("$ %s ", cmd)
	}
	return utils.ExecCmd(cmd)
}

func (c *RunCtx) KubeCleanup() error {
	cmd := fmt.Sprintf("kubectl delete pod,networkpolicy -l \"runid=%s\"", c.id)
	if !c.quiet {
		log.Printf("$ %s ", cmd)
	}
	return utils.ExecCmd(cmd)
}