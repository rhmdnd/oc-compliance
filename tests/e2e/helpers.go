package e2e

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const defaultOCWaitTimeout = "--timeout=60s"
const defaultOCLongWaitTimeout = "--timeout=10m"

func do(cmd string, args ...string) string {
	execcmd := exec.Command(cmd, args...)
	output, err := execcmd.CombinedOutput()
	Expect(err).ShouldNot(HaveOccurred(),
		"The command '%s' shouldn't fail.\n- Arguments: %v\n- Output: %s", cmd, args, output)
	return strings.Trim(string(output), "\n")
}

func oc(args ...string) string {
	return do("oc", args...)
}

func ocApplyFromString(obj string) string {
	tmpfile, err := ioutil.TempFile("", "oc-create")
	Expect(err).ShouldNot(HaveOccurred(), "Creating a temp file shouldn't fail")
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	_, writeErr := io.WriteString(tmpfile, obj)
	Expect(writeErr).ShouldNot(HaveOccurred(), "Writing kube object to temp file shouldn't fail")
	return oc("apply", "-f", tmpfile.Name())
}

func ocWaitFor(args ...string) string {
	return oc(append([]string{"wait", defaultOCWaitTimeout, "--for"}, args...)...)
}

func ocWaitLongFor(args ...string) string {
	return oc(append([]string{"wait", defaultOCLongWaitTimeout, "--for"}, args...)...)
}

// Will set up a scan with the given name and wait for it to be done.
// The scan will be done for the CIS benchmark.
func withCISScan(scan string) {
	By("Creating a ScanSettingBinding for this test")
	oc("compliance", "bind", "--name", scan, "profile/ocp4-cis")

	time.Sleep(5 * time.Second)
	ocWaitFor("condition=ready", "scansettingbinding", scan)

	By("Waiting for scan to be ready")
	ocWaitLongFor("condition=ready", "compliancesuite", scan)
}
