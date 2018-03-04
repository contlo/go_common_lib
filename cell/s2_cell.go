package cell

// https://github.com/topfreegames/apm
// https://elithrar.github.io/article/running-go-applications-in-the-background/
import (
  "github.com/golang/geo/s2"
  "github.com/golang/geo/s1"
)

type Cell struct {
  Token string
  Center []float64
}

func FetchCell(latitude float64, longitude float64, level int) Cell {
	cellId := s2.CellIDFromLatLng(s2.LatLngFromDegrees(latitude, longitude))
	cellAtLevel := cellId.Parent(level)
  rect := s2.CellFromCellID(cellAtLevel).RectBound()
  lat := float64(s1.Angle(rect.Center().Lat) / s1.Degree)
	lng := float64(s1.Angle(rect.Center().Lng) / s1.Degree)
	return Cell{Token: cellAtLevel.ToToken(), Center: []float64{lat, lng}}
}

func FetchCellToken(latitude float64, longitude float64, level int) string {
  return FetchCell(latitude, longitude, level).Token
}
