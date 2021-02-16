package example2

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"math"
	"net"
	"sync"
	"time"
)

var featuresStorage = []*Feature{
	{Location: &Point{Latitude: 0, Longitude: 0}},
	{Location: &Point{Latitude: 1, Longitude: 2}},
	{Location: &Point{Latitude: 5, Longitude: 6}},
	{Location: &Point{Latitude: 7, Longitude: 8}},
	{Location: &Point{Latitude: 3, Longitude: 4}},
	{Location: &Point{Latitude: 0, Longitude: 0}},
	{Location: &Point{Latitude: 3, Longitude: 3}},
	{Location: &Point{Latitude: 3, Longitude: 3}},
	{Location: &Point{Latitude: 5, Longitude: 4}},
	{Location: &Point{Latitude: 4, Longitude: 4}},
	{Location: &Point{Latitude: 3, Longitude: 1}},
	{Location: &Point{Latitude: 2, Longitude: 1}},
	{Location: &Point{Latitude: 1, Longitude: 1}},
	{Location: &Point{Latitude: 2, Longitude: 2}},
	{Location: &Point{Latitude: 0, Longitude: 0}},
}

var routeNotesStorage = map[string][]*RouteNote{}

type routeGuideServer struct {
	mu sync.Mutex
}

func (s *routeGuideServer) GetFeature(ctx context.Context, point *Point) (*Feature, error) {
	_ = ctx
	for _, feature := range featuresStorage {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	return &Feature{Name: "", Location: point}, nil
}

func (s *routeGuideServer) ListFeatures(rect *Rectangle, stream RouteGuide_ListFeaturesServer) error {
	for _, feature := range featuresStorage {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *routeGuideServer) RecordRoute(stream RouteGuide_RecordRouteServer) error {
	var pointsCount, featureCount, distance int32
	var lastPoint *Point
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&RouteSummary{
				PointCount:   pointsCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointsCount++
		for _, feature := range featuresStorage {
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}

func (s *routeGuideServer) RouteChat(stream RouteGuide_RouteChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)

		s.mu.Lock()
		routeNotesStorage[key] = append(routeNotesStorage[key], in)
		rn := make([]*RouteNote, len(routeNotesStorage[key]))
		copy(rn, routeNotesStorage[key])
		s.mu.Unlock()

		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func serialize(point *Point) string {
	return fmt.Sprintf("%d %d", point.Latitude, point.Longitude)
}

func inRange(point *Point, rect *Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	return float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top

}

func calcDistance(p1 *Point, p2 *Point) int32 {
	const (
		CordFactor float64 = 1e7
		R                  = float64(6371000)
	)
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	lat := lat2 - lat1
	lng := lng2 - lng1
	a := math.Sin(lat/2)*math.Sin(lat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(lng/2)*math.Sin(lng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

func CreateAndRunServer(listener net.Listener) *grpc.Server {

	s := grpc.NewServer()
	RegisterRouteGuideServer(s, &routeGuideServer{})

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return s
}
