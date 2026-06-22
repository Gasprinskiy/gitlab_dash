package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Credentials struct {
		Host          string `yaml:"host"`
		PersonalToken string `yaml:"personal_token"`
	} `yaml:"credentials"`
	ProjectsData struct {
		ProjectIdList               []int  `yaml:"project_id_list"`
		IgnoreTestBranchCompareList []int  `yaml:"ignore_test_branch_compare_list"`
		TestBranchName              string `yaml:"test_branch_name"`
	} `yaml:"projects_data"`
}

func NewConfig() *Config {
	cfg := &Config{}

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("could not define home path: ", err)
	}

	pathToConfig := fmt.Sprintf("%s/.gitlab_dash/config.yaml", homePath)
	err = cleanenv.ReadConfig(pathToConfig, cfg)
	if err != nil {
		log.Fatal("could not parse config: ", err)
	}

	return cfg
}
