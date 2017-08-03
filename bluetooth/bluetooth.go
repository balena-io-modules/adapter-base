package bluetooth

import (
	"context"
	"fmt"
	"time"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
	"github.com/currantlabs/ble/linux/hci"
	"github.com/currantlabs/ble/linux/hci/cmd"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

var (
	doneChannel chan struct{}
)

func OpenDevice() error {
	device, err := linux.NewDevice()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to open bluetooth device")
		return err
	}

	updateLinuxParam(device)
	ble.SetDefaultDevice(device)

	return nil
}

func CloseDevice() error {
	return ble.Stop()
}

func Connect(id string, timeout time.Duration) (ble.Client, error) {
	client, err := ble.Dial(
		ble.WithSigHandler(
			context.WithTimeout(
				context.Background(),
				timeout*time.Second),
		),
		hci.RandomAddress{
			Addr: ble.NewAddr(id),
		},
	)
	if err != nil {
		return nil, err
	}

	if _, err := client.ExchangeMTU(ble.MaxMTU); err != nil {
		return nil, err
	}

	doneChannel = make(chan struct{})
	go func() {
		<-client.Disconnected()
		close(doneChannel)
	}()

	return client, nil
}

func Disconnect(client ble.Client) error {
	if err := client.ClearSubscriptions(); err != nil {
		return err
	}

	if err := client.CancelConnection(); err != nil {
		return err
	}
	<-doneChannel

	return nil
}

func GetName(id string, timeout time.Duration) (string, error) {
	client, err := Connect(id, timeout)
	if err != nil {
		return "", err
	}

	name, err := GetCharacteristic("2a00", ble.CharRead+ble.CharWrite, 0x02, 0x03)
	if err != nil {
		return "", err
	}

	resp, err := client.ReadCharacteristic(name)
	if err != nil {
		return "", err
	}

	if err := Disconnect(client); err != nil {
		return "", err
	}

	return string(resp), nil
}

func GetCharacteristic(uuid string, property ble.Property, handle, vhandle uint16) (*ble.Characteristic, error) {
	parsedUUID, err := ble.Parse(uuid)
	if err != nil {
		return nil, err
	}

	characteristic := ble.NewCharacteristic(parsedUUID)
	characteristic.Property = property
	characteristic.Handle = handle
	characteristic.ValueHandle = vhandle

	return characteristic, nil
}
func WriteCharacteristic(client ble.Client, characteristic *ble.Characteristic, value []byte, noRsp bool, timeout time.Duration) error {
	err := make(chan error)
	go func() {
		err <- client.WriteCharacteristic(characteristic, value, noRsp)
	}()

	select {
	case done := <-err:
		return done
	case <-time.After(timeout):
		return fmt.Errorf("Write characteristic timed out")
	}
}

func ReadCharacteristic(client ble.Client, characteristic *ble.Characteristic, timeout time.Duration) ([]byte, error) {
	type Result struct {
		Val []byte
		Err error
	}

	result := make(chan Result)
	go func() {
		result <- func() Result {
			val, err := client.ReadCharacteristic(characteristic)
			return Result{val, err}
		}()
	}()

	select {
	case done := <-result:
		return done.Val, done.Err
	case <-time.After(timeout):
		return nil, fmt.Errorf("Read characteristic timed out")
	}
}

func GetDescriptor(uuid string, handle uint16) (*ble.Descriptor, error) {
	parsedUUID, err := ble.Parse(uuid)
	if err != nil {
		return nil, err
	}

	descriptor := ble.NewDescriptor(parsedUUID)
	descriptor.Handle = handle

	return descriptor, nil
}

func WriteDescriptor(client ble.Client, descriptor *ble.Descriptor, value []byte, timeout time.Duration) error {
	err := make(chan error)
	go func() {
		err <- client.WriteDescriptor(descriptor, value)
	}()

	select {
	case done := <-err:
		return done
	case <-time.After(timeout):
		return fmt.Errorf("Write descriptor timed out")
	}
}

func updateLinuxParam(device *linux.Device) error {
	if err := device.HCI.Send(&cmd.LESetScanParameters{
		LEScanType:           0x00,   // 0x00: passive, 0x01: active
		LEScanInterval:       0x0060, // 0x0004 - 0x4000; N * 0.625msec
		LEScanWindow:         0x0060, // 0x0004 - 0x4000; N * 0.625msec
		OwnAddressType:       0x01,   // 0x00: public, 0x01: random
		ScanningFilterPolicy: 0x00,   // 0x00: accept all, 0x01: ignore non-white-listed.
	}, nil); err != nil {
		return errors.Wrap(err, "can't set scan param")
	}

	if err := device.HCI.Option(hci.OptConnParams(
		cmd.LECreateConnection{
			LEScanInterval:        0x0060,    // 0x0004 - 0x4000; N * 0.625 msec
			LEScanWindow:          0x0060,    // 0x0004 - 0x4000; N * 0.625 msec
			InitiatorFilterPolicy: 0x00,      // White list is not used
			PeerAddressType:       0x00,      // Public Device Address
			PeerAddress:           [6]byte{}, //
			OwnAddressType:        0x00,      // Public Device Address
			ConnIntervalMin:       0x0028,    // 0x0006 - 0x0C80; N * 1.25 msec
			ConnIntervalMax:       0x0038,    // 0x0006 - 0x0C80; N * 1.25 msec
			ConnLatency:           0x0000,    // 0x0000 - 0x01F3; N * 1.25 msec
			SupervisionTimeout:    0x002A,    // 0x000A - 0x0C80; N * 10 msec
			MinimumCELength:       0x0000,    // 0x0000 - 0xFFFF; N * 0.625 msec
			MaximumCELength:       0x0000,    // 0x0000 - 0xFFFF; N * 0.625 msec
		})); err != nil {
		return errors.Wrap(err, "can't set connection param")
	}
	return nil
}
