package service

import (
    "fmt"

    helper "github.com/mataharimall/micro-api/commons"
    loket "github.com/mataharimall/micro-api/components/Loket"
    entity "github.com/mataharimall/micro-api/entities"
)

type EventService struct {
    Response struct {
        Result interface{}
    }
}

func (self *EventService) List() (err error) {

    api := loket.New()
    api.GetAuth()
    token := fmt.Sprintf(`{"token": "%s"}`, api.Token)
    api.Post("v3", "event", token)

    ev := new(entity.Events)
    api.SetStruct(ev)

    x, err := helper.JsMap(ev)
    fmt.Println(x)

    self.Response.Result = x

    return nil
}
