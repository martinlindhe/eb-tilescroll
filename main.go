// Copyright 2016 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var (
	thePlayer            = &player{}
	gophersImage         *ebiten.Image
	repeatedGophersImage *ebiten.Image
	groundImage          *ebiten.Image
)

type player struct {
	x16 int
	y16 int
}

func round(x float64) float64 {
	return math.Floor(x + 0.5)
}

func (p *player) MoveForward() {

	w, h := gophersImage.Size()
	mx := w * 16
	my := h * 16
	s, c := math.Sincos(2 * math.Pi)
	p.x16 += int(round(16*c) * 2)
	p.y16 += int(round(16*s)*2) + screenHeight/8
	for mx <= p.x16 {
		p.x16 -= mx
	}
	for my <= p.y16 {
		p.y16 -= my
	}
	for p.x16 < 0 {
		p.x16 += mx
	}
	for p.y16 < 0 {
		p.y16 += my
	}
}

func (p *player) Position() (int, int) {
	return p.x16, p.y16
}

func updateGroundImage(ground *ebiten.Image) error {

	if err := ground.Clear(); err != nil {
		return err
	}

	x16, y16 := thePlayer.Position()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-x16)/16, float64(-y16)/16)
	if err := ground.DrawImage(repeatedGophersImage, op); err != nil {
		return err
	}
	return nil
}

func drawGroundImage(screen *ebiten.Image, ground *ebiten.Image) error {

	op := &ebiten.DrawImageOptions{}
	if err := screen.DrawImage(ground, op); err != nil {
		return err
	}
	return nil
}

func update(screen *ebiten.Image) error {

	thePlayer.MoveForward()

	if err := updateGroundImage(groundImage); err != nil {
		return err
	}
	if err := drawGroundImage(screen, groundImage); err != nil {
		return err
	}

	msg := fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS())
	if err := ebitenutil.DebugPrint(screen, msg); err != nil {
		return err
	}
	return nil
}

func main() {

	var err error
	gophersImage, _, err = ebitenutil.NewImageFromFile("desert.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	w, h := gophersImage.Size()
	const repeat = 5
	repeatedGophersImage, err = ebiten.NewImage(w*repeat, h*repeat, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(w*i), float64(h*j))
			if err := repeatedGophersImage.DrawImage(gophersImage, op); err != nil {
				log.Fatal(err)
			}
		}
	}
	groundImage, err = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "infinite scroll"); err != nil {
		log.Fatal(err)
	}
}
