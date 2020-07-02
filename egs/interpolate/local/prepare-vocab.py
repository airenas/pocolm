import argparse
import sys

def print_word(args, word, wc):
    if args.index:
        print("%s %d" % (word, wc))
    else:
        print("%s" % word)    

def main(argv):
    parser = argparse.ArgumentParser(description="This script makes vocab from text file",
                                     epilog="E.g. cat input.txt | " + sys.argv[0] + " > result.txt",
                                     formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    parser.add_argument("--add-unk", default='', type=str, help="Unknown word to add if any")
    parser.add_argument("--add-eps", default='', type=str, help="Epsilon word to add if any")
    parser.add_argument("--all-words", action='store_true', help="Do not skip words like '<xxx>'")
    parser.add_argument("--index", action='store_true', help="Add word index next to word")
    args = parser.parse_args(args=argv)

    print("Starting", file=sys.stderr)
    words = set()
    lc = 0
    wc = 0
    sc = 0
    for line in sys.stdin:
        lc += 1
        ws = line.split()
        for w in ws:
            wc += 1
            if args.all_words or (not (w.startswith("<") and w.endswith(">"))):
                words.add(w)
            else:
                sc += 1    
    print("Read %d lines, %d words, %d distinct, %d skipped." % (lc, wc, len(words), sc), file=sys.stderr)

    wc = 0
    if args.add_eps != '':
        print("Add eps: %s" % args.add_eps, file=sys.stderr)
        print_word(args, args.add_eps, wc)
        wc += 1
    for w in sorted(words):
        print_word(args, w, wc)
        wc += 1
    if args.add_unk != '' and args.add_unk not in words:
        print("Add unk: %s" % args.add_unk, file=sys.stderr)
        print_word(args, args.add_unk, wc)
        wc += 1

    print("Wrote %d words" % wc, file=sys.stderr)
    print("Done", file=sys.stderr)


if __name__ == "__main__":
    main(sys.argv[1:])
