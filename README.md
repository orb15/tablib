# Tablib


## Overview
Tablib is short for "table library".  It is a library to execute an 'extensible random table generation engine' using a grammar and mechanisms similar to but different from a similar project in my repo (https://github.com/orb15/tabproc). I eventually found the ANTLR4-based approach in that project unsatisfactory because of the many complexities I encountered when trying to consume the IPP3 grammar. Despite this setback, the goal remains the same. I desire to replicate significant portions of the functionality present in the Windows Desktop Application [Inspiration Pad Pro3](http://www.nbos.com/products/inspiration-pad-pro), but move the execution engine off the Windows desktop and onto a linux server fronted with with REST services. This library will be the foundation of this implementation.

## Approach
Since my first implementation, I have switched languages (Java to Go) and most importantly, have given up on parity with the IPP3 grammar. The grammar is rich but still limited and I believe I can address both the complexity issue I encountered with the ANTLR4 approach as well as address several shortcomings in the IPP3 execution engine. The approach now is to use as much off-the-shelf technology as possible, building only what is needed and staying away from formal grammars in favor of implicit grammars defined contextually in YAML. Moving to YAML also allows for greater clarity (but less brevity) in certain table expressions and also provides free parsing and structural validation prior to any semantic validation and post-ingestion parsing, freeing me to focus on the important parts of the project.

## Example IPP3 Input
For an example of what this code will attempt to emulate, have a look at my [IPP3 project](https:/github.com/orb15/ipp3).  This project holds "code" for the IPP3 Windows Desktop client.
