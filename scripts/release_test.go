package scripts

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// run runs a script of this directory in dir and returns what it prints.
func run(t *testing.T, dir, script string, args ...string) string {
	t.Helper()
	path, err := filepath.Abs(script)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("bash", append([]string{path}, args...)...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v: %v\n%s", script, args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// git runs git in dir.
func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "tag.gpgSign=false", "-c", "commit.gpgSign=false"}, args...)...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// fixture is a repository with one commit per tag, tagged in push order; each commit's subject is
// "change for <tag>".
func fixture(t *testing.T, tags ...string) string {
	t.Helper()
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	git(t, dir, "commit", "-q", "--allow-empty", "-m", "first")
	for _, tag := range tags {
		git(t, dir, "commit", "-q", "--allow-empty", "-m", "change for "+tag)
		git(t, dir, "tag", tag)
	}
	return dir
}

func TestReleaseTags(t *testing.T) {
	cases := []struct {
		name string
		tags []string
		ref  string
		want string
	}{
		{
			name: "first release takes its line",
			tags: []string{"v0.1.0"},
			ref:  "v0.1.0",
			want: "0.1.0 0.1",
		},
		{
			// The issue's acceptance: v0.1.1 moves 0.1, and not an older line.
			name: "a patch moves its own line and no other",
			tags: []string{"v0.0.3", "v0.1.0", "v0.1.1"},
			ref:  "v0.1.1",
			want: "0.1.1 0.1",
		},
		{
			name: "a patch on an older line moves only that line",
			tags: []string{"v0.1.0", "v0.2.0", "v0.1.1"},
			ref:  "v0.1.1",
			want: "0.1.1 0.1",
		},
		{
			// Parallel or re-run tags: v0.1.1's run sees v0.1.2 already pushed.
			name: "an older patch of the line never moves it back",
			tags: []string{"v0.1.1", "v0.1.2"},
			ref:  "v0.1.1",
			want: "0.1.1",
		},
		{
			name: "numeric, not lexical, order",
			tags: []string{"v0.1.9", "v0.1.10"},
			ref:  "v0.1.10",
			want: "0.1.10 0.1",
		},
		{
			name: "a pre-release gets only its own tag",
			tags: []string{"v0.1.1", "v0.1.2-rc.1"},
			ref:  "v0.1.2-rc.1",
			want: "0.1.2-rc.1",
		},
		{
			name: "a pre-release does not hold back the line",
			tags: []string{"v0.1.2-rc.1", "v0.1.1"},
			ref:  "v0.1.1",
			want: "0.1.1 0.1",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := fixture(t, c.tags...)
			if got := run(t, dir, "release-tags.sh", c.ref); got != c.want {
				t.Errorf("release-tags.sh %s = %q, want %q", c.ref, got, c.want)
			}
		})
	}
}

func TestReleaseTagsRefusesNonVersion(t *testing.T) {
	dir := fixture(t, "v0.1")
	path, _ := filepath.Abs("release-tags.sh")
	cmd := exec.Command("bash", path, "v0.1")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("release-tags.sh v0.1 succeeded: %s", out)
	}
}

func TestReleaseNotes(t *testing.T) {
	cases := []struct {
		name string
		tags []string
		ref  string
		want []string // the change lines, in order
		not  []string
	}{
		{
			name: "first release lists every commit",
			tags: []string{"v0.1.0"},
			ref:  "v0.1.0",
			want: []string{"change for v0.1.0", "first"},
		},
		{
			name: "a release lists the commits since the previous release",
			tags: []string{"v0.1.0", "v0.1.1"},
			ref:  "v0.1.1",
			want: []string{"change for v0.1.1"},
			not:  []string{"change for v0.1.0", "first"},
		},
		{
			name: "a pre-release is not a previous release",
			tags: []string{"v0.1.0", "v0.1.1-rc.1", "v0.1.1"},
			ref:  "v0.1.1",
			want: []string{"change for v0.1.1", "change for v0.1.1-rc.1"},
			not:  []string{"change for v0.1.0"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := fixture(t, c.tags...)
			got := run(t, dir, "release-notes.sh", c.ref)
			at := 0
			for _, w := range c.want {
				i := strings.Index(got[at:], "- "+w+" (")
				if i < 0 {
					t.Fatalf("notes for %s lack %q after offset %d:\n%s", c.ref, w, at, got)
				}
				at += i + 1
			}
			for _, n := range c.not {
				if strings.Contains(got, "- "+n+" (") {
					t.Errorf("notes for %s list %q:\n%s", c.ref, n, got)
				}
			}
		})
	}
}

// The patch on an older line is diffed against that line, not against the newest release.
func TestReleaseNotesOlderLine(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q", "-b", "main")
	git(t, dir, "commit", "-q", "--allow-empty", "-m", "first")
	git(t, dir, "tag", "v0.1.0")
	git(t, dir, "checkout", "-q", "-b", "release/0.1")
	git(t, dir, "commit", "-q", "--allow-empty", "-m", "fix on 0.1")
	git(t, dir, "checkout", "-q", "main")
	git(t, dir, "commit", "-q", "--allow-empty", "-m", "feature on main")
	git(t, dir, "tag", "v0.2.0")
	git(t, dir, "tag", "v0.1.1", "release/0.1")
	got := run(t, dir, "release-notes.sh", "v0.1.1")
	if !strings.Contains(got, "- fix on 0.1 (") || strings.Contains(got, "feature on main") || strings.Contains(got, "- first (") {
		t.Errorf("notes for v0.1.1:\n%s", got)
	}
}
