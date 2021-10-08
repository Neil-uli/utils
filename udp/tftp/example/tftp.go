package main

import (
	"flag"
	"io/ioutil"
	"log"
    "os"

	"github.com/ulicod3/utils/udp/tftp"
)

var (
    address = flag.String("a", "127.0.0.1:69", "listen address")
    payload = flag.String("p", "payload.jpeg", "file to serve to clients")
    help = flag.Bool("h, --help", false, "Show help")
)

func main() {
   flag.Parse()

   if *help {
        flag.PrintDefaults()
        os.Exit(0) 
   }

   p, err := ioutil.ReadFile(*payload)
   if err != nil { log.Fatal(err) }
    
   s := tftp.Server{Payload: p}
   log.Fatal(s.ListenAndServe(*address))
}

