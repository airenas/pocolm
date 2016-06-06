#!/usr/bin/env python

# we're using python 3.x style print but want it to work in python 2.x,
from __future__ import print_function
import re, os, argparse, sys, math, warnings, subprocess, shutil
from collections import defaultdict

parser = argparse.ArgumentParser(description="This script takes an lm-dir, as produced by make_lm_dir.py, "
                                 "that should not have the counts split up into pieces, and it prunes "
                                 "the counts and writes out to a new lm-dir.")

parser.add_argument("--steps", type=str,
                    default='prune*4.0 EM EM EM prune*2.0 EM EM prune*1.0 EM EM ',
                    help='This string specifies a sequence of steps in the pruning sequence.'
                    'prune*X, with X >= 1.0, tells it to prune with X times the threshold '
                    'specified with the --threshold option.  EM specifies one iteration of '
                    'E-M on the model. ')
parser.add_argument("--remove-zeros", type=str, choices=['true','false'],
                    default='true', help='Set this to false to disable an optimization. '
                    'Only useful for debugging purposes.')
parser.add_argument("lm_dir_in",
                    help="Source director, for the input language model.")
parser.add_argument("threshold", type=float,
                    help="Threshold for pruning, e.g. 0.5, 1.0, 2.0, 4.0.... "
                    "larger threshold will give you more highly-pruned models."
                    "Threshold is interpreted as entropy-change times overall "
                    "weighted data count, for each parameter.  It should be "
                    "larger if you have more data, assuming you want the "
                    "same-sized models.")
parser.add_argument("lm_dir_out",
                    help="Output directory where the language model is created.")


args = parser.parse_args()

# Add the script dir and the src dir to the path.
os.environ['PATH'] = (os.environ['PATH'] + os.pathsep +
                      os.path.abspath(os.path.dirname(sys.argv[0])) + os.pathsep +
                      os.path.abspath(os.path.dirname(sys.argv[0])) + "/../src");

if os.system("validate_lm_dir.py " + args.lm_dir_in) != 0:
    sys.exit("prune_lm_dir.py: failed to validate input LM-dir")

if os.path.exists(args.lm_dir_in + "/num_splits"):
    sys.exit("prune_lm_dir.py: input LM-dir is split ({0}/num_splits exists). "
             "This script currently only works with non-split LM-dirs.".format(
            args.lm_dir_in))

if args.threshold <= 0.0:
    sys.exit("prune_lm_dir.py: --threshold must be positive: got " + str(args.threshold))

work_dir = args.lm_dir_out + "/work"

steps = args.steps.split()

if len(steps) == 0:
    sys.exit("prune_lm_dir.py: 'steps' cannot be empty.")

# returns num-words in this lm-dir.
def GetNumWords(lm_dir_in):
    command = "tail -n 1 {0}/words.txt".format(lm_dir_in)
    line = subprocess.check_output(command, shell = True)
    try:
        a = line.split()
        assert len(a) == 2
        ans = int(a[1])
    except:
        sys.exit("prune_lm_dir: error: unexpected output '{0}' from command {1}".format(
                line, command))
    return ans

def GetNgramOrder(lm_dir_in):
    f = open(lm_dir_in + "/ngram_order");
    return int(f.readline())

def RunCommand(command):
    # print the command for logging
    print(command, file=sys.stderr)
    if os.system(command) != 0:
        sys.exit("get_objf_and_derivs_split.py: error running command: " + command)




# This script creates work/protected.all (listing protected
# counts which may not be removed); it requires work/float.all
# to exist.
def CreateProtectedCounts(work):
    command = ('float-counts-to-histories <{0}/float.all | LC_ALL=C sort |'
               ' histories-to-null-counts >{0}/protected.all'.format(work))
    RunCommand(command)


def SoftLink(src, dest):
    if os.path.exists(dest):
        os.remove(dest)
    try:
        os.symlink(os.path.abspath(src), dest)
    except:
        sys.exit("prune_lm_dir.py: error linking {0} to {1}".format(src, dest))

def CreateInitialWorkDir():
    # Creates float.all, stats.all, and protected.all in work_dir/iter
    work0dir = work_dir + "/iter0"
    # create float.all
    if not os.path.isdir(work0dir):
        os.makedirs(work0dir)
    SoftLink(args.lm_dir_in + "/float.all", work0dir + "/float.all")

    # create protected.all
    CreateProtectedCounts(work0dir)

    # create stats.all
    # e.g. command = 'float-counts-to-float-stats 20000 foo/work/iter0/stats.1 '
    #                'foo/work/iter0/stats.2 <foo/work/iter0/float.all'
    command = ("float-counts-to-float-stats {0} ".format(num_words) +
               ' '.join([ "{0}/stats.{1}".format(work0dir, n)
                          for n in range(1, ngram_order + 1) ]) +
               " <{0}/float.all".format(work0dir))
    RunCommand(command)


def RunPruneStep(work_in, work_out, threshold):
    # set float_star = 'work_out/float.1 work_out/float.2 ...'
    float_star = " ".join([ '{0}/float.{1}'.format(work_out, n)
                            for n in range(1, ngram_order + 1) ])
    # create work_out/float.{1,2,..}
    command = ("float-counts-prune {threshold} {num_words} {work_in}/float.all "
               "{work_in}/protected.all ".format(threshold = threshold,
                                                 num_words = num_words,
                                                 work_in = work_in) +
               float_star)
    RunCommand(command)
    if args.remove_zeros == 'false':
        # create work_out/float.all.
        command = 'merge-float-counts {0} >{1}/float.all'.format(float_star, work_out)
        RunCommand(command)
        command = 'rm ' + float_star
        RunCommand(command)
        # soft-link work_out/stats.all to work_in/stats.all
        SoftLink(work_in + "/stats.all",
                 work_out + "/stats.all")
    else:
        # in this case we pipe the output of merge-float-counts into
        # float-counts-stats-remove-zeros.
        # set stats_star = 'work_out/stats.1 work_out/stats.2 ..'
        stats_star = " ".join([ '{0}/stats.{1}'.format(work_out, n)
                                for n in range(1, ngram_order + 1) ])

        command = ('merge-float-counts {float_star} | float-counts-stats-remove-zeros /dev/stdin '
                   '{work_in}/stats.all {work_out}/float.all {stats_star}'.format(
                float_star = float_star, work_in = work_in, work_out = work_out,
                stats_star = stats_star))
        RunCommand(command)

        # create work_out/stats.all
        command = 'merge-float-counts {0} >{1}/stats.all'.format(stats_star, work_out)
        RunCommand(command)
        command = 'rm ' + float_star + ' ' + stats_star
        RunCommand(command)

    # create work_out/protected.all
    CreateProtectedCounts(work_out)



def RunEmStep(work_in, work_out):
    # set float_star = 'work_out/float.1 work_out/float.2 ...'
    float_star = " ".join([ '{0}/float.{1}'.format(work_out, n)
                            for n in range(1, ngram_order + 1) ])

    command = ('float-counts-estimate {num_words} {work_in}/float.all {work_in}/stats.all '
               '{float_star}'.format(num_words = num_words, work_in = work_in,
                                     float_star = float_star))
    RunCommand(command)
    command = 'merge-float-counts {0} >{1}/float.all'.format(float_star, work_out)
    RunCommand(command)
    command = 'rm ' + float_star
    RunCommand(command)
    # soft-link work_out/stats.all to work_in/stats.all
    SoftLink(work_in + "/stats.all",
             work_out + "/stats.all")
    # soft-link work_out/protected.all to work_in/protected.all
    SoftLink(work_in + "/protected.all",
             work_out + "/protected.all")


# runs one of the numbered steps.  step_number >= 0 is the number of the work
# directory we'll get the input from (the output will be that plus one).
def RunStep(step_number):
    work_in = work_dir + "/iter" + str(step_number)
    work_out = work_dir + "/iter" + str(step_number + 1)
    if not os.path.isdir(work_out):
        os.makedirs(work_out)
    step_text = steps[step_number]
    if step_text[0:6] == 'prune*':
        try:
            scale = float(step_text[6:])
            assert scale >= 1.0
        except:
            sys.exit("prune_lm_dir.py: invalid step (wrong --steps "
                     "option): '{0}'".format(step_text))
        RunPruneStep(work_in, work_out, args.threshold * scale)

    elif step_text == 'EM':
        RunEmStep(work_in, work_out)
    else:
        sys.exit("prune_lm_dir.py: invalid step (wrong --steps "
                 "option): '{0}'".format(step_text))


def FinalizeOutput(final_work_out):
    try:
        shutil.move(final_work_out + "/float.all",
                    args.lm_dir_out + "/float.all")
    except:
        sys.exit("prune_lm_dir.py: error moving {0}/float.all to {1}/float.all".format(
                final_work_out, args.lm_dir_out))
    f = open(args.lm_dir_out + "/was_pruned")
    print("true", file=f)
    f.close()
    for f in [ 'names', 'words.txt', 'ngram_order', 'metaparameters' ]:
        try:
            shutil.copy(args.lm_dir_in + "/" + f,
                        args.lm_dir_out + "/" + f)
        except:
            sys.exit("prune_lm_dir.py: error copying {0}/{1} to {2}/{1}".format(
                    args.lm_dir_in, f, args.lm_dir_out))


if not os.path.isdir(work_dir):
    try:
        os.makedirs(work_dir)
    except:
        sys.exit("prune_lm_dir.py: error creating directory " + work_dir)

num_words = GetNumWords(args.lm_dir_in)
ngram_order = GetNgramOrder(args.lm_dir_in)

CreateInitialWorkDir()
for step in range(len(steps)):
    RunStep(step)
FinalizeOutput(work_dir + "/iter" + str(len(step) + 1))


