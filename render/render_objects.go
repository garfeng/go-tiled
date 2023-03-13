package render

import (
	"github.com/disintegration/imaging"
	"github.com/lafriks/go-tiled"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sort"
)

// RenderVisibleGroups renders all visible groups
func (r *Renderer) RenderVisibleGroups() error {
	for _, group := range r.m.Groups {
		if !group.Visible {
			continue
		}
		r._renderGroup(group)
	}
	return nil
}

// RenderGroup renders single group.
func (r *Renderer) RenderGroup(groupIdx int) error {
	group := r.m.Groups[groupIdx]
	return r._renderGroup(group)
}

func (r *Renderer) _renderGroup(group *tiled.Group) error {
	for _, layer := range group.Layers {
		if !layer.Visible {
			continue
		}
		err := r._renderLayer(layer)
		if err != nil {
			return err
		}
	}

	for _, objectGroup := range group.ObjectGroups {
		if !objectGroup.Visible {
			continue
		}
		err := r._renderObjectGroup(objectGroup)
		if err != nil {
			return err
		}
	}

	return nil
}

// RenderVisibleLayersAndObjectGroups render all layers and object groups, layer first, objectGroup second
// so the order may be incorrect,
// you may put them into different groups, then call RenderVisibleGroups
func (r *Renderer) RenderVisibleLayersAndObjectGroups() error {
	// TODO: The order maybe incorrect

	err := r.RenderVisibleLayers()
	if err != nil {
		return err
	}
	return r.RenderVisibleObjectGroups()
}

// RenderVisibleObjectGroups renders all visible object groups
func (r *Renderer) RenderVisibleObjectGroups() error {
	for i, layer := range r.m.ObjectGroups {
		if layer.Visible {
			err := r.RenderObjectGroup(i)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RenderObjectGroup renders a single object group
func (r *Renderer) RenderObjectGroup(i int) error {
	layer := r.m.ObjectGroups[i]
	return r._renderObjectGroup(layer)
}

func (r *Renderer) _renderObjectGroup(objectGroup *tiled.ObjectGroup) error {
	objs := objectGroup.Objects
	objs = SortAny(objs, sortObjs)
	for _, obj := range objs {
		err := r.renderOneObjectUpperTile(objectGroup, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

// RenderGroupObjectGroup renders single object group in a certain group.
func (r *Renderer) RenderGroupObjectGroup(groupIdx, objectGroupId int) error {
	group := r.m.Groups[groupIdx]
	layer := group.ObjectGroups[objectGroupId]
	return r._renderObjectGroup(layer)
}

func (r *Renderer) renderOneObjectUpperTile(layer *tiled.ObjectGroup, o *tiled.Object) error {
	if !o.Visible {
		return nil
	}

	if o.GID == 0 {
		// TODO: o.GID == 0
		return nil
	}

	tile, err := r.m.TileGIDToTile(o.GID)
	if err != nil {
		return err
	}

	img, err := r.GetTileImage(tile)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	srcSize := bounds.Size()
	dstSize := image.Pt(int(o.Width), int(o.Height))

	if !srcSize.Eq(dstSize) {
		img = imaging.Resize(img, dstSize.X, dstSize.Y, imaging.NearestNeighbor)
	}

	if o.Rotation != 0 {
		img = imaging.Rotate(img, -o.Rotation, color.RGBA{})
	}

	n := math.Floor(o.Y / float64(r.m.TileHeight))
	upperY := n * float64(r.m.TileHeight)
	leftTop := o.Y - o.Height

	upperHeight := upperY - leftTop

	bounds = img.Bounds()
	bounds.Max.Y = int(upperHeight)

	pos := bounds.Add(image.Pt(int(o.X), int(leftTop)))

	if layer.Opacity < 1 {
		mask := image.NewUniform(color.Alpha{uint8(layer.Opacity * 255)})

		draw.DrawMask(r.Result, pos, img, img.Bounds().Min, mask, mask.Bounds().Min, draw.Over)
	} else {
		draw.Draw(r.Result, pos, img, img.Bounds().Min, draw.Over)
	}

	return nil
}

func (r *Renderer) renderOneObject(layer *tiled.ObjectGroup, o *tiled.Object) error {
	if !o.Visible {
		return nil
	}

	if o.GID == 0 {
		// TODO: o.GID == 0
		return nil
	}

	tile, err := r.m.TileGIDToTile(o.GID)
	if err != nil {
		return err
	}

	img, err := r.GetTileImage(tile)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	srcSize := bounds.Size()
	dstSize := image.Pt(int(o.Width), int(o.Height))

	if !srcSize.Eq(dstSize) {
		img = imaging.Resize(img, dstSize.X, dstSize.Y, imaging.NearestNeighbor)
	}

	if o.Rotation != 0 {
		img = imaging.Rotate(img, -o.Rotation, color.RGBA{})
	}

	bounds = img.Bounds()
	pos := bounds.Add(image.Pt(int(o.X), int(o.Y-o.Height)))

	if layer.Opacity < 1 {
		mask := image.NewUniform(color.Alpha{uint8(layer.Opacity * 255)})

		draw.DrawMask(r.Result, pos, img, img.Bounds().Min, mask, mask.Bounds().Min, draw.Over)
	} else {
		draw.Draw(r.Result, pos, img, img.Bounds().Min, draw.Over)
	}

	return nil
}

func sortObjs(a, b *tiled.Object) bool {
	if a.Y != b.Y {
		return a.Y < b.Y
	}

	return a.X < b.X
}

func SortAny[T any](data []T, lessMethod func(a, b T) bool) []T {
	s := &Sortable[T]{
		Data:       data,
		LessMethod: lessMethod,
	}
	sort.Sort(s)
	return s.Data
}

type Sortable[T any] struct {
	Data       []T
	LessMethod func(a, b T) bool
}

func (s *Sortable[T]) Swap(i, j int) {
	tmp := (s.Data)[i]
	(s.Data)[i] = (s.Data)[j]
	(s.Data)[j] = tmp
}

func (s *Sortable[T]) Less(i, j int) bool {
	return s.LessMethod(s.Data[i], s.Data[j])
}

func (s *Sortable[T]) Len() int {
	return len(s.Data)
}
