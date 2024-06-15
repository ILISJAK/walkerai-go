package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"gonum.org/v1/gonum/floats"
)

// Helper functions
func angleBetweenPoints(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y2-y1, x2-x1)
}

func clampAngle(angle, min, max float64) float64 {
	for angle < min {
		angle += 2 * math.Pi
	}
	for angle > max {
		angle -= 2 * math.Pi
	}
	return angle
}

func distanceBetweenPoints(x1, y1, x2, y2 float64) float64 {
	return floats.Distance([]float64{x1, y1}, []float64{x2, y2}, 2.0)
}

type Point struct {
	x, y       float64
	oldX, oldY float64
	mass       float64
}

type Stick struct {
	p0, p1  int
	length  float64
	radiusX float64 // Major axis
	radiusY float64 // Minor axis
}

type JointConstraint struct {
	minAngle float64
	maxAngle float64
}

type Ragdoll struct {
	points           []*Point
	sticks           []*Stick
	jointConstraints map[int]JointConstraint
	draggedPoint     *Point
}

func NewRagdoll() *Ragdoll {
	// Initialize points with mass
	points := []*Point{
		{400, 100, 400, 100, 4.0}, // Head
		{400, 160, 400, 160, 8.0}, // Upper torso
		{400, 220, 400, 220, 8.0}, // Lower torso
		{370, 140, 370, 140, 2.0}, // Left upper arm
		{370, 180, 370, 180, 1.0}, // Left lower arm
		{430, 140, 430, 140, 2.0}, // Right upper arm
		{430, 180, 430, 180, 1.0}, // Right lower arm
		{390, 260, 390, 260, 6.0}, // Left upper leg
		{390, 300, 390, 300, 4.0}, // Left lower leg
		{410, 260, 410, 260, 6.0}, // Right upper leg
		{410, 300, 410, 300, 4.0}, // Right lower leg
	}

	// Initialize sticks with ellipses (major and minor radii)
	sticks := []*Stick{
		{0, 1, 60, 10, 8},  // Head to upper torso (neck)
		{1, 2, 60, 10, 8},  // Upper torso to lower torso
		{1, 3, 40, 8, 6},   // Upper torso to left upper arm
		{3, 4, 40, 8, 6},   // Left upper arm to left lower arm
		{1, 5, 40, 8, 6},   // Upper torso to right upper arm
		{5, 6, 40, 8, 6},   // Right upper arm to right lower arm
		{2, 7, 60, 10, 8},  // Lower torso to left upper leg
		{7, 8, 40, 10, 8},  // Left upper leg to left lower leg
		{2, 9, 60, 10, 8},  // Lower torso to right upper leg
		{9, 10, 40, 10, 8}, // Right upper leg to right lower leg
	}

	jointConstraints := map[int]JointConstraint{
		// Head to upper torso (neck)
		1: {minAngle: -0.5 * math.Pi, maxAngle: 0.5 * math.Pi},

		// Upper torso to lower torso
		2: {minAngle: -0.5 * math.Pi, maxAngle: 0.5 * math.Pi},

		// Upper torso to upper arms
		3: {minAngle: -0.5 * math.Pi, maxAngle: 0.5 * math.Pi}, // Left shoulder
		5: {minAngle: -0.5 * math.Pi, maxAngle: 0.5 * math.Pi}, // Right shoulder

		// Upper arms to lower arms
		4: {minAngle: 0, maxAngle: math.Pi}, // Left elbow
		6: {minAngle: 0, maxAngle: math.Pi}, // Right elbow

		// Lower torso to upper legs
		7: {minAngle: -0.5 * math.Pi, maxAngle: 0.5 * math.Pi}, // Left hip
		9: {minAngle: -0.5 * math.Pi, maxAngle: 0.5 * math.Pi}, // Right hip

		// Upper legs to lower legs
		8:  {minAngle: 0, maxAngle: math.Pi}, // Left knee
		10: {minAngle: 0, maxAngle: math.Pi}, // Right knee
	}

	return &Ragdoll{
		points:           points,
		sticks:           sticks,
		jointConstraints: jointConstraints,
	}
}

func (r *Ragdoll) Update() {
	HandleMouseInteractions(r) // Call the mouse interaction handler

	for _, point := range r.points {
		// Apply gravity
		vy := (point.y-point.oldY)*(1-0.02) + gravity*point.mass
		vx := (point.x - point.oldX) * (1 - 0.02)

		point.oldX = point.x
		point.oldY = point.y
		point.x += vx
		point.y += vy

		// Constrain within screen bounds
		if point.y > 600 {
			point.y = 600
			point.oldY = point.y + vy*-0.5
		}
		if point.x < 0 {
			point.x = 0
			point.oldX = point.x + vx*-0.5
		}
		if point.x > 800 {
			point.x = 800
			point.oldX = point.x + vx*-0.5
		}
	}

	for i := 0; i < 5; i++ {
		for _, stick := range r.sticks {
			p0 := r.points[stick.p0]
			p1 := r.points[stick.p1]

			dx := p1.x - p0.x
			dy := p1.y - p0.y
			distance := math.Sqrt(dx*dx + dy*dy)
			difference := stick.length - distance
			percent := difference / distance / 2
			offsetX := dx * percent
			offsetY := dy * percent

			p0.x -= offsetX
			p0.y -= offsetY
			p1.x += offsetX
			p1.y += offsetY

			// Apply joint constraints
			if stick.p0 != 1 && stick.p1 != 1 { // Adjust for non-root points (torso is root)
				if constraint, ok := r.jointConstraints[stick.p0]; ok {
					currentAngle := angleBetweenPoints(p0.x, p0.y, p1.x, p1.y)
					clampedAngle := clampAngle(currentAngle, constraint.minAngle, constraint.maxAngle)
					if currentAngle != clampedAngle {
						distance := stick.length
						p1.x = p0.x + math.Cos(clampedAngle)*distance
						p1.y = p0.y + math.Sin(clampedAngle)*distance
					}
				}
			}
		}
	}

	// Prevent overlapping volumes
	for i := 0; i < 3; i++ { // Multiple iterations for stable resolution
		for _, stickA := range r.sticks {
			for _, stickB := range r.sticks {
				if stickA == stickB {
					continue
				}

				p0A := r.points[stickA.p0]
				p1A := r.points[stickA.p1]
				p0B := r.points[stickB.p0]
				p1B := r.points[stickB.p1]

				// Check if the volumes intersect
				centerA := Point{(p0A.x + p1A.x) / 2, (p0A.y + p1A.y) / 2, (p0A.x + p1A.x) / 2, (p0A.y + p1A.y) / 2, (p0A.mass + p1A.mass) / 2}
				centerB := Point{(p0B.x + p1B.x) / 2, (p0B.y + p1B.y) / 2, (p0B.x + p1B.x) / 2, (p0B.y + p1B.y) / 2, (p0B.mass + p1B.mass) / 2}
				dist := distanceBetweenPoints(centerA.x, centerA.y, centerB.x, centerB.y)
				minDist := math.Max(stickA.radiusX, stickA.radiusY) + math.Max(stickB.radiusX, stickB.radiusY)

				if dist < minDist {
					// Move points to separate volumes
					overlap := minDist - dist
					moveX := overlap * (centerA.x - centerB.x) / dist / 2
					moveY := overlap * (centerA.y - centerB.y) / dist / 2

					p0A.x += moveX
					p0A.y += moveY
					p1A.x += moveX
					p1A.y += moveY

					p0B.x -= moveX
					p0B.y -= moveY
					p1B.x -= moveX
					p1B.y -= moveY
				}
			}
		}
	}
}

func (r *Ragdoll) Draw(screen *ebiten.Image) {
	for _, stick := range r.sticks {
		p0 := r.points[stick.p0]
		p1 := r.points[stick.p1]
		ebitenutil.DrawLine(screen, p0.x, p0.y, p1.x, p1.y, color.White)

		// Draw ellipses around the sticks to represent body parts
		centerX := (p0.x + p1.x) / 2
		centerY := (p0.y + p1.y) / 2
		drawEllipse(screen, centerX, centerY, stick.radiusX, stick.radiusY, color.White)
	}

	// Draw volume around points to represent body parts
	for _, point := range r.points {
		ebitenutil.DrawRect(screen, point.x-5, point.y-5, 10, 10, color.White)
	}

	// Draw a larger circle for the head
	head := r.points[0]
	ebitenutil.DrawCircle(screen, head.x, head.y, 20, color.RGBA{255, 0, 0, 255}) // Red color for visibility
}

// DrawEllipse draws an ellipse on the screen
func drawEllipse(screen *ebiten.Image, centerX, centerY, radiusX, radiusY float64, clr color.Color) {
	const points = 100
	var theta float64
	for i := 0; i < points; i++ {
		theta = float64(i) * 2 * math.Pi / points
		x := radiusX * math.Cos(theta)
		y := radiusY * math.Sin(theta)
		ebitenutil.DrawRect(screen, centerX+x, centerY+y, 1, 1, clr)
	}
}
