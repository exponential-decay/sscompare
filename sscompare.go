package main

import (
	"os"
	"fmt"
   "flag"
   "time"
   "path/filepath"
)

var (
   fuzzy, compare, compute, all, storeHashes bool
   file1, file2, hash1, hash2, string1, string2, dir string
   pathhashes []pathhash
   start time.Time
)

func init() {
   flag.BoolVar(&fuzzy, "fuzzy", false, "Generate a fuzzy hash for a file or string.")
   flag.BoolVar(&compare, "compare", false, "Compare two hashes and return the percentage (%) familiarity.")
   flag.BoolVar(&compute, "compute", false, "Compare all file hashes and output all comparisons for a given directory.")
   flag.StringVar(&file1, "file1", "false", "[Conditional] File or string to generate and/or compare a hash for.")
   flag.StringVar(&file2, "file2", "false", "[Conditional] File to compare file1 to.")
   flag.StringVar(&string1, "string1", "false", "[Conditional] File or string to generate and/or compare a hash for.")
   flag.StringVar(&string2, "string2", "false", "[Conditional] String to compare string1 to.")
   flag.StringVar(&hash1, "hash1", "false", "[Conditional] Hash to run a comparison against. The needle.")
   flag.StringVar(&hash2, "hash2", "false", "[Conditional] Hash to compare a hash1 to. The haystack.")
   flag.StringVar(&dir, "dir", "false", "[Conditional] Directory to run a full 1:1 comparison against.")
   flag.BoolVar(&all, "all", false, "[Optional] Output all files, including zero matches and duplicates.")
}

func fileExists(path string) (bool, os.FileInfo) {
   f, err := os.Open(path)
   defer f.Close()
   if err != nil {
      fmt.Fprintln(os.Stderr, "Error:", err)
      os.Exit(1)
   } 
   fi, err := f.Stat()
   if err != nil {
      fmt.Fprintln(os.Stderr, "Error:", err)
      os.Exit(1)
   }
   return true, fi
}

func computeall(path string, all bool) { 
   f1, fi := fileExists(path)
   if f1 {
      switch mode := fi.Mode(); {
      case mode.IsDir():
         //store hashes for comparison
         storeHashes = true         
         //walk the dir
         filepath.Walk(path, readFile)
      default:
         fmt.Fprintf(os.Stderr, "Warning: Cannot compute all values on a non-directory.\n")
         os.Exit(1)
      }
   }
   //if we complete processing of a directory, generate comparison table
   if len(pathhashes) > 0 {
      createcomparisontable(pathhashes, all)
   }
}

func readFile(path string, fi os.FileInfo, err error) error {
   f1, _ := fileExists(path) 
   if f1 {
      switch mode := fi.Mode(); {
      case mode.IsRegular():
         hash, err := createfilehash(path)
         if err != nil {
            return err
         }
         if storeHashes == true {
            pathhashes = append(pathhashes, pathhash{hash, path})
         } else {
            //output for $sscompare -fuzzy -file1 <filepath.file>
            fmt.Fprintf(os.Stdout, "%s,%s\n", path, hash)
         }
      }
   }
   return nil
}

func main() {

   flag.Parse()

   if flag.NFlag() < 2 {    // can access args w/ len(os.Args[1:]) too
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-fuzzy] [-file1 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-fuzzy] [-string1 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compare] [-file1 ...] [-file2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compare] [-string1 ...] [-string2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compare] [-hash1 ...] [-hash2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compute] [-dir ...] [OPTIONAL] [-all]")
      fmt.Fprintln(os.Stderr, "Output: [CSV] 'file1','hash'")
      fmt.Fprintln(os.Stderr, "Output: [CSV] 'score','hash1','hash2','string compare fail'")
      fmt.Fprintln(os.Stderr, "Output: [CSV] 'score','file1','file2','string compare fail','sha1 compare fail'")
      flag.Usage()
      os.Exit(0)
   }
   
   //create fuzzy hashes for a file or a string
   if (fuzzy == true && file1 != "false") {
      err := filepath.Walk(file1, readFile)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error: %s\n", err)
         os.Exit(1)
      }
   }
   if (fuzzy == true && string1 != "false") {
      hash, err := hashString(string1)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error: %s\n", err)
         os.Exit(1)
      }
      fmt.Println(hash)
   }

   //compare fuzzy hashes for a file or a string
   if (compare == true && file1 != "false" && file2 != "false") {
      r, err := Comparefiles(file1, file2)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error: %s\n", err)
         os.Exit(1)
      }
      outputresults(r)
   }

   //compare two strings and output comparison result
   if (compare == true && string1 != "false" && string2 != "false") {
      r, err := CompareStrings(string1, string2)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error: %s\n", err)
         os.Exit(1)
      }
      outputresults(r)
   }

   //compare two hashes and output similarity
   if (compare == true && hash1 != "false" && hash2 != "false") {
      r, err := Comparehashes(hash1, hash2)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error: %s\n", err)
         os.Exit(1)
      }
      outputresults(r)
   } 

   //compute all values (n*n) for a dir
   if (compute == true && dir != "false") {
      start = time.Now()
      computeall(dir, all)
      elapsed := time.Since(start)
      fmt.Fprintln(os.Stderr, elapsed)
   }
}

