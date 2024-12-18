package scanner

var DefaultIgnorePatterns = []string{
	".git",
	"go.sum",
}

type Config struct {
	Output         string
	IgnorePatterns []string
}

func DefaultConfig() Config {
	return Config{
		Output:         "project_knowledge.md",
		IgnorePatterns: DefaultIgnorePatterns,
	}
}
