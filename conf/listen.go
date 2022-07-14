package conf

import "time"

type ListenYaml struct {
	Web          string        `yaml:"web"`
	WebMethod    string        `yaml:"webMethod"`
	WebNetwork   string        `yaml:"webNetwork"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	Identify     bool          `yaml:"identify"`
}

func (ly ListenYaml) GetReadTimeout() time.Duration {
	if ly.ReadTimeout.Seconds() < 1 {
		return 1 * time.Second
	} else {
		return ly.ReadTimeout
	}
}

func (ly ListenYaml) GetWriteTimeout() time.Duration {
	if ly.WriteTimeout.Seconds() < 1 {
		return 1 * time.Second
	} else {
		return ly.WriteTimeout
	}
}
