# The Go test results as a Markdown job summary, from gotestsum's JSON (one test2json event per
# line, read with -s). $module (`go list -m`) is cut from package paths.
# Usage: jq -rs --arg module "$(go list -m)" -f .github/scripts/go-tests.jq .tmp/go-test.json
def done: .Action == "pass" or .Action == "fail" or .Action == "skip";
def short: sub("^" + $module + "/?"; "") | if . == "" then "." else . end;
def count($a): map(select(.Action == $a)) | length;

[.[] | select(.Test and done)] as $tests
| [.[] | select(.Test == null and done)] as $pkgs
# a failed test whose subtests all passed is where the failure is; its parents only repeat it
| [$tests[] | select(.Action == "fail")] as $failed
| [$failed[] | . as $f | select([$failed[] | select(.Package == $f.Package and (.Test | startswith($f.Test + "/")))] | length == 0)] as $leaves
| [.[] | select(.Action == "output")] as $out
| ($tests | count("pass")) as $passed | ($tests | count("fail")) as $nfailed | ($tests | count("skip")) as $skipped
| ([$pkgs[] | select(.Action == "fail")] | length) as $pkgfail
| "### Go tests\n\n"
  + (if $nfailed + $pkgfail == 0 then "✅" else "❌" end)
  + " **\($passed) passed**, \($nfailed) failed, \($skipped) skipped\n\n"
  + "| | Package | Passed | Failed | Skipped | Time |\n|---|---|--:|--:|--:|--:|\n"
  + ([$pkgs[] | . as $p
      | [$tests[] | select(.Package == $p.Package)] as $t
      | select(($t | length) > 0 or $p.Action == "fail")
      | "| \(if $p.Action == "fail" then "❌" elif ($t | length) == ($t | count("skip")) then "⏭️" else "✅" end) "
        + "| `\($p.Package | short)` | \($t | count("pass")) | \($t | count("fail")) | \($t | count("skip")) "
        + "| \($p.Elapsed // 0)s |"
    ] | join("\n"))
  + "\n"
  # the output of each failure: the lines of the test itself, or of the package when no test ran (a build error)
  + ([($leaves[] | {Package, Test}),
      ($pkgs[] | select(.Action == "fail") | . as $p
        | select([$tests[] | select(.Package == $p.Package)] | length == 0) | {Package, Test: null})]
     | map(. as $f
       | "\n<details><summary>❌ <code>\($f.Package | short)</code>\(if $f.Test then " " + $f.Test else " (no tests ran)" end)</summary>\n\n```\n"
         + ([$out[] | select(.Package == $f.Package and .Test == $f.Test) | .Output] | add // "" | rtrimstr("\n"))
         + "\n```\n\n</details>\n")
     | join(""))
