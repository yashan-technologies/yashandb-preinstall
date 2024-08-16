package confdef

import "encoding/json"

type YashanDBUser struct {
	SkipCreatingUser bool     `toml:"skip_creating_user"`
	User             string   `toml:"user"`
	Group            string   `toml:"group"`
	Home             string   `toml:"home"`
	AdditionalGroups []string `toml:"additional_groups"`
	Limits           []string `toml:"limits"`
}

type Commands struct {
	Bash         string   `toml:"bash"`
	Useradd      string   `toml:"useradd"`
	Usermod      string   `toml:"usermod"`
	Groupadd     string   `toml:"groupadd"`
	Systemctl    string   `toml:"systemctl"`
	IPTables     string   `toml:"iptables"`
	Sudo         string   `toml:"sudo"`
	Echo         string   `toml:"echo"`
	Cat          string   `toml:"cat"`
	Su           string   `toml:"su"`
	Timedatectl  string   `toml:"timedatectl"`
	Ulimit       string   `toml:"ulimit"`
	Sysctl       string   `toml:"sysctl"`
	Swapoff      string   `toml:"swapoff"`
	UpdateGrub   string   `toml:"update_grub"`
	GrubMkConfig []string `toml:"grub_mkconfig"`
	GrubCgf      string   `toml:"grub_cfg"`
}

type Hardware struct {
	CPUMinCores             int      `toml:"cpu_min_cores"`
	MemoryMinGB             int      `toml:"memory_min_gb"`
	InstallPathMinFreeGB    int      `toml:"install_path_min_free_gb"`
	InstallPathFsTypes      []string `toml:"install_path_fs_types"`
	NetworkMinBandWidthMbps int      `toml:"network_min_bandwidth_mbps"`
}

type Software struct {
	CentosMinVersion string `toml:"centos_min_version"`
	Charset          string `toml:"charset"`
}

type Limit struct {
	Hardware Hardware `toml:"hardware"`
	Software Software `toml:"software"`
}

type HostSetting struct {
	Sysctl        []string `toml:"sysctl"`
	Timezone      string   `toml:"timezone"`
	DiskScheduler string   `toml:"disk_scheduler"`
}

type Fio struct {
	Size      string   `toml:"size"`
	RWMode    []string `toml:"rw"`
	BlockSize string   `toml:"bs"`
	NumJobs   string   `toml:"numjobs"`
	RunTime   string   `toml:"runtime"`
	IODepth   string   `toml:"iodepth"`
	Direct    string   `toml:"direct"`
}

type Config struct {
	LogLevel     string       `toml:"log_level"`
	YashanDBUser YashanDBUser `toml:"yashandb_user"`
	HostSetting  HostSetting  `toml:"host_setting"`
	Fio          Fio          `toml:"fio"`
	Commands     Commands     `toml:"commands"`
	Limit        Limit        `toml:"limit"`
}

func (c Config) ToJSON() string {
	data, _ := json.MarshalIndent(c, "", "    ")
	return string(data)
}

func Conf() Config {
	return _conf
}
