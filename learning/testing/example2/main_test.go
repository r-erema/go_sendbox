package example2

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"testing"
)

type smsServiceMock struct {
	mock.Mock
}

func (m *smsServiceMock) SendChargeNotification(value int) error {
	fmt.Println("Mocked charge notification function")
	fmt.Printf("Value passed in: %d\n", value)
	args := m.Called(value)
	return args.Error(0)
}

func TestMyService_ChargeCustomer(t *testing.T) {
	smsService := new(smsServiceMock)
	smsService.On("SendChargeNotification", 100).Return(nil)
	myService := MyService{smsService}
	_ = myService.ChargeCustomer(100)
	smsService.AssertExpectations(t)
}
