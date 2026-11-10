# Go coverage per package as a Markdown job summary, from go-test-coverage's breakdown
# (`file;statements;covered` per line, its own numbers: excludes and coverage-ignore applied),
# against the thresholds in .testcoverage.yml.
# Usage: sort .tmp/coverage-go.txt | awk -F';' -f .github/scripts/go-coverage.awk
# Portable awk: also runs on mawk, Ubuntu's default.

function pct(c, t) { return t ? 100 * c / t : 100 }

BEGIN {
  # the `threshold:` block of .testcoverage.yml: `  file: 100` and so on
  while ((getline line < ".testcoverage.yml") > 0) {
    if (line ~ /^threshold:/) { on = 1; continue }
    if (on && line ~ /^[^ #]/) on = 0
    if (on && line ~ /^ +[a-z]+: *[0-9]+/) {
      split(line, kv, ":"); key = kv[1]; gsub(/ /, "", key); val = kv[2]; gsub(/[^0-9]/, "", val)
      thr[key] = val + 0
    }
  }
  close(".testcoverage.yml")
}

{
  pkg = $1; sub("/[^/]*$", "", pkg)
  if (!(pkg in total)) order[++n] = pkg
  total[pkg] += $2; covered[pkg] += $3; T += $2; C += $3
  if (pct($3, $2) < thr["file"]) low[++m] = sprintf("`%s` %.1f%%", $1, pct($3, $2))
}

END {
  printf "### Go coverage\n\n%s **%.1f%%** of %d statements (thresholds: total %d%%, package %d%%, file %d%%)\n\n",
    (pct(C, T) >= thr["total"] ? "✅" : "❌"), pct(C, T), T, thr["total"], thr["package"], thr["file"]
  print "| | Package | Statements | Covered | Coverage |"
  print "|---|---|--:|--:|--:|"
  for (i = 1; i <= n; i++) {
    p = order[i]; v = pct(covered[p], total[p])
    printf "| %s | `%s` | %d | %d | %.1f%% |\n", (v >= thr["package"] ? "✅" : "❌"), p, total[p], covered[p], v
  }
  if (m) { printf "\nFiles below %d%%:\n\n", thr["file"]; for (i = 1; i <= m; i++) print "- " low[i] }
}
