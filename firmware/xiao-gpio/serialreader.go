package main

import "machine"

type serialReader struct{}

func (serialReader) Read(p []byte) (int, error) {
	for machine.Serial.Buffered() == 0 {
	}
	return machine.Serial.Read(p)
}

func (serialReader) ReadByte() (byte, error) {
	for machine.Serial.Buffered() == 0 {
	}

	return machine.Serial.ReadByte()
}
