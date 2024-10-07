# Constrained PRFs for Inner Product Predicates

This repository contains implementations of Constrained Pseudorandom Functions (CPRFs) for inner product predicates from [this paper](https://eprint.iacr.org/2024/58) (to appear at AsiaCrypt 2024).
It includes two different constructions: one based on random oracles and another based on the Decisional Diffie-Hellman (DDH) assumption. The implementation can be used to reproduce Tables 2 & 3 from the paper. 

## Code Organization

| Directory | Description |
| :--- | :--- |
| [ro-cprf/](ro-cprf/) | Random oracle based CPRF construction |
| [ddh-cprf/](ddh-cprf/) | DDH (Naor-Reingold) based CPRF construction |

## Prerequisites

- Go (version 1.20 or later)

## Running Benchmarks

To run benchmarks for each implementation:

1. Choose the implementation you want to benchmark:
   - Random Oracle based: `cd ro-cprf`
   - DDH based: `cd ddh-cprf`

2. Run the benchmarks:
   ```
   go test -bench=.
   ```

## Interpreting the Results

The benchmark results are presented in the following format:
```
BenchmarkEval/length=10-10                555568              2150 ns/op
BenchmarkEval/length=50-10                118142              9941 ns/op
BenchmarkEval/length=100-10                59997             19887 ns/op
BenchmarkEval/length=500-10                12088             97314 ns/op
BenchmarkEval/length=1000-10                5764            196743 ns/op
```

- `length=X`: Indicates the vector length used in the CPRF evaluation.
- `-Y`: Indicates number of benchmark iterations.
- `Z ns/op`: Represents the average time in nanoseconds it takes for a single CPRF evaluation operation.

For example, `BenchmarkEval/length=100-10 60027 19896 ns/op` means:
- Vector length: 100
- Benchmark iterations: 10
- Average time per evaluation: 19,896 nanoseconds (about 0.02 milliseconds)

## Future Improvements

- [ ] Implement the VDLPN-based construction from [the paper](https://eprint.iacr.org/2024/58).
- [ ] Optimize implementations (there's room for performance improvements).

## Acknowledgements
We thank [Maxime Bombar](https://github.com/mbombar) for reviewing and providing feedback on the code.

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

**This implementation is intended for _research purposes only_. The code has NOT been vetted by security experts. As such, no portion of the code should be used in any real-world or production setting!**
