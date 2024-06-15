package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Detects and handles mouse interactions with the ragdoll
func HandleMouseInteractions(r *Ragdoll) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		mouseX, mouseY := float64(x), float64(y)

		if r.draggedPoint == nil {
			// Find the closest point to drag
			closestPoint := findClosestPoint(r, mouseX, mouseY)
			if distanceBetweenPoints(closestPoint.x, closestPoint.y, mouseX, mouseY) < 20 {
				r.draggedPoint = closestPoint
			}
		}

		if r.draggedPoint != nil {
			// Drag the point
			r.draggedPoint.oldX = r.draggedPoint.x
			r.draggedPoint.oldY = r.draggedPoint.y
			r.draggedPoint.x = mouseX
			r.draggedPoint.y = mouseY
		}
	} else {
		r.draggedPoint = nil
	}
}

// Finds the closest point to the given coordinates
func findClosestPoint(r *Ragdoll, x, y float64) *Point {
	var closestPoint *Point
	minDist := math.MaxFloat64

	for _, point := range r.points {
		dist := distanceBetweenPoints(point.x, point.y, x, y)
		if dist < minDist {
			minDist = dist
			closestPoint = point
		}
	}

	return closestPoint
}
