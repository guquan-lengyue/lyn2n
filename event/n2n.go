package event

var EmptyEvenVar = EmptySignal{}

type EmptySignal struct {
}

var N2NConnectedEvent = Event[any]{}
var N2NDisConnectedEvent = Event[any]{}
var N2NConnectedErr = Event[any]{}
