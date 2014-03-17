package gowork

import (
    "fmt"
    "reflect"
    "sync"
    "testing"
)

var (
    SimpleWorker Worker
    EchoWork     BusyWork
    Adder        BusyWork
    once         sync.Once
)

func setup() {
    SimpleWorker = Worker{
        Name: "Simple Worker",
        WorkData: map[string]string{
            "type": "echo",
        },
    }
    EchoWork = func(args ...interface{}) (interface{}, error) {
        if len(args) < 1 {
            return nil, fmt.Errorf("Invalid args")
        }

        var retVals []interface{}
        for _, value := range args {
            retVals = append(retVals, value)
        }

        return retVals, nil
    }
    Adder = func(args ...interface{}) (interface{}, error) {
        if len(args) < 1 {
            return nil, fmt.Errorf("Invalid args, at least on value required")
        }

        var sum int
        for _, value := range args {
            sum += value.(int)
        }

        return sum, nil
    }
}

func checkWorkerEqual(t *testing.T, expected *Worker, actual *Worker) {
    if actual == nil {
        t.Error("No actual worker found")
    }
    if actual.Name != expected.Name {
        t.Errorf("Name mismatch. Expected [%v] Found [%v]", expected.Name, actual.Name)
    }
    if reflect.TypeOf(actual.WorkData) != reflect.TypeOf(expected.WorkData) {
        t.Errorf(
            "WorkData mismatch. Expected [%v] Found [%v]",
            reflect.TypeOf(expected.WorkData),
            reflect.TypeOf(actual.WorkData),
        )
    }
}

func checkDeepEquality(t *testing.T, expected []interface{}, actual []interface{}) {
    if len(expected) != len(actual) {
        t.Fatalf("Return value length mismatch. Expected [%v] Found [%v]", len(expected), len(actual))
    }
    for i, _ := range expected {
        if expected[i] != actual[i] {
            t.Errorf("Value mismatch at [%v]. Expected [%v] Found [%v]", i, expected[i], actual[i])
        }
    }
}

func TestWorkerEquality(t *testing.T) {
    once.Do(setup)

    w := Worker{}
    w.Name = "Simple Worker"
    w.WorkData = make(map[string]string)
    w.WorkData.(map[string]string)["type"] = "echo"

    checkWorkerEqual(t, &w, &SimpleWorker)
}

func TestBusyWork(t *testing.T) {
    var bw BusyWork = func(args ...interface{}) (interface{}, error) {
        return nil, nil
    }
    eRet, eErr := bw()
    if eRet != nil {
        t.Errorf("Expected BW return value: [%v] Actual: [%v]", nil, eRet)
    }

    if eErr != nil {
        t.Errorf("Expected BW return error: [%v] Actual: [%v]", nil, eErr)
    }
}

func TestEcho1(t *testing.T) {
    once.Do(setup)

    var intArr []interface{}
    for i := 1; i < 25; i++ {
        intArr = append(intArr, i)
        retVal, err := SimpleWorker.Process(EchoWork, intArr...)
        if err != nil {
            t.Errorf("Unexpected error: [%v]", err)
        }
        checkDeepEquality(t, intArr, retVal.([]interface{}))
    }
}

func TestEcho2(t *testing.T) {
    once.Do(setup)

    var stringArr string
    for i := 1; i < 25; i++ {
        stringArr += "h"
        retVal, err := SimpleWorker.Process(EchoWork, stringArr)
        if err != nil {
            t.Errorf("Unexpected error: [%v]", err)
        }
        checkDeepEquality(t, []interface{}{stringArr}, retVal.([]interface{}))
    }
}

func TestAdder1(t *testing.T) {
    once.Do(setup)

    var intArr []interface{}
    sum := 0
    for i := 1; i < 25; i++ {
        intArr = append(intArr, i)
        sum += i
        retVal, err := SimpleWorker.Process(Adder, intArr...)
        if err != nil {
            t.Errorf("Unexpected error: [%v]", err)
        }
        if retVal != sum {
            t.Errorf("Expected Adder return error: [%v] Actual: [%v]", sum, retVal)
        }
    }
}
