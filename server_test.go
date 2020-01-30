package gmqtt

import (
	"context"
	"net"
	"reflect"
	"testing"

	"github.com/DrmagicE/gmqtt/pkg/packets"
)

func TestHooks(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:1883")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var hooksStr string
	hooks := Hooks{
		OnAccept: func(ctx context.Context, conn net.Conn) bool {
			hooksStr += "Accept"
			return true
		},
		OnConnect: func(ctx context.Context, client Client) (code uint8) {
			hooksStr += "OnConnect"
			return packets.CodeAccepted
		},
		OnConnected: func(ctx context.Context, client Client) {
			hooksStr += "OnConnected"
		},
		OnSessionCreated: func(ctx context.Context, client Client) {
			hooksStr += "OnSessionCreated"
		},
		OnSubscribe: func(ctx context.Context, client Client, topic packets.Topic) (qos uint8) {
			hooksStr += "OnSubscribe"
			return packets.QOS_1
		},
		OnMsgArrived: func(ctx context.Context, client Client, msg packets.Message) (valid bool) {
			hooksStr += "OnMsgArrived"
			return true
		},
		OnSessionTerminated: func(ctx context.Context, client Client, reason SessionTerminatedReason) {
			hooksStr += "OnSessionTerminated"
		},
		OnClose: func(ctx context.Context, client Client, err error) {
			hooksStr += "OnClose"
		},
		OnStop: func(ctx context.Context) {
			hooksStr += "OnStop"
		},
	}
	srv := NewServer(
		WithTCPListener(ln),
		WithHook(hooks))
	srv.Run()

	c, err := net.Dial("tcp", "127.0.0.1:1883")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	w := packets.NewWriter(c)
	r := packets.NewReader(c)
	w.WriteAndFlush(defaultConnectPacket())
	r.ReadPacket(0x04)

	sub := &packets.Subscribe{
		PacketID: 10,
		Topics: []packets.Topic{
			{Name: "name", Qos: packets.QOS_1},
		},
	}
	w.WriteAndFlush(sub)
	r.ReadPacket(0x04) //suback

	pub := &packets.Publish{
		Dup:       false,
		Qos:       packets.QOS_1,
		Retain:    false,
		TopicName: []byte("ok"),
		PacketID:  10,
		Payload:   []byte("payload"),
	}
	w.WriteAndFlush(pub)
	r.ReadPacket(0x04) //puback
	srv.Stop(context.Background())
	want := "AcceptOnConnectOnConnectedOnSessionCreatedOnSubscribeOnMsgArrivedOnSessionTerminatedOnCloseOnStop"
	if hooksStr != want {
		t.Fatalf("hooksStr error, want %s, got %s", want, hooksStr)
	}
}

func TestZeroBytesClientId(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:1883")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	srv := NewServer(WithTCPListener(ln))
	defer srv.Stop(context.Background())
	srv.Run()
	c, err := net.Dial("tcp", "127.0.0.1:1883")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	w := packets.NewWriter(c)
	r := packets.NewReader(c)
	connect := defaultConnectPacket()
	connect.CleanSession = true
	connect.ClientID = make([]byte, 0)
	w.WriteAndFlush(connect)
	p, err := r.ReadPacket(0x04)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if ack, ok := p.(*packets.Connack); ok {
		if ack.Code != packets.CodeAccepted {
			t.Fatalf("connack.Code error, want %d, but got %d", packets.CodeAccepted, ack.Code)
		}
	} else {
		t.Fatalf("invalid type, want %v, got %v", reflect.TypeOf(&packets.Connack{}), reflect.TypeOf(p))
	}
	c2, err := net.Dial("tcp", "127.0.0.1:1883")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	w2 := packets.NewWriter(c2)
	r2 := packets.NewReader(c2)
	connect2 := defaultConnectPacket()
	connect2.CleanSession = true
	connect2.ClientID = make([]byte, 0)
	w2.WriteAndFlush(connect2)
	p, err = r2.ReadPacket(0x04)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if ack, ok := p.(*packets.Connack); ok {
		if ack.Code != packets.CodeAccepted {
			t.Fatalf("connack.Code error, want %d, but got %d", packets.CodeAccepted, ack.Code)
		}
	} else {
		t.Fatalf("invalid type, want %v, got %v", reflect.TypeOf(&packets.Connack{}), reflect.TypeOf(p))
	}
}

func TestRandUUID(t *testing.T) {
	uuids := make(map[string]struct{})
	for i := 0; i < 100; i++ {
		uuids[getRandomUUID()] = struct{}{}
	}
	if len(uuids) != 100 {
		t.Fatalf("duplicated ID")
	}
}
