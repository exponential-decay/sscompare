package main

import (
   "io/ioutil"
   "crypto/sha1"
   "encoding/hex"
)

func hashfile(path string) (string, error) {
   hasher := sha1.New()
   s, err := ioutil.ReadFile(path)    
   hasher.Write(s)
   if err != nil {
      return "", err
   }
   return hex.EncodeToString(hasher.Sum(nil)), nil
}
