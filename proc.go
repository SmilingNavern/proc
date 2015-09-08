package main

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "strings"
    //"os"
)

const ProcPath = "/proc"

type Process struct {
    Name      string
    Pid       int
    Ppid      int
    State     RunState
    Tty       int
    Priority  int
    Nice      int
    Processor int
}

type RunState byte

const (
    RunStateSleep   = 'S'
    RunStateRun     = 'R'
    RunStateStop    = 'T'
    RunStateZombie  = 'Z'
    RunStateIdle    = 'D'
    RunStateUnknown = '?'
)

func (p Process) String() string {
    return fmt.Sprintf("%s %d %d %s %d %d %d %d",
        p.Name, p.Pid, p.Ppid, string(p.State),
        p.Tty, p.Priority, p.Nice, p.Processor)
}

func GetProcess(pid string) (Process, error) {
    procstat := fmt.Sprintf("%s/%s/stat", ProcPath, pid)
    data, err := ioutil.ReadFile(procstat)
    var p Process

    if err != nil {
        return p, err
    }

    fields := strings.Fields(string(data))

    p.Name = fields[1][1 : len(fields[1])-1] // strip ()'s

    p.Pid, _ = strconv.Atoi(pid)

    p.Ppid, _ = strconv.Atoi(fields[3])

    p.State = RunState(fields[2][0])

    p.Tty, _ = strconv.Atoi(fields[6])

    p.Priority, _ = strconv.Atoi(fields[17])

    p.Nice, _ = strconv.Atoi(fields[18])

    p.Processor, _ = strconv.Atoi(fields[38])

    return p, nil
}

func GetProcessList() ([]Process, error) {
    var processlist []Process

    procdir, _ := ioutil.ReadDir(ProcPath)

    for _, p := range procdir {
        if !p.IsDir() {
            continue
        }

        if _, err := strconv.Atoi(p.Name()); err != nil {
            //if not numeric, than it's not pid
            continue
        }

        proc, err := GetProcess(p.Name())
        if err != nil {
            continue
        }

        processlist = append(processlist, proc)

    }

    return processlist, nil
}

func main() {
    processlist, _ := GetProcessList()

    for _, p := range processlist {
        fmt.Println(p)
    }
}
