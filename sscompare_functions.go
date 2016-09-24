package main

import (
   "os"
   "fmt"
   "time"
   "strings"
   "github.com/dutchcoders/gossdeep"
)

var rescache bool = false

type result struct {
   paths       bool
   score       int
   s1          string   //path or hash
   s2          string
   strflag     bool     //standard string compare not equal
   shaflag     bool     //sha1 comparison not equal
}

func outputresults(r result) {
   if !r.paths {
      fmt.Fprintf(os.Stdout, "%d,%s,%s,%t\n", r.score, r.s1, r.s2, r.strflag) 
   } else {
      fmt.Fprintf(os.Stdout, "%d,%s,%s,%t,%t\n", r.score, r.s1, r.s2, r.strflag, r.shaflag) 
   }
}

func hashString(str string) (string, error) {
   hash, err := ssdeep.HashString(str)
   if err != nil {
      return "", err
   }
   return hash, nil
}

//Creates a fizzy hash for a single file
func createfilehash(path string) (string, error) {
   f, err := os.Open(path)
   defer f.Close()
   if err != nil {
      return "", err
   } 
   hash, err := ssdeep.HashFilename(path)   //confusing function title, not mine!
   if err != nil {
      return "", nil
   }
   return hash, nil
}

func compareStrings(str1 string, str2 string) (result, error) {
   var r result
   var err error
   r.s1, err = hashString(str1)
   if err != nil {
      return r, err
   }
   r.s2, err = hashString(str2)
   if err != nil {
      return r, err
   }
   r.score, err = ssdeep.Compare(hash1, hash2)
   if err != nil {
      return r, err
   }
   if r.score == 100 {     //100 spotted in the wild for non-identifcal files
      if strings.Compare(r.s1, r.s2) != 0 {
         r.strflag = true
      }
   }
   return r, nil
}

func comparefiles(file1 string, file2 string) (result, error) {
   var r result
   r.paths = true
   var err error
   f1, _ := fileExists(file1)
   f2, _ := fileExists(file2)
   if !f1 || !f2 {
      return r, fmt.Errorf("Warning: Cannot find file.\n")
   }
   r.s1, err = createfilehash(file1)
   if err != nil {
      return r, err
   }
   r.s2, err = createfilehash(file2)
   if err != nil {
      return r, err
   }
   r.score, err = ssdeep.Compare(r.s1, r.s2)
   if err != nil {
      return r, err
   }
   if r.score == 100 {     //100 spotted in the wild for non-identifcal files
      if strings.Compare(r.s1, r.s2) != 0 {
         r.strflag = true
      }
      shaval1, err := hashfile(file1)
      if err != nil {
         return r, err_sha1_file1
      } 
      shaval2, err := hashfile(file2)
      if err != nil {
         return r, err_sha1_file2
      }
      if strings.Compare(shaval1, shaval2) != 0 {
         r.shaflag = true
      }
   }
   return r, nil 
}

func comparehashes(hash1 string, hash2 string) (result, error) {
   var r result
   var err error
   r.s1 = hash1
   r.s2 = hash2
   r.score, err = ssdeep.Compare(hash1, hash2)
   if err != nil {
      return r, err_sscomp  
   }
   if r.score == 100 {     //100 spotted in the wild for non-identifcal files
      if strings.Compare(r.s1, r.s2) != 0 {
         r.strflag = true
      }
   }
   return r, nil
}

func newVariant(str1 string, str2 string, results_cache [][]string) bool {
   add := false   
   exist1 := false
   exist2 := false
   for idx, _ := range results_cache {
      if len(results_cache[idx]) != 0 {
         file1 = results_cache[idx][0]
         file2 = results_cache[idx][1]
         if str1 == file1 || str1 == file2 {
            exist1 = true
         }
         if str2 == file1 || str2 == file2 {
            exist2 = true
            break
         }
      }
   }
   if exist1 != true || exist2 != true {   
      add = true  
   }
   return add
}

func handleComputeResults(score int, hash1 []string, hash2 []string, all bool) int {
   
   added := 0
   hfile1 := hash1[1]
   hfile2 := hash2[1]

   if all != true {
      row := []string{hfile1, hfile2}
      if rescache == false {
         results_cache = make([][]string, 512)
         results_cache = append(results_cache, row)
         fmt.Fprintf(os.Stdout, "%d,%s,%s\n", score, hfile1, hfile2)
         added = 1
         rescache = true
      } else {

         if newVariant(hfile1, hfile2, results_cache) {
            results_cache = append(results_cache, row)
            fmt.Fprintf(os.Stdout, "%d,%s,%s\n", score, hfile1, hfile2)
            added = 1
         }
      }
   } else {
      fmt.Fprintf(os.Stdout, "%d,%s,%s\n", score, hfile1, hfile2)
      added = 1
   }
   return added
}

func generateComparisonTable(hashes [][]string, all bool) {
   fmt.Fprintln(os.Stdout, "score,source,target")   
   total := 0
   x := len(hashes)
   for hash, _ := range hashes {
      x = x - 1
      if len(hashes[hash]) > 1 {    //we have (x) block slice, we may have empties
         hash1 := hashes[hash][0]
         found := false
         for h, _ := range hashes {
            if len(hashes[h]) > 1 {    //preferable to delete empty slices? 
               hash2 = hashes[h][0]
               score, err := ssdeep.Compare(hash1, hash2)
               if err == nil {
                  if score == 100 && found == false { //ignore first identical (itself)
                     found = true
                  } else {
                     if score != 0 {
                        total += handleComputeResults(score, hashes[hash], hashes[h], all)
                     }  //N.B. can opt to log zeroes, if we really care
                  }
               }         
            }
         }
      }
   }
   elapsed := time.Since(start)
   fmt.Fprintln(os.Stderr, total, elapsed)
}
