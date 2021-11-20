package env

type EnvData struct {
	Debug         bool
	PrivateKeys   []string `required:"true" split_words:"true"`
	DBUser        string   `required:"true" split_words:"true"`
	DBPassword    string   `required:"true" split_words:"true"`
	DBHost        string   `default:"127.0.0.1" split_words:"true"`
	DBPort        uint64   `default:"5432" split_words:"true"`
	DBName        string   `required:"true" split_words:"true"`
	QueueUser     string   `required:"true" split_words:"true"`
	QueuePassword string   `required:"true" split_words:"true"`
	QueueHost     string   `default:"127.0.0.1" split_words:"true"`
	QueuePort     uint64   `default:"5672" split_words:"true"`
	QueueName     string   `required:"true" split_words:"true"`
	NodeURL       string   `required:"true" split_words:"true"`
	QueryPK       string   `required:"true" split_words:"true"`
	ChainID       int64    `required:"true" split_words:"true"`
	ServerPort    uint64   `default:"8080" split_words:"true"`
	ServerOn      bool     `default:"true" split_words:"true"`
	WorkerOn      bool     `default:"true" split_words:"true"`
}
