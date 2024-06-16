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
		{0, 1, 60, 10, 8}, // Head to upper torso (neck)
		{1, 2, 60, 10, 8}, // Upper torso to lower torso
		{1, 3, 40, 8, 6},  // Upper torso to left upper arm
		{3, 4, 40, 8, 6},  // Left upper arm to left lower arm
		{1, 5, 40, 8, 6},  // Upper torso to right upper arm
		{5, 6, 40, 8, 6},  // Right upper arm to right lower arm
		{2, 7, 40, 8, 6},  // Lower torso to left upper leg
		{7, 8, 40, 8, 6},  // Left upper leg to left lower leg
		{2, 9, 40, 8, 6},  // Lower torso to right upper leg
		{9, 10, 40, 8, 6}, // Right upper leg to right lower leg
	}

	jointConstraints := map[int]JointConstraint{
		1: {math.Pi / 4, 3 * math.Pi / 4}, // Neck constraint
		// Add more constraints as necessary
	}

	return &Ragdoll{
		points:           points,
		sticks:           sticks,
		jointConstraints: jointConstraints,
	}
}

// Exposed methods
func (r *Ragdoll) GetPoints() []*Point {
	return r.points
}

func (r *Ragdoll) GetSticks() []*Stick {
	return r.sticks
}

func (r *Ragdoll) GetJointConstraints() map[int]JointConstraint {
	return r.jointConstraints
}

func (r *Ragdoll) SetPointPosition(index int, x, y float64) {
	if index >= 0 && index < len(r.points) {
		r.points[index].x = x
		r.points[index].y = y
	}
}

func (r *Ragdoll) SetStickLength(index int, length float64) {
	if index >= 0 && index < len(r.sticks) {
		r.sticks[index].length = length
	}
}

func (r *Ragdoll) Update() {
	HandleMouseInteractions(r)
	// Update points based on the previous positions
	for _, point := range r.points {
		vx := (point.x - point.oldX) * 0.99 // Apply friction
		vy := (point.y - point.oldY) * 0.99 // Apply friction
		point.oldX = point.x
		point.oldY = point.y
		point.x += vx
		point.y += vy
		point.y += 0.5 // Gravity
	}

	// Constrain the points within the screen bounds
	for _, point := range r.points {
		if point.x < 0 {
			point.x = 0
		}
		if point.y < 0 {
			point.y = 0
		}
		if point.x > 800 {
			point.x = 800
		}
		if point.y > 600 {
			point.y = 600
		}
	}

	// Simulate the physics
	r.SimulatePhysics()
}

func (r *Ragdoll) SimulatePhysics() {
	// Simulate physics for the ragdoll (similar to the existing code)
	for _, stick := range r.sticks {
		p0 := r.points[stick.p0]
		p1 := r.points[stick.p1]
		dx := p1.x - p0.x
		dy := p1.y - p0.y
		distance := distanceBetweenPoints(p0.x, p0.y, p1.x, p1.y)
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
