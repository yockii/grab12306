package major

import (
	"bytes"

	"github.com/yockii/grab12306/domain"
)

func generatePassengerTicket(seatType string, passengers []*domain.Passenger) string {
	var bf bytes.Buffer
	for _, p := range passengers {
		bf.WriteString(seatType)
		bf.WriteString(",0,")
		bf.WriteString(p.PassengerType)
		bf.WriteString(",")
		bf.WriteString(p.PassengerName)
		bf.WriteString(",")
		bf.WriteString(p.PassengerIDTypeCode)
		bf.WriteString(",")
		bf.WriteString(p.PassengerIDNo)
		bf.WriteString(",")
		if p.PhoneNo != "" {
			bf.WriteString(p.PhoneNo)
		}
		bf.WriteString(",N_")
	}
	s := bf.String()
	return s[:len(s)-1]
}

func generateOldPassenger(passengers []*domain.Passenger) string {
	var bf bytes.Buffer
	for _, p := range passengers {
		bf.WriteString(p.PassengerName)
		bf.WriteString(",")
		bf.WriteString(p.PassengerIDTypeCode)
		bf.WriteString(",")
		bf.WriteString(p.PassengerIDNo)
		bf.WriteString(",")
		bf.WriteString(p.PassengerType)
		bf.WriteString("_")
	}
	return bf.String()
}
