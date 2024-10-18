package w365_tt

import "fmt"

// Get the full tag of a group in the form "class.class-group", or just
// "class" if there is no class-group.

func groupTagFull(idmap IdMap, gid string) string {
	gnode, ok := idmap.Id2Node[gid]
	if !ok {
		return fmt.Sprintf("*!%s!*", gid)
	}
	if n, ok := gnode.(Group); ok {
		c, ok := idmap.Group2Class[gid]
		if !ok {
			return fmt.Sprintf("*(%s)*", gid)
		}
		return fmt.Sprintf("%d%s.%s", c.Level, c.Letter, n.Shortcut)
	} else if n, ok := gnode.(Class); ok {
		return fmt.Sprintf("%d%s", n.Level, n.Letter)
	} else {
		return fmt.Sprintf("*?%s?*", gid)
	}
}
