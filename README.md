# pocolm

---

## About the fork

The fork has some aditional example for Lithuanian LM. See [egs/dlkt](egs/dlkt).
The samples require SRILM for comparison and perplexity calculations.
To run the sample, configure [egs/dlkt/Makefile.options](egs/dlkt/Makefile.options):

```bash
cd egs/dlkt
#to build adaptation model using pocolm
# result will be saved in data/final/
make build-pocolm 
#to build adaptation model using srilm
make build-srilm
#to calculate perplexity
make calc-ppl-pocolm
make calc-ppl-srilm
```

---

Small language toolkit for creation, interpolation and pruning of ARPA language models

See [INSTALL.md](INSTALL.md) for installation instructions.

To use this toolkit: after installing, look at the example scripts for how to
run it.  egs/swbd/ and egs/swbd_fisher/ will show you how to run it on real
data.

[Motivation for this project](docs/motivation.md)
