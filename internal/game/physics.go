package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	gravity       = 0.5
	maxAngularVel = 0.2
	damping       = 0.98
)

type Segment struct {
	x, y       float64
	width      float64
	height     float64
	velocityX  float64
	velocityY  float64
	angle      float64
	angularVel float64
}

func NewSegment(x, y, width, height float64) *Segment {
	return &Segment{x: x, y: y, width: width, height: height}
}

func (s *Segment) Update() {
	s.velocityY += gravity
	s.x += s.velocityX
	s.y += s.velocityY
	s.angle += s.angularVel

	// Apply damping
	s.velocityX *= damping
	s.velocityY *= damping
	s.angularVel *= damping

	// Clamp angular velocity
	if s.angularVel > maxAngularVel {
		s.angularVel = maxAngularVel
	} else if s.angularVel < -maxAngularVel {
		s.angularVel = -maxAngularVel
	}

	// Collision with ground
	if s.y+s.height/2 > 600 {
		s.y = 600 - s.height/2
		s.velocityY = 0
		s.angularVel = 0
	}
}

func (s *Segment) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-s.width/2, -s.height/2)
	op.GeoM.Rotate(s.angle)
	op.GeoM.Translate(s.x, s.y)
	img := ebiten.NewImage(int(s.width), int(s.height))
	img.Fill(color.White)
	screen.DrawImage(img, op)
}

func (s *Segment) Intersects(other *Segment) bool {
	return math.Abs(s.x-other.x) < (s.width+other.width)/2 && math.Abs(s.y-other.y) < (s.height+other.height)/2
}

func (s *Segment) ResolveCollision(other *Segment) {
	if s.Intersects(other) {
		overlapX := (s.width+other.width)/2 - math.Abs(s.x-other.x)
		overlapY := (s.height+other.height)/2 - math.Abs(s.y-other.y)

		if overlapX < overlapY {
			if s.x < other.x {
				s.x -= overlapX / 2
				other.x += overlapX / 2
			} else {
				s.x += overlapX / 2
				other.x -= overlapX / 2
			}
		} else {
			if s.y < other.y {
				s.y -= overlapY / 2
				other.y += overlapY / 2
			} else {
				s.y += overlapY / 2
				other.y -= overlapY / 2
			}
		}
	}
}

type Joint struct {
	segmentA *Segment
	segmentB *Segment
	length   float64
}

func NewJoint(segmentA, segmentB *Segment, length float64) *Joint {
	return &Joint{
		segmentA: segmentA,
		segmentB: segmentB,
		length:   length,
	}
}

func (j *Joint) Update() {
	dx := j.segmentB.x - j.segmentA.x
	dy := j.segmentB.y - j.segmentA.y
	distance := math.Sqrt(dx*dx + dy*dy)
	difference := j.length - distance
	percent := difference / distance / 2 // Adjust this value to strengthen the joint
	offsetX := dx * percent
	offsetY := dy * percent

	j.segmentA.x -= offsetX
	j.segmentA.y -= offsetY
	j.segmentB.x += offsetX
	j.segmentB.y += offsetY

	// Apply forces to create rotation
	offsetRotationX := offsetX * 0.5
	offsetRotationY := offsetY * 0.5
	j.segmentA.angularVel -= offsetRotationX * 0.1
	j.segmentB.angularVel += offsetRotationX * 0.1
	j.segmentA.angularVel -= offsetRotationY * 0.1
	j.segmentB.angularVel += offsetRotationY * 0.1
}

func (j *Joint) Draw(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, j.segmentA.x, j.segmentA.y, j.segmentB.x, j.segmentB.y, color.White)
}
