package main
import (
	"reflect"
	"fmt"
	"errors"
)


type Data struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool
	FormCompletionTime int
}


type Dimension struct {
	Width  string
	Height string
}


func NewData(ev Event) *Data {
	d := new(Data)
	d.ResizeTo = ev.ResizeTo
	d.ResizeFrom = ev.ResizeFrom
	d.WebsiteUrl = gh(ev.WebsiteUrl) 	//change to hash
	
	d.CopyAndPaste = make(map[string]bool,3)
	if ev.Pasted {
		d.CopyAndPaste[ev.FormId] = true
	}
	d.FormCompletionTime = ev.Time
	d.SessionId = ev.SessionId
	return d
}


func (d *Data) isCompleted() bool {
	if d.FormCompletionTime > 0 {
		return true
	}
	return false
}



func (d *Data) updateSession(ev Event) error {
	
	switch ev.EventType {
	case "copyAndPaste":
		d.CopyAndPaste[ev.FormId] = true
	case "screenResize":
		d.ResizeFrom = ev.ResizeFrom
		d.ResizeTo = ev.ResizeTo
	case "timeTaken":
		d.FormCompletionTime = ev.Time
	default :
		return errors.New("Event type not supported " + ev.EventType)
	}
	return nil
}



func (d *Data) Print() {
	v := reflect.ValueOf(*d)
	typeOfS := v.Type()
	fmt.Printf("Data struct: \n")
	for i := 0; i< v.NumField(); i++ {
		fmt.Printf("\t%s\t\t\t: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
	fmt.Printf("{\n")
}