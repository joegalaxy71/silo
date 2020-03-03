package types

type DbItem struct {
	Id       string `db:"id"`
	ServerId int    `db:"server_id"`
	Price    string `db:"price"`
	ChatId   int64  `db:"chat_id"`
	UserName string `db:"user_name"`
}

type UpdateInfo struct {
	Version string `yaml:"Version"`
	Build   int    `yaml:"Build"`
}

// =====================================================================================================================
// generated struct (generated from https://mholt.github.io/json-to-go/)

type SbItems struct {
	Hash         string `json:"hash"`
	MinMaxValues struct {
		MinPrice     int     `json:"minPrice"`
		MaxPrice     float64 `json:"maxPrice"`
		MinRAM       int     `json:"minRam"`
		MaxRAM       int     `json:"maxRam"`
		MinHDDSize   int     `json:"minHDDSize"`
		MaxHDDSize   int     `json:"maxHDDSize"`
		MinHDDCount  int     `json:"minHDDCount"`
		MaxHDDCount  int     `json:"maxHDDCount"`
		MinBenchmark int     `json:"minBenchmark"`
		MaxBenchmark int     `json:"maxBenchmark"`
	} `json:"minMaxValues"`
	Server []Server `json:"server"`
}

type Server struct {
	Key          int           `json:"key"`
	Name         string        `json:"name"`
	Description  []string      `json:"description"`
	CPU          string        `json:"cpu"`
	CPUBenchmark int           `json:"cpu_benchmark"`
	CPUCount     int           `json:"cpu_count"`
	IsHighio     bool          `json:"is_highio"`
	IsEcc        bool          `json:"is_ecc"`
	Traffic      string        `json:"traffic"`
	Dist         []string      `json:"dist"`
	Bandwith     int           `json:"bandwith"`
	RAM          int           `json:"ram"`
	Price        string        `json:"price"`
	PriceV       string        `json:"price_v"`
	RAMHr        string        `json:"ram_hr"`
	SetupPrice   string        `json:"setup_price"`
	HddSize      int           `json:"hdd_size"`
	HddCount     int           `json:"hdd_count"`
	HddHr        string        `json:"hdd_hr"`
	FixedPrice   bool          `json:"fixed_price"`
	NextReduce   int           `json:"next_reduce"`
	NextReduceHr string        `json:"next_reduce_hr"`
	Datacenter   []string      `json:"datacenter"`
	Specials     []interface{} `json:"specials"`
	SpecialHdd   string        `json:"specialHdd"`
	Freetext     string        `json:"freetext"`
}

type ServerInfo struct {
	Event  string
	Server Server
	DbItem DbItem
}
