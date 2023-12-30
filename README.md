# Constrained PRFs for Inner Product Predicates

Prototype implementation of constrained PRFs for inner product predicates from [this paper](https://eprint.iacr.org/2024/TBD). 

| **Code organization** ||
| :--- | :---|
| [ro-cprf/](ro-cprf/) | Random oracle based CPRF construction|
| [ddh-cprf/](ddh-cprf/) | DDH (Naor-Reingold) based CPRF construction|
| [owf-cprf/](owf-cprf/) | One-way function based CPRF construction|


## Running benchmarks
1. Go to sub-directory containing the implementation of interest.
2. Run ```go test -bench=.```

## TODOs (nice-to-haves)
- [ ] implement the VDLPN-based construction from [the paper](https://eprint.iacr.org/2024/TBD). 
- [ ] improve performance of the OWF-based construction by using a faster universal hashing technique
- [ ] optimize (lots of room for possible optimizations currently left on the table)


## ⚠️ Important Warning

<b>This implementation of is intended for _research purposes only_. The code has NOT been vetted by security experts.
As such, no portion of the code should be used in any real-world or production setting!</b>

## License

Copyright © 2024 Sacha Servan-Schreiber

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
