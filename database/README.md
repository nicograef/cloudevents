# database

**database** is an event-sourcing database written in Go. It stores events conforming the cloudevent specification.
The persistence is implemented by appending the events in a json format to a newline-delimited json file on the disk/volume.