package conf

import (
	"github.com/jinzhu/configor"
)

type userDataBase struct {
	HostAndPort string `required:"true"`
	Name        string `required:"true"`
	User        string
	Passwd      string
}

type contentDatabase struct {
	Name     string `required:"true"`
	Protocol string `required:"true"`
	Addr     string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
}

type appConfig struct {
	UserDB          userDataBase    `required:"true"`
	ContentDB       contentDatabase `required:"true"`
	Host            string          `required:"true"`
	Port            uint            `required:"true"`
	Debug           bool            `required:"true"`
	Secret          string          `required:"true"`
	ContentServer   string          `required:"true"`
	BookImagesURL   string          `required:"true"`
	PaperImagesURL  string          `required:"true"`
	HTMLURL         string          `required:"true"`
	ProblemDocURL   string          `required:"true"`
	StudentFilesDir string          `required:"true"`
	StudentFilesURL string          `required:"true"`
	FilesServer     string          `required:"true"`
}

var (
	// AppConfig is the config for app
	AppConfig = appConfig{}
)

func init() {
	configor.Load(&AppConfig, "./conf/config.yml")
}
