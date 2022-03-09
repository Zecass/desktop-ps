package process

func buchheim(tree *ProcessTree) *ProcessTree {
	dt := firstWalk(tree)
	secondWalk(dt, 0, 0)
	// if min != nil && *min < 0 {
	// 	thirdWalk(tree, -*min)
	// }
	return dt
}

func firstWalk(tree *ProcessTree) *ProcessTree {
	if len(tree.Childrens) == 0 {
		if tree.leftMostSibling() != nil {
			tree.X = tree.leftBrother().X + 1
		} else {
			tree.X = 0
		}
	} else {
		defaultAncestor := tree.Childrens[0]

		for _, c := range tree.Childrens {
			firstWalk(c)
			defaultAncestor = apportion(c, defaultAncestor)
		}

		executeShifts(tree)

		midpoint := (tree.Childrens[0].X + tree.Childrens[len(tree.Childrens)-1].X) / 2

		b := tree.leftBrother()
		if b != nil {
			tree.X = b.X + 1
			tree.mod = tree.X - midpoint
		} else {
			tree.X = midpoint
		}
	}

	return tree
}

func apportion(v *ProcessTree, defaultAncestor *ProcessTree) *ProcessTree {
	w := v.leftBrother()
	if w != nil {
		vir := v
		vor := v

		vil := w
		vol := v.leftMostSibling()

		sir := v.mod
		sor := v.mod

		sil := vil.mod
		sol := vol.mod

		for vil.nextRight() != nil && vir.nextLeft() != nil {
			vil = vil.nextRight()
			vir = vir.nextLeft()

			vol = vol.nextLeft()
			vor = vor.nextRight()

			vor.ancestor = v

			shift := (vil.X + sil) - (vir.X + sir) + 1

			if shift > 0 {
				a := ancestor(vil, v, defaultAncestor)
				moveSubtree(a, v, shift)

				sir += shift
				sor += shift
			}

			sil += vil.mod
			sir += vir.mod

			sol += vol.mod
			sor += vor.mod
		}

		if vil.nextRight() != nil && vor.nextRight() == nil {
			vor.link = vil.nextRight()
			vor.mod += sil - sor
		} else {
			if vir.nextLeft() != nil && vol.nextLeft() == nil {
				vol.link = vir.nextLeft()
				vol.mod += sir - sol
			}

			defaultAncestor = v
		}
	}

	return defaultAncestor
}

func moveSubtree(wl, wr *ProcessTree, shift float64) {
	subtrees := float64(wr.number - wl.number)

	wr.change -= shift / subtrees
	wr.shift += shift

	wl.change += shift / subtrees
	wr.X += shift
	wr.mod += shift
}

func executeShifts(v *ProcessTree) {
	shift := 0.0
	change := 0.0

	s := make([]*ProcessTree, len(v.Childrens))

	copy(s, v.Childrens)

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	for _, w := range s {
		w.X += shift
		w.mod += shift

		change += w.change
		shift += w.shift + change
	}

}

func ancestor(vil, v, defaultAncestor *ProcessTree) *ProcessTree {
	for _, c := range v.Parent.Childrens {
		if c == vil.ancestor {
			return c
		}
	}

	return defaultAncestor
}

func secondWalk(v *ProcessTree, m, depth float64) {
	v.X += m
	v.Depth = depth

	for _, w := range v.Childrens {
		secondWalk(w, m+v.mod, depth+1)
	}
}
