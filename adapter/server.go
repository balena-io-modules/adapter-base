package adapter

type Server struct {
	concurrency int
	workers     map[string]*Worker
}

func CreateServer(concurrency int, verbose bool) *Server {
	server := &Server{
		concurrency: concurrency,
		workers:     make(map[string]*Worker),
	}

	return server
}
