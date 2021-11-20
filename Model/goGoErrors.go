package model

import (
    "fmt"
    "reflect"
)

func NoSuchEntityError(ent interface{}, extraDetails, err error) error{
    entType := reflect.TypeOf(ent).Name()
    return fmt.Errorf("A %s with the %s doesn't exist. Error: %w", entType, extraDetails, err)
}

func EntityAlreadyExists(ent interface{}, extraDetails string, err error) error{
    entType := reflect.TypeOf(ent).Name()
    return fmt.Errorf("A %s with the %s already exist. Error: %w", entType, extraDetails, err)
}

func InvalidInput(ent interface{}, extraDetails string, err error) error{
    entType := reflect.TypeOf(ent).Name()
    return fmt.Errorf("Value %s is not a legal %s. Error: %w", extraDetails, entType, err)
}

func TechnicalFailrue(extraDetails string, err error) error{
    return fmt.Errorf("Technical Failure occured.%s Error: %w", extraDetails, err)
}
