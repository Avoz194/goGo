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

//GoGoError
//  ErrorNum: The type of the error.
//  EntityType: The type of the Entity on which the error occurs (Person, Task etc).
//  ErrorOnKey: The invalid field.
//  ErrorOnValue: The worng value.
//  AdditionalMsg: Optional information (func name).
//  Err: The original error received.
type GoGoError struct{
    ErrorNum ErrorNum
    EntityType interface{}
    ErrorOnKey string
    ErrorOnValue string
    AdditionalMsg string
    Err error
}

// GetError returning the goGoError message, according to the ErrorNum.
func (goErr *GoGoError)GetError() error{
    if goErr.EntityType != nil {
        log.Println(goErr.Err)
        entType := reflect.TypeOf(goErr.EntityType).Name()
        switch goErr.ErrorNum {
        case NoSuchEntityError:
            return fmt.Errorf("A %s with the %s '%s' does not exist.", entType, goErr.ErrorOnKey, goErr.ErrorOnValue)
        case EntityAlreadyExists:
            return fmt.Errorf("A %s with the %s '%s' already exist.", entType, goErr.ErrorOnKey, goErr.ErrorOnValue)
        case InvalidInput:
            return fmt.Errorf("Value '%s' is not a legal %s.", goErr.ErrorOnValue, goErr.ErrorOnKey)
        case FailedCommitingRequest:
            return fmt.Errorf("Failed commiting request: %s.", goErr.AdditionalMsg)
        case TechnicalFailrue:
            return fmt.Errorf("Technical Failure occured: %s.", goErr.AdditionalMsg)
        }
    }
    return goErr.Err
}

func (goErr *GoGoError) Error() string {
    return goErr.GetError().Error()
}
