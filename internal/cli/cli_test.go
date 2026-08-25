package cli

import (
	"bytes"
	"strings"
	"testing"
)

const bowJSON = `{
  "name": "cubic-bow",
  "controlPoints": [[0,0],[0.45,1.4],[1.55,1.4],[2,0]],
  "offsetDistance": 0.25
}`

const cuspJSON = `{"controlPoints":[[0,0],[0,0],[0,0],[1,0]]}`

func run(t *testing.T, stdin string, args ...string) (int, string, string) {
	t.Helper()
	var out, errBuf bytes.Buffer
	env := Env{Stdin: strings.NewReader(stdin), Stdout: &out, Stderr: &errBuf}
	code := Run(args, env)
	return code, out.String(), errBuf.String()
}

func TestCLISampleBow(t *testing.T) {
	code, out, errS := run(t, bowJSON, "sample", "-")
	if code != ExitOK {
		t.Fatalf("sample exit = %d, stderr: %s", code, errS)
	}
	for _, want := range []string{"arc length", "kappa", "offset polyline", "cubic-bow"} {
		if !strings.Contains(out, want) {
			t.Errorf("sample output missing %q", want)
		}
	}
	if !strings.Contains(out, "1.545") {
		t.Errorf("sample output should contain relative length ~1.545, got:\n%s", out)
	}
}

func TestCLIArcLengthBow(t *testing.T) {
	code, out, errS := run(t, bowJSON, "arclength", "-")
	if code != ExitOK {
		t.Fatalf("arclength exit = %d, stderr: %s", code, errS)
	}
	for _, want := range []string{"adaptive Simpson", "gauss", "chord sum", "converged=true"} {
		if !strings.Contains(out, want) {
			t.Errorf("arclength output missing %q", want)
		}
	}
}

func TestCLICurvatureBow(t *testing.T) {
	code, out, errS := run(t, bowJSON, "curvature", "-", "-t", "0.5")
	if code != ExitOK {
		t.Fatalf("curvature exit = %d, stderr: %s", code, errS)
	}
	if !strings.Contains(out, "1.553937") {
		t.Errorf("curvature output missing expected kappa ~1.553937, got:\n%s", out)
	}
}

func TestCLIOffsetBow(t *testing.T) {
	code, out, errS := run(t, bowJSON, "offset", "-", "-d", "0.25", "-n", "8")
	if code != ExitOK {
		t.Fatalf("offset exit = %d, stderr: %s", code, errS)
	}
	for _, want := range []string{"offset distance", "rejected points", "0 of 9"} {
		if !strings.Contains(out, want) {
			t.Errorf("offset output missing %q", want)
		}
	}
}

func TestCLIValidateInvalidExit(t *testing.T) {
	three := `[[0,0],[1,1],[2,0]]`
	code, out, errS := run(t, three, "validate", "-")
	if code == ExitOK {
		t.Error("validate with 3 points must exit non-zero")
	}
	if !strings.Contains(errS, "need exactly 4 control points") {
		t.Errorf("stderr should mention 4 control points, got: %s", errS)
	}
	if out != "" {
		t.Errorf("stdout should be empty on failure, got: %s", out)
	}
}

func TestCLICuspCurvatureExit(t *testing.T) {
	code, _, errS := run(t, cuspJSON, "curvature", "-", "-t", "0")
	if code == ExitOK {
		t.Error("curvature at cusp must exit non-zero")
	}
	if !strings.Contains(errS, "curvature undefined") {
		t.Errorf("stderr should mention curvature undefined, got: %s", errS)
	}
}

func TestCLIMissingFile(t *testing.T) {
	code, _, errS := run(t, "", "sample", "no-such-file.json")
	if code == ExitOK {
		t.Error("missing file must exit non-zero")
	}
	if !strings.Contains(errS, "no such file") {
		t.Errorf("stderr should mention missing file, got: %s", errS)
	}
}

func TestCLIUnknownCommand(t *testing.T) {
	code, _, errS := run(t, "", "frobnicate")
	if code != ExitUsage {
		t.Errorf("unknown command exit = %d, want %d", code, ExitUsage)
	}
	if !strings.Contains(errS, "unknown command") {
		t.Errorf("stderr should mention unknown command, got: %s", errS)
	}
}

func TestCLICheckBow(t *testing.T) {
	code, out, errS := run(t, bowJSON, "check", "-")
	if code != ExitOK {
		t.Fatalf("check exit = %d, stderr: %s", code, errS)
	}
	if !strings.Contains(out, "5/5 applicable invariants hold") {
		t.Errorf("check output should report all applicable invariants, got:\n%s", out)
	}
}

func TestCLIValidateOK(t *testing.T) {
	code, out, errS := run(t, bowJSON, "validate", "-")
	if code != ExitOK {
		t.Fatalf("validate exit = %d, stderr: %s", code, errS)
	}
	if !strings.Contains(out, "OK: 4 control points") {
		t.Errorf("validate output should confirm 4 points, got: %s", out)
	}
}

func TestCLINoArgs(t *testing.T) {
	code, _, _ := run(t, "", "help")
	if code != ExitOK {
		t.Errorf("help exit = %d, want 0", code)
	}
	code2, _, _ := run(t, "")
	if code2 != ExitUsage {
		t.Errorf("no-args exit = %d, want %d", code2, ExitUsage)
	}
}
