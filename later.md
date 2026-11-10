# Later: stricter checks

Ideas to make the checks stricter, collected on 2026-11-10 while building the todo page.
`npm audit` stays at "high" on purpose: "moderate" would fail on 12 known advisories in dev tools
that have no fix.

## Possible, but costs code churn (measure first)

Turn each on temporarily, count what fails, and decide before changing anything.

1. **typescript-eslint `strictTypeChecked`** instead of `recommended`: many more type-aware rules
   (unnecessary conditions, unsafe `any`, non-null assertions). It would flag `result.data!` and
   `dropdown.querySelector(...)!`, so it means a cleanup pass.
