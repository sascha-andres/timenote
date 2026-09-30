package cache

import "time"

// NeedUpdate returns true if the projects or clients cache for workspace
// needs a refresh.
func (c *Cache) NeedUpdate(workspace int) bool {
	needProjects := needsRefresh(projectsPath(c.dir, workspace))
	needClients := needsRefresh(clientsPath(c.dir, workspace))
	return needProjects || needClients || c.AccountNeedUpdate()
}

func needsRefresh(path string) bool {
	var rec struct{ Meta MetaData }
	if err := readJSON(path, &rec); err != nil {
		return true
	}
	if rec.Meta.NextUpdate.IsZero() {
		return true
	}
	return time.Now().After(rec.Meta.NextUpdate)
}
