package main

import (
   "os"
   "fmt"
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

//Generates hashes for two strings and compares those values
func CompareStrings(str1 string, str2 string) (result, error) {
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
   r.score, err = ssdeep.Compare(r.s1, r.s2)
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

//Generates hashes for two files and compares those values
func Comparefiles(file1 string, file2 string) (result, error) {
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

//Runs ssdeep compare for two pre-existing hash strings
func Comparehashes(hash1 string, hash2 string) (result, error) {
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


//Return the hash of a single string.
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

