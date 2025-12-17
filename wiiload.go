package GoWiiload

/*
 * Copyright (C) 2025 Cody Shimizu (OniKuma666)
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// Wiiload protocol structure
type Wiiload struct {
	Magic    [4]byte // HAXX bytes
	Version  uint32  // Wiiload version
	Size     uint32
	Filename [256]byte
	Args     [256]byte
}

const HBC_VERSION_MAJOR = 0
const HBC_VERSION_MINOR = 5

func wiiload_grab_ip() (string, error) { // Looks for the WII environment variable and returns the ip as a string

	ip, exist := os.LookupEnv("WII")

	if !exist || ip == "" {
		return "", fmt.Errorf("WII environment variable not found!")
	}

	return ip, nil
}
func wiiload_grab_file(path string) ([]byte, error) { // Has a map of all the valid extensions.  Zip file sending will come soon
	valid := map[string]bool{
		".dol":  true,
		".elf":  true,
		".wuhb": true,
		".rpx":  true,
	}

	ext := strings.ToLower(filepath.Ext(path))
	if !valid[ext] {
		return nil, fmt.Errorf("Unknown Extension!  Wiiload only takes dol, elf, rpx, or wuhb files!")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading: %w", err)
	}
	return file, nil

}
func wiiload_connect(ip, path string) error {

	data, err := wiiload_grab_file(path)
	if err != nil {
		return fmt.Errorf("Unable to grab file: %w", err)
	}

	if ip == "" {
		ip, err = wiiload_grab_ip()
		if err != nil {
			return fmt.Errorf("unable to grab IP: %w", err)
			}
	}

	parsed := net.ParseIP(ip)
	if parsed == nil || parsed.To4() == nil {
		return fmt.Errorf("invalid IP (only IPv4 is supported): %s", ip)
	}
	ip = parsed.To4().String()

	const port = 4299

	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer conn.Close()
	// Construct the header
	var header Wiiload
	copy(header.Magic[:], []byte("HAXX")) // HAXX is the magic name
	header.Version = uint32((HBC_VERSION_MAJOR << 16) | HBC_VERSION_MINOR)
	header.Size = uint32(len(data))
	copy(header.Filename[:], []byte(filepath.Base(path)))

	// Serialize the header
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.BigEndian, header); err != nil {
		return fmt.Errorf("Failed to serialize: %w", err)
	}
	wiiload_send(conn, data, buffer.Bytes())

	return nil
}
func wiiload_send(conn net.Conn, data []byte, buffer []byte) error {
	// Sending the header and file to the wii
	_, err1 := conn.Write(buffer)
	if err1 != nil {
		return fmt.Errorf("Header failed to send: %w", err1)
	}
	_, err2 := conn.Write(data)
	if err2 != nil {
		return fmt.Errorf("Payload failed to send: %w", err2)
	}
	return nil
}
