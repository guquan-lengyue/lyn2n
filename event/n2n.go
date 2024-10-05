package event

var EmptyEvenVar = EmptySignal{}

type EmptySignal struct {
}

var N2NConnectedEvent = make(chan EmptySignal, 1)
var N2NDisConnectedEvent = make(chan EmptySignal, 1)
var N2NConnectedErr = make(chan EmptySignal, 1)
