package example2

import "fmt"

type MessageService interface {
	SendChargeNotification(int) error
}

type SMSService struct{}

func (sms *SMSService) SendChargeNotification(value int) error {
	_ = value
	fmt.Println("Sending Production Charge Notification")
	return nil
}

type MyService struct {
	messageService MessageService
}

func (ms *MyService) ChargeCustomer(value int) error {
	_ = ms.messageService.SendChargeNotification(value)
	fmt.Printf("Charging Customer For the value of %d\n", value)
	return nil
}
