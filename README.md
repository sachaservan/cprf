# Constrained PRFs for Inner Product Predicates

**Paper:** https://eprint.iacr.org/2024/TBD

| **Code organization** ||
| :--- | :---|
| [ro-cprf/](ro-cprf/) | Random oracle based CPRF construction|
| [ddh-cprf/](ddh-cprf/) | DDH (Naor-Reingold) based CPRF construction|
| [owf-cprf/](owf-cprf/) | One-way function based CPRF construction|


## Running benchmarks
1. Go to sub-directory containing the implementation of interest.
2. Run ```go test -bench=.``` 
