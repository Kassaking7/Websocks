package Websocks

import (
	"encoding/binary"
	"net"
)

type LsServer struct {
	Cipher     *Cipher
	ListenAddr *net.TCPAddr
}

func NewLsServer(password string, listenAddr string) (*LsServer, error) {
	bsPassword, err := ParsePassword(password)
	if err != nil {
		return nil, err
	}
	structListenAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &LsServer{
		Cipher:     NewCipher(bsPassword),
		ListenAddr: structListenAddr,
	}, nil

}

func (lsServer *LsServer) Listen(didListen func(listenAddr *net.TCPAddr)) error {
	return ListenEncryptedTCP(lsServer.ListenAddr, lsServer.Cipher, lsServer.handleConn, didListen)
}

// https://www.ietf.org/rfc/rfc1928.txt
func (lsServer *LsServer) handleConn(localConn *SecureTCPConn) {
	defer localConn.Close()
	buf := make([]byte, 256)
	//socks5 handshake
	_, err := localConn.DecodeRead(buf)
	if err != nil || buf[0] != 0x05 {
		return
	}
	// no verification
	localConn.EncodeWrite([]byte{0x05, 0x00})

	// Get the remote address
	n, err := localConn.DecodeRead(buf)
	if err != nil || n < 7 {
		return
	}
	// CMD
	// CONNECT X'01'
	if buf[1] != 0x01 {
		return
	}
	var dIP []byte
	// aType
	switch buf[3] {
	case 0x01:
		//	IP V4 address: X'01'
		dIP = buf[4 : 4+net.IPv4len]
	case 0x03:
		//	DOMAINNAME: X'03'
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))
		if err != nil {
			return
		}
		dIP = ipAddr.IP
	case 0x04:
		//	IP V6 address: X'04'
		dIP = buf[4 : 4+net.IPv6len]
	default:
		return
	}
	dPort := buf[n-2:]
	dstAddr := &net.TCPAddr{
		IP:   dIP,
		Port: int(binary.BigEndian.Uint16(dPort)),
	}
	// Connect to remote address
	dstServer, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		return
	} else {
		defer dstServer.Close()
		dstServer.SetLinger(0)
		localConn.EncodeWrite([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}
	go func() {
		err := localConn.DecodeCopy(dstServer)
		if err != nil {
			// case for network timeout
			localConn.Close()
			dstServer.Close()
		}
	}()
	(&SecureTCPConn{
		Cipher:          localConn.Cipher,
		ReadWriteCloser: dstServer,
	}).EncodeCopy(localConn)
}
