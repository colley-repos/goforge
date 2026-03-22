package spatial

import "testing"

func TestGridBasics(t *testing.T) {
	g := NewGrid(8, 8)

	if !g.InBounds(0, 0) || !g.InBounds(7, 7) {
		t.Error("corner cells should be in bounds")
	}
	if g.InBounds(-1, 0) || g.InBounds(8, 0) {
		t.Error("out of bounds should return false")
	}

	if !g.IsWalkable(3, 3) {
		t.Error("all cells should start walkable")
	}

	g.SetWalkable(3, 3, false)
	if g.IsWalkable(3, 3) {
		t.Error("should be blocked after SetWalkable(false)")
	}

	g.SetOccupant(4, 4, 99)
	if g.IsWalkable(4, 4) {
		t.Error("occupied cell should not be walkable")
	}

	g.SetOccupant(4, 4, 0)
	if !g.IsWalkable(4, 4) {
		t.Error("cleared occupant should restore walkability")
	}
}

func TestGridCover(t *testing.T) {
	g := NewGrid(4, 4)
	g.SetCover(2, 2, CoverHalf)
	g.SetCover(3, 3, CoverFull)

	if g.GetCover(2, 2) != CoverHalf {
		t.Error("expected half cover")
	}
	if g.GetCover(3, 3) != CoverFull {
		t.Error("expected full cover")
	}
	if g.GetCover(0, 0) != CoverNone {
		t.Error("expected no cover")
	}
}

func TestBFSReachable(t *testing.T) {
	g := NewGrid(5, 5)

	reachable := g.BFSReachable(2, 2, 2)

	// Should include the center
	if _, ok := reachable[[2]int{2, 2}]; !ok {
		t.Error("start should be in reachable set")
	}

	// Should include adjacent tiles
	if _, ok := reachable[[2]int{3, 2}]; !ok {
		t.Error("(3,2) should be reachable with 2 steps")
	}

	// Should NOT include tiles 3+ steps away
	if _, ok := reachable[[2]int{4, 4}]; ok {
		t.Error("(4,4) should NOT be reachable with 2 steps")
	}

	// Block a path
	g.SetWalkable(3, 2, false)
	reachable2 := g.BFSReachable(2, 2, 1)
	if _, ok := reachable2[[2]int{3, 2}]; ok {
		t.Error("blocked cell should not be reachable")
	}
}

func TestBFSPath(t *testing.T) {
	g := NewGrid(5, 5)

	path := g.BFSPath(0, 0, 3, 0)
	if path == nil {
		t.Fatal("path should exist")
	}
	if len(path) != 3 {
		t.Fatalf("expected path length 3, got %d: %v", len(path), path)
	}
	// Should end at goal
	if path[len(path)-1] != [2]int{3, 0} {
		t.Errorf("path should end at (3,0), got %v", path[len(path)-1])
	}

	// Path to self
	path = g.BFSPath(2, 2, 2, 2)
	if path == nil || len(path) != 0 {
		t.Errorf("path to self should be empty slice, got %v", path)
	}

	// Blocked path
	g.SetWalkable(2, 0, false)
	g.SetWalkable(2, 1, false)
	g.SetWalkable(2, 2, false)
	g.SetWalkable(2, 3, false)
	g.SetWalkable(2, 4, false)
	path = g.BFSPath(0, 0, 4, 0)
	if path != nil {
		t.Errorf("path should be nil when fully blocked, got %v", path)
	}
}

func TestManhattanDistance(t *testing.T) {
	if d := ManhattanDistance(0, 0, 3, 4); d != 7 {
		t.Errorf("expected 7, got %d", d)
	}
	if d := ManhattanDistance(5, 5, 5, 5); d != 0 {
		t.Errorf("expected 0, got %d", d)
	}
}

func TestGridNeighbors(t *testing.T) {
	g := NewGrid(3, 3)

	// Center has 4 neighbors
	n := g.Neighbors(1, 1)
	if len(n) != 4 {
		t.Errorf("center should have 4 neighbors, got %d", len(n))
	}

	// Corner has 2 neighbors
	n = g.Neighbors(0, 0)
	if len(n) != 2 {
		t.Errorf("corner should have 2 neighbors, got %d", len(n))
	}
}
