// lastro - a proxy for ruining connections
//
// Copyright (c) 2014, zahpee
// All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// 
// * Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
// 
// * Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
// 
// * Neither the name of lastro nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
    "time"
    "net"
    "log"
    "flag"
    "sync"
)

// address to bind proxy, it can be a port in the format :1212
var bindAddr = flag.String("bind", ":1212", "address to bind proxy")
// target address that should be accessed
var targetAddr = flag.String("target", "localhost:80", "target address example localhost:80")
// sleep time in ms between socket reads
var sleepMillis = flag.Int("sleep", 0, "sleep between sockreads in ms")

// Copy from one stream to another until error is returned.
func streamCopy(from, to net.Conn, wg *sync.WaitGroup) {
    defer wg.Done()
    buff := make([]byte, 512)
    for {
        if *sleepMillis > 0 {
            time.Sleep(time.Duration(*sleepMillis) * time.Millisecond)
        }
        n, err := from.Read(buff)
        if err != nil {
            log.Printf("Error reading. from_remote= %s err= %s",
                from.RemoteAddr(), err.Error())
            to.Close()
            return
        }

        _, err = to.Write(buff[0:n])
        if err != nil {
            log.Printf("Error writing. to_remote= %s err= %s",
                to.RemoteAddr(), err.Error())
            return
        }
    }
}

// Create a slow proxy for the accepted connection. This function opens
// a new connection with the target server and copies data from one
// stream to another.
func slowProxy(conn net.Conn) {
    // open new connection with target
    targetConn, err := net.Dial("tcp", *targetAddr)
    if err != nil {
        log.Printf("Could not connect to target= %s", *targetAddr)
        return
    }

    var wg sync.WaitGroup
    wg.Add(2)
    go streamCopy(conn, targetConn, &wg)
    go streamCopy(targetConn, conn, &wg)
    wg.Wait()
    log.Printf("Finished one slow tunnel from_remote= %s target_remote= %s",
        conn.RemoteAddr(), targetConn.RemoteAddr())
}

func main() {
    flag.Parse()

    log.Printf("============================================")
    log.Printf("Starting lastro")
    log.Printf("Binding addr: %s", *bindAddr)
    log.Printf("Target addr: %s", *targetAddr)
    log.Printf("Sleep between sockread: %dms", *sleepMillis)
    log.Printf("============================================")

    ln, err := net.Listen("tcp", *bindAddr)
    if err != nil {
        log.Printf("Could no bind proxy: %s", *bindAddr)
        return
    }

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("Error accepting connection: %s", err.Error())
            return
        }
        log.Printf("Accepted local= %s remote= %s", conn.LocalAddr(), conn.RemoteAddr())
        go slowProxy(conn)
    }
}
