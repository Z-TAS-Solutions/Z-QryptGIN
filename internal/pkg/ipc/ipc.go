package ipc

import (
	"encoding/binary"
	"io"
	"log"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/grpc/zproto"
	"github.com/natefinch/npipe"
	"google.golang.org/protobuf/proto"
)

func StartServer(pipeName string, handle func(conn io.ReadWriteCloser)) {
	listener, err := npipe.Listen(pipeName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Go ZIPC server running on ....", pipeName)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to Connect : ", err)
			continue
		}
		go handle(conn)
	}
}

func HandleConnection(conn io.ReadWriteCloser) {
	defer conn.Close()

	lengthBuffer := make([]byte, 4)
	_, err := io.ReadFull(conn, lengthBuffer)
	if err != nil {
		log.Println("Failed To Read Msg Length : ", err)
		return
	}
	length := int(binary.LittleEndian.Uint32(lengthBuffer))

	payload := make([]byte, length)
	_, err = io.ReadFull(conn, payload)
	if err != nil {
		log.Println("Failed To Read Msg : ", err)
		return
	}

	zpacket := &zproto.ZIPCPacket{}
	err = proto.Unmarshal(payload, zpacket)
	if err != nil {
		log.Println("Failed To Decode Payload : ", err)
		return
	}

	log.Println("Received ZIPCPacket:", zpacket)

}
