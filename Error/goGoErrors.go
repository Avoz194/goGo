package Error

import (
    "fmt"
    "log"
    "reflect"
)

func NoSuchEntityError(ent interface{}, extraDetails string, err error) error{
    entType := reflect.TypeOf(ent).Name()
    log.Println(err)
    return fmt.Errorf("A %s with the %s does not exist. Error: %w", entType, extraDetails, err)
}

func EntityAlreadyExists(ent interface{}, extraDetails string, err error) error{
    entType := reflect.TypeOf(ent).Name()
    log.Println(err)
    return fmt.Errorf("A %s with the %s already exist. Error: %w", entType, extraDetails, err)
}

func InvalidInput(ent interface{}, extraDetails string, err error) error{
    entType := reflect.TypeOf(ent).Name()
    log.Println(err)
    return fmt.Errorf("Value %s is not a legal %s. Error: %w", extraDetails, entType, err)
}

func FailedCommitingRequest(extraDetails string, err error) error{
    log.Println(err)
    return fmt.Errorf("Failed commiting %s request. Error: %w", extraDetails, err)
}

func TechnicalFailrue(extraDetails string, err error) error{
    log.Fatal(err)
    return fmt.Errorf("Technical Failure occured.%s Error: %w", extraDetails, err)
}
