package scanner

type Config struct {
	Output         string
	IgnorePatterns []string
}

func DefaultConfig() Config {
	return Config{
		Output: "project_knowledge.md",
		IgnorePatterns: []string{
			".git",
			".idea",
			"node_modules",
			".idea",
			"vendor",
			"dist",
			"build",
			"go.sum",
		},
	}
}
