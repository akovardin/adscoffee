package placements

type Cache struct {
	repo *Repo
}

func NewCache(repo *Repo) *Cache {
	return &Cache{
		repo: repo,
	}
}
