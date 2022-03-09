package image

import (
	"math"

	"github.com/Zecass/desktop-ps/process"
	"github.com/Zecass/desktop-ps/settings"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

type dot struct {
	nx, ny, r float64

	x, y float64
}

type link struct {
	nx1, ny1, r1 float64
	nx2, ny2, r2 float64
}

type label struct {
	angle float64
	x, y  float64
	text  string
}

type circle struct {
	dots   []*dot
	links  []*link
	labels []*label
}

type nVector2 struct {
	nx, ny, r float64
}

var (
	config *settings.Settings

	halfW float64
	halfH float64

	rotation float64

	face *canvas.FontFace
)

func GenerateImageFromProcess(ps *process.ProcessTree, settings *settings.Settings) (err error) {
	config = settings

	face, err = setFontFace()
	if err != nil {
		return err
	}

	halfW = config.ResolutionX / 2
	halfH = config.ResolutionY / 2

	c := canvas.New(config.ResolutionX, config.ResolutionY)
	ctx := canvas.NewContext(c)

	ctx.SetFillColor(config.BackgroundColorRGBA)
	ctx.DrawPath(0, 0, canvas.Rectangle(config.ResolutionX, config.ResolutionY))

	ctx.SetStrokeColor(config.Link.ColorRGBA)
	ctx.SetStrokeWidth(config.Link.Width)

	rotation = 360 / (ps.MaxX - ps.MinX + 1)
	rotation = rotation * (math.Pi / 180)

	drawTree(ctx, ps)

	renderers.Write("wallpaper.png", c)

	return nil
}

func drawTree(ctx *canvas.Context, ps *process.ProcessTree) {
	circle := &circle{
		dots:   []*dot{},
		links:  []*link{},
		labels: []*label{},
	}
	sketch(ps, nil, circle)

	draw(ctx, circle)
}

func sketch(ps *process.ProcessTree, parentPos *nVector2, circle *circle) {
	if ps.ProcessName != "fakeRoot" {
		r := 150 * ps.Depth
		r = (config.ResolutionY * r) / 1800

		angle := ps.X * rotation
		nx := math.Cos(angle)
		ny := math.Sin(angle)

		d := &dot{
			nx: nx,
			ny: ny,
			r:  r,

			x: nx*r + halfW,
			y: ny*r + halfH,
		}

		circle.dots = append(circle.dots, d)

		if parentPos != nil {
			link := &link{
				nx1: parentPos.nx,
				ny1: parentPos.ny,
				r1:  parentPos.r,

				nx2: nx,
				ny2: ny,
				r2:  r,
			}

			circle.links = append(circle.links, link)
		}

		label := &label{
			angle: angle,
			x:     nx*(r+7) + halfW,
			y:     ny*(r+7) + halfH,

			text: ps.ProcessName,
		}

		circle.labels = append(circle.labels, label)

		parentPos = &nVector2{nx, ny, r}
	}

	for _, c := range ps.Childrens {
		sketch(c, parentPos, circle)
	}
}

func draw(ctx *canvas.Context, circle *circle) {
	frame := canvas.Rectangle(config.ResolutionX, config.ResolutionY)

	for _, l := range circle.links {
		drawLink(frame, l)
	}

	ctx.DrawPath(0, 0, frame)

	for _, l := range circle.labels {
		drawLabel(ctx, l)
	}

	for _, d := range circle.dots {
		drawDot(ctx, d)
	}
}

func drawDot(ctx *canvas.Context, d *dot) {
	x := d.nx*d.r + halfW
	y := d.ny*d.r + halfH

	ctx.SetStrokeColor(config.Dot.BorderColorRGBA)
	ctx.SetStrokeWidth(config.Dot.BorderWidth)

	ctx.SetFillColor(config.Dot.ColorRGBA)
	ctx.DrawPath(x, y, canvas.Circle(config.Dot.Radius))
}

func drawLabel(ctx *canvas.Context, l *label) {
	textAllign := canvas.Left
	ctx.RotateAbout(l.angle*180/math.Pi, l.x, l.y)
	if l.x < halfW {
		ctx.ReflectYAbout(l.y)
		ctx.ReflectXAbout(l.x)
		textAllign = canvas.Right
	}

	ctx.DrawText(l.x, l.y, canvas.NewTextBox(face, l.text, 0, 0, textAllign, canvas.Center, 0, 0))
	ctx.ResetView()
}

func drawLink(ctx *canvas.Path, l *link) {
	x1 := l.nx1*l.r1 + halfW
	y1 := l.ny1*l.r1 + halfH

	x2 := l.nx2*l.r2 + halfW
	y2 := l.ny2*l.r2 + halfH

	cp := (l.r2 - l.r1) / 2

	cx1 := l.nx1*(l.r1+cp) + halfW
	cy1 := l.ny1*(l.r1+cp) + halfH

	cx2 := l.nx2*(l.r2-cp) + halfW
	cy2 := l.ny2*(l.r2-cp) + halfH

	ctx.MoveTo(x1, y1)
	ctx.CubeTo(cx1, cy1, cx2, cy2, x2, y2)
}

func setFontFace() (*canvas.FontFace, error) {
	font := canvas.NewFontFamily("OpenSans")
	err := font.LoadFontFile("font/OpenSans-Medium.ttf", canvas.FontMedium)
	if err != nil {
		return &canvas.FontFace{}, err
	}

	face = font.Face((config.ResolutionY*30)/1800, config.Font.ColorRGBA, canvas.FontMedium, canvas.FontNormal)

	return face, nil
}
