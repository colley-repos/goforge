// Package spatial provides grid and spatial query abstractions.
//
// The default implementation is a 2D tile grid with BFS pathfinding,
// occupancy tracking, and cover support — a direct port of Dystopia's GridManager.
package spatial

// CoverLevel represents the cover provided by a tile.
type CoverLevel int

const (
	CoverNone CoverLevel = iota
	CoverHalf
	CoverFull
)

// Cell holds data for a single grid tile.
type Cell struct {
	Walkable bool
	Occupant uint64     // entity ID, 0 = empty
	Cover    CoverLevel
}

// Grid is a 2D tile-based spatial structure.
type Grid struct {
	Width  int
	Height int
	Cells  []Cell // row-major: index = row * Width + col
}

// NewGrid creates a grid with all cells walkable and empty.
func NewGrid(width, height int) *Grid {
	cells := make([]Cell, width*height)
	for i := range cells {
		cells[i].Walkable = true
	}
	return &Grid{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// InBounds returns true if (col, row) is within the grid.
func (g *Grid) InBounds(col, row int) bool {
	return col >= 0 && col < g.Width && row >= 0 && row < g.Height
}

// CellAt returns a pointer to the cell at (col, row), or nil if out of bounds.
func (g *Grid) CellAt(col, row int) *Cell {
	if !g.InBounds(col, row) {
		return nil
	}
	return &g.Cells[row*g.Width+col]
}

// IsWalkable returns true if (col, row) is in bounds, walkable, and unoccupied.
func (g *Grid) IsWalkable(col, row int) bool {
	cell := g.CellAt(col, row)
	if cell == nil {
		return false
	}
	return cell.Walkable && cell.Occupant == 0
}

// SetOccupant places an entity on a tile (0 to clear).
func (g *Grid) SetOccupant(col, row int, entityID uint64) {
	if cell := g.CellAt(col, row); cell != nil {
		cell.Occupant = entityID
	}
}

// SetWalkable marks a tile as walkable or blocked.
func (g *Grid) SetWalkable(col, row int, walkable bool) {
	if cell := g.CellAt(col, row); cell != nil {
		cell.Walkable = walkable
	}
}

// SetCover sets the cover level for a tile.
func (g *Grid) SetCover(col, row int, cover CoverLevel) {
	if cell := g.CellAt(col, row); cell != nil {
		cell.Cover = cover
	}
}

// GetCover returns the cover level at (col, row).
func (g *Grid) GetCover(col, row int) CoverLevel {
	if cell := g.CellAt(col, row); cell != nil {
		return cell.Cover
	}
	return CoverNone
}

// Neighbors returns the 4 cardinal neighbors of (col, row) that are in bounds.
func (g *Grid) Neighbors(col, row int) [][2]int {
	dirs := [4][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	result := make([][2]int, 0, 4)
	for _, d := range dirs {
		nc, nr := col+d[0], row+d[1]
		if g.InBounds(nc, nr) {
			result = append(result, [2]int{nc, nr})
		}
	}
	return result
}

// BFSReachable returns all tiles reachable from (startCol, startRow)
// within maxSteps moves, respecting walkability and occupancy.
// Returns a map of [col, row] → distance.
func (g *Grid) BFSReachable(startCol, startRow, maxSteps int) map[[2]int]int {
	type node struct {
		col, row, dist int
	}

	visited := make(map[[2]int]int)
	queue := []node{{startCol, startRow, 0}}
	visited[[2]int{startCol, startRow}] = 0

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.dist >= maxSteps {
			continue
		}

		for _, n := range g.Neighbors(cur.col, cur.row) {
			pos := [2]int{n[0], n[1]}
			if _, seen := visited[pos]; seen {
				continue
			}
			if !g.IsWalkable(n[0], n[1]) {
				continue
			}
			dist := cur.dist + 1
			visited[pos] = dist
			queue = append(queue, node{n[0], n[1], dist})
		}
	}

	return visited
}

// BFSPath returns the shortest path from (startCol, startRow) to (goalCol, goalRow),
// or nil if no path exists. The path includes the goal but not the start.
func (g *Grid) BFSPath(startCol, startRow, goalCol, goalRow int) [][2]int {
	type node struct {
		col, row int
	}

	if !g.InBounds(goalCol, goalRow) {
		return nil
	}

	start := node{startCol, startRow}
	goal := node{goalCol, goalRow}

	if start == goal {
		return [][2]int{}
	}

	parent := map[node]node{start: start}
	queue := []node{start}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur == goal {
			// Reconstruct path
			path := [][2]int{}
			at := goal
			for at != start {
				path = append([][2]int{{at.col, at.row}}, path...)
				at = parent[at]
			}
			return path
		}

		for _, n := range g.Neighbors(cur.col, cur.row) {
			next := node{n[0], n[1]}
			if _, seen := parent[next]; seen {
				continue
			}
			// Goal doesn't need to be walkable (might be occupied by target)
			if next != goal && !g.IsWalkable(n[0], n[1]) {
				continue
			}
			parent[next] = cur
			queue = append(queue, next)
		}
	}

	return nil // no path
}

// ManhattanDistance returns the Manhattan distance between two grid positions.
func ManhattanDistance(c1, r1, c2, r2 int) int {
	dc := c1 - c2
	dr := r1 - r2
	if dc < 0 {
		dc = -dc
	}
	if dr < 0 {
		dr = -dr
	}
	return dc + dr
}
