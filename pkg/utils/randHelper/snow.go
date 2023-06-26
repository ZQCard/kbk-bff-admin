package randHelper

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sony/sonyflake"
)

func GenerateSnowflakeID() uint64 {
	rand.Seed(time.Now().UnixNano())
	machineID := uint16(time.Now().Unix() % 1024)
	machineID = (machineID << 6) | (uint16(rand.Intn(1<<6)) & 0x3f)
	// Create a new Sonyflake instance with the unique machine ID
	sf := sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (uint16, error) {
			return machineID, nil
		},
	})

	// Generate a new ID
	id, _ := sf.NextID()
	fmt.Println(id)
	return id
}
