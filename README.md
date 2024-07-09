# Constrained PRFs for Inner Product Predicates

Prototype implementation of constrained PRFs for inner product predicates from [this paper](https://eprint.iacr.org/2024/58). 

| **Code organization** ||
| :--- | :---|
| [ro-cprf/](ro-cprf/) | Random oracle based CPRF construction|
| [ddh-cprf/](ddh-cprf/) | DDH (Naor-Reingold) based CPRF construction|
| [owf-cprf/](owf-cprf/) | One-way function based CPRF construction|


## Running benchmarks
1. Go to sub-directory containing the implementation of interest.
2. Run ```go test -bench=.```

## TODOs (nice-to-haves)
- [ ] implement the VDLPN-based construction from [the paper](https://eprint.iacr.org/2024/58). 
- [ ] improve performance of the OWF-based construction by using a faster universal hashing technique
- [ ] optimize (lots of room for possible optimizations currently left on the table)

# Citations
```
@misc{cprfs,
      author = {Sacha Servan-Schreiber},
      title = {Constrained Pseudorandom Functions for Inner-Product Predicates from Weaker Assumptions},
      howpublished = {Cryptology ePrint Archive, Paper 2024/058},
      year = {2024},
      note = {\url{https://eprint.iacr.org/2024/058}},
      url = {https://eprint.iacr.org/2024/058}
}
```

## ⚠️ Important Warning

<b>This implementation is intended for _research purposes only_. The code has NOT been vetted by security experts.
As such, no portion of the code should be used in any real-world or production setting!</b>
