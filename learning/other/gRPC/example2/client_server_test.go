package example2

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

var (
	client    RouteGuideClient
	ctx       context.Context
	ctxCancel context.CancelFunc
)

func TestMain(m *testing.M) {

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := CreateAndRunServer(listener)
	defer server.GracefulStop()

	client = CreateClient(listener.Addr().String())
	ctx, ctxCancel = context.WithTimeout(context.Background(), time.Second*60)
	defer ctxCancel()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestSimpleRPC(t *testing.T) {

	tests := []struct {
		name string
		args struct{ point *Point }
		want *Feature
	}{
		{
			name: "Point [3;4]",
			args: struct{ point *Point }{&Point{Latitude: 3, Longitude: 4}},
			want: &Feature{Location: &Point{Latitude: 3, Longitude: 4}},
		},
		{
			name: "Point [1;2]",
			args: struct{ point *Point }{&Point{Latitude: 1, Longitude: 2}},
			want: &Feature{Location: &Point{Latitude: 1, Longitude: 2}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			feature, err := client.GetFeature(ctx, test.args.point)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, test.want, feature)
		})
	}
}

func TestServerSideStreamingRpc(t *testing.T) {
	tests := []struct {
		name            string
		rect            *Rectangle
		wantPointsCount int
	}{
		{
			name:            "Twelve points rectangle",
			rect:            &Rectangle{Lo: &Point{Latitude: 4, Longitude: 0}, Hi: &Point{Latitude: 0, Longitude: 4}},
			wantPointsCount: 12,
		},
		{
			name:            "Three points rectangle",
			rect:            &Rectangle{Lo: &Point{Latitude: 9, Longitude: 0}, Hi: &Point{Latitude: 4, Longitude: 10}},
			wantPointsCount: 4,
		},
		{
			name:            "Eight points rectangle",
			rect:            &Rectangle{Lo: &Point{Latitude: 7, Longitude: 2}, Hi: &Point{Latitude: 1, Longitude: 7}},
			wantPointsCount: 8,
		},
		{
			name:            "One point rectangle",
			rect:            &Rectangle{Lo: &Point{Latitude: 8, Longitude: 7}, Hi: &Point{Latitude: 6, Longitude: 10}},
			wantPointsCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream, err := client.ListFeatures(ctx, test.rect)
			if err != nil {
				t.Error(err)
			}

			var receivedFeatures []*Feature
			for {
				feature, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Error(err)
				}
				receivedFeatures = append(receivedFeatures, feature)
			}

			assert.Equal(t, test.wantPointsCount, len(receivedFeatures))
		})
	}
}

func TestClientSideStreamingRPC(t *testing.T) {
	tests := []struct {
		name         string
		points       []*Point
		wantDistance int32
	}{
		{
			name: "Three points",
			points: []*Point{
				{Latitude: 5, Longitude: 61},
				{Latitude: 47, Longitude: 101},
				{Latitude: 31, Longitude: 50},
				{Latitude: 311, Longitude: 1},
				{Latitude: 301, Longitude: 12},
				{Latitude: 31, Longitude: 1},
				{Latitude: 0, Longitude: 92},
			},
			wantDistance: 7,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream, err := client.RecordRoute(context.Background())
			if err != nil {
				t.Error(err)
			}
			for _, point := range test.points {
				err = stream.Send(point)
				if err != nil {
					t.Error(err)
				}
			}
			routeSummary, err := stream.CloseAndRecv()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, test.wantDistance, routeSummary.Distance)
		})
	}

}

func TestBidirectionalStreamingRPC(t *testing.T) {

	tests := []struct {
		name       string
		points     []*RouteNote
		notesCount int
	}{
		{
			name: "Nine notes",
			points: []*RouteNote{
				{Location: &Point{Latitude: 0, Longitude: 1}, Message: "First message"},
				{Location: &Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
				{Location: &Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
				{Location: &Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
				{Location: &Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
				{Location: &Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
			},
			notesCount: 9,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			stream, err := client.RouteChat(context.Background())
			if err != nil {
				t.Error(err)
			}
			var receivedNotes []*RouteNote
			wait := make(chan bool)
			go func() {
				for {
					in, err := stream.Recv()
					if err == io.EOF {
						close(wait)
						return
					}
					if err != nil {
						log.Fatalf("Failed to receive a note : %v", err)
					}
					receivedNotes = append(receivedNotes, in)
				}
			}()
			for _, note := range test.points {
				if err := stream.Send(note); err != nil {
					log.Fatalf("Failed to send a note: %v", err)
				}
			}
			_ = stream.CloseSend()
			<-wait
			assert.Equal(t, test.notesCount, len(receivedNotes))
		})
	}

}
