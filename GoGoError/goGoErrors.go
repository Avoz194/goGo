package GoGoError

import (
    "fmt"
    "log"
    "reflect"
)

type ErrorNum int

const (
    NoSuchEntityError ErrorNum = iota
    EntityAlreadyExists
    InvalidInput
    FailedCommitingRequest
    TechnicalFailrue
)


type GoGoError struct{
    ErrorNum ErrorNum
    EntityType interface{}
    ErrorOnKey string
    ErrorOnValue string
    AdditionalMsg string
    Err error
}

func (goErr *GoGoError)GetError() error{
    log.Println(goErr.Err)
    entType := reflect.TypeOf(goErr.EntityType).Name()
    switch  goErr.ErrorNum {
    case NoSuchEntityError: return fmt.Errorf("A %s with the %s '%s' does not exist. Error: %w", entType, goErr.ErrorOnKey, goErr.ErrorOnValue, goErr.Err)
    case EntityAlreadyExists: return fmt.Errorf("A %s with the %s %s' already exist. Error: %w", entType, goErr.ErrorOnKey, goErr.ErrorOnValue, goErr.Err)
    case InvalidInput: return fmt.Errorf("Value '%s' is not a legal %s. Error: %w", goErr.ErrorOnKey, goErr.ErrorOnValue, goErr.Err)
    case FailedCommitingRequest: return fmt.Errorf("Failed commiting request: %s . Error: %w", goErr.AdditionalMsg, goErr.Err)
    case TechnicalFailrue: return fmt.Errorf("Technical Failure occured. %s Error: %w", goErr.AdditionalMsg, goErr.Err)
    }
    return goErr.Err
}

func (goErr *GoGoError) Error() string {
    return goErr.GetError().Error()
}
