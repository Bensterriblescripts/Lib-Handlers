package network

import (
	"io"
	"net"
	"os"
	"path/filepath"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"

	"golang.org/x/crypto/ssh"
)

func SSHTunnel(client *ssh.Client, localAddr, remoteAddr string) net.Listener {
	TraceLog("Starting SSH tunnel " + localAddr + " -> " + remoteAddr)
	ln := PanicError(net.Listen("tcp", localAddr))
	go func() {
		for {
			lc, err := ln.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne == nil {
					continue
				}
				return // listener closed
			}
			go func() {
				rc := PanicError(client.Dial("tcp", remoteAddr))
				go func() { _, _ = io.Copy(rc, lc) }()
				_, _ = io.Copy(lc, rc)
				PrintErr(lc.Close())
				PrintErr(rc.Close())
			}()
		}
	}()
	return ln
}
func LoadDefaultPrivateKeys() ssh.Signer {
	var path string
	home := PanicError(os.UserHomeDir())
	candidates := []string{
		filepath.Join(home, ".ssh", "id_ed25519"),
		filepath.Join(home, ".ssh", "id_rsa"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			path = p
			break
		}
	}
	if path == "" {
		Panic("No SSH private key found.")
	}

	key := PanicError(os.ReadFile(path))
	return PanicError(ssh.ParsePrivateKey(key))
}
