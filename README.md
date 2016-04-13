# sscompare

Bulk compare fuzzy hashes for file objects using spamsum algorithm. Unlike out of the box ssdeep will
compute everything in memory and on the fly enabling point and type bulk capability. 

Rationale for the use of Fuzzy Hashing in archives: https://gist.github.com/ross-spencer/4c2579f4ad2991785485 

Currently uses cgo version of ssdeep: https://godoc.org/github.com/dutchcoders/gossdeep

###Usage

      Usage:  sscompare [-fuzz] [-file1 ...]
      Usage:  sscompare [-fuzz] [-string1 ...]
      Usage:  sscompare [-compare] [-file1 ...] [-file2 ...]
      Usage:  sscompare [-compare] [-string1 ...] [-string2 ...]
      Usage:  sscompare [-compare] [-hash1 ...] [-hash2 ...]
      Usage:  sscompare [-compute] [-dir1 ...] [OPTIONAL] [-all]
      Output: [CSV] 'file1','hash'
      Output: [CSV] 'similarity','file1 | string1 | hash1','file2 | string1 | hash2',
      Usage of sscompare:
        -all
          	[Optional] Output all files, including zero matches and duplicates.
        -compare
          	Compare two hashes and return the percentage (%) familiarity.
        -compute
          	Compare all file hashes and output all comparisons for a given directory.
        -dir1 string
          	[Conditional] Directory to run a full 1:1 comparison against.
        -file1 string
          	[Conditional] File or string to generate and/or compare a hash for.
        -file2 string
          	[Conditional] File to compare file1 to.
        -fuzz
          	Generate a fuzzy hash for a file or string.
        -hash1 string
          	[Conditional] Hash to run a comparison against. The needle.
        -hash2 string
          	[Conditional] Hash to compare a hash1 to. The haystack.
        -string1 string
          	[Conditional] File or string to generate and/or compare a hash for.
        -string2 string
          	[Conditional] String to compare string1 to.

###Requirements

While this utilizes CGO there are still a handful of complications, including easy compilation on Windows.
Instructions for Windows have not yet been created. 

On Linux, compilation is much easier, but you might need libfyzzy, e.g. 

    sudo apt-get install libfuzzy2.2.13

Those using graphfuzzy.py (below) require [networkx](https://networkx.github.io/), I recommend using [PIP](https://pypi.python.org/pypi/pip). 

    pip install networkx

###GraphFuzzy see: ['graph-fuzzy/'](https://github.com/ross-spencer/sscompare/tree/master/graph-fuzzy)

Utilise the output from sscompare in a tool like [SocNetV](http://socnetv.sourceforge.net/) by creating GraphML

    python graphfuzzy.py --csv fuzzy-report.csv --score 0.5
  
    ** --score can be used to filter the report to return only values >= to the number (float) given

Usage: 

      usage: graphfuzzy.py [-h] [--csv CSV] [--score SCORE]

      Convert results of a fuzzy hash computation to a network graph, GraphML.
      Outputs using CSV filename with XML suffix.

      optional arguments:
        -h, --help            show this help message and exit
        --csv CSV, --results CSV
                              CSV export from sscompare tool.
        --score SCORE         Filter out results less than this match score.


![Network Graph Example](https://raw.githubusercontent.com/ross-spencer/rs-misc-scripts/master/anon-fuzzy-analysis.png)

Share your results with me! And let me know other tools that you've found useful for working with GraphML.

My first two results with SocNetV: [GovDocs and OPF Corpus](https://twitter.com/beet_keeper/status/719512360264683520)

###License

Copyright (c) 2016 Ross Spencer

This software is provided 'as-is', without any express or implied warranty. In no event will the authors be held liable for any damages arising from the use of this software.

Permission is granted to anyone to use this software for any purpose, including commercial applications, and to alter it and redistribute it freely, subject to the following restrictions:

The origin of this software must not be misrepresented; you must not claim that you wrote the original software. If you use this software in a product, an acknowledgment in the product documentation would be appreciated but is not required.

Altered source versions must be plainly marked as such, and must not be misrepresented as being the original software.

This notice may not be removed or altered from any source distribution.
