package transaction

type transaction struct {
	config *Config
}

func New(config *Config) *transaction {
	return &transaction{
		config: config,
	}
}

func (t *transaction) Save() {}

func (t *transaction) GetAll() {}

func (t *transaction) Update() {}

func (t *transaction) Delete() {}
