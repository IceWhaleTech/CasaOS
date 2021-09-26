package docker_base

type MysqlConfig struct {
	DataBaseHost     string `json:"database_host"`
	DataBasePort     string `json:"database_port"`
	DataBaseUser     string `json:"database_user"`
	DataBasePassword string `json:"data_base_password"`
	DataBaseDB       string `json:"data_base_db"`
}


