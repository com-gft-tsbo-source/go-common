package microservice

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/com-gft-tsbo-source/go-common/ms-framework/dispatcher"
)

// DBConfiguration ...
type DBConfiguration struct {
	DBName string `json:"db"`
}

// IDBConfiguration ...
type IDBConfiguration interface {
	GetDBName() string
}

// ServiceConfiguration ...
type ServiceConfiguration struct {
	// Name     string `json:"name"`
	// Hostname string `json:"hostname"`
	// Version  string `json:"version"`
}

// IServiceConfiguration ...
type IServiceConfiguration interface {
	// GetName() string
	// GetHostname() string
	// GetVersion() string
}

// FileConfiguration ...
type FileConfiguration struct {
	ConfigurationFile string
}

// IFileConfiguration ...
type IFileConfiguration interface {
	GetConfigurationFile() string
}

// Configuration ...
type Configuration struct {
	dispatcher.Configuration
	DBConfiguration
	ServiceConfiguration
	FileConfiguration
}

// IConfiguration ...
type IConfiguration interface {
	dispatcher.IConfiguration
	IDBConfiguration
	IServiceConfiguration
	IFileConfiguration
}

// GetDBName ...
func (cfg DBConfiguration) GetDBName() string { return cfg.DBName }

// // GetName ...
// func (cfg ServiceConfiguration) GetName() string { return cfg.Name }

// // GetHostname ...
// func (cfg ServiceConfiguration) GetHostname() string { return cfg.Hostname }

// // GetVersion ...
// func (cfg ServiceConfiguration) GetVersion() string { return cfg.Version }

// GetConfigurationFile ...
func (cfg FileConfiguration) GetConfigurationFile() string { return cfg.ConfigurationFile }

// ---------------------------------------------------------------------------

// InitConfigurationFromArgs ...
func InitConfigurationFromArgs(cfg *Configuration, args []string, flagset *flag.FlagSet) {
	var requestHeaders dispatcher.HeaderList
	var responseHeaders dispatcher.HeaderList
	var copyHeaders dispatcher.HeaderList
	var configurationFile string

	DefaultPort := 8080
	DefaultHost := "0.0.0.0"

	if flagset == nil {
		flagset = flag.NewFlagSet("ms", flag.PanicOnError)
	}

	flagset.Var(&requestHeaders, "requestheader", "Add these headers for outgoing requests.")
	flagset.Var(&responseHeaders, "responseheader", "Add these headers for outgoing replies.")
	flagset.Var(&copyHeaders, "copyheader", "Add these headers for outgoing replies.")
	phost := flagset.String("host", "", "Ip address to listen on.")
	pport := flagset.Int("port", -1, "Listen port.")
	pname := flagset.String("name", "", "Name of the service.")
	phostname := flagset.String("hostname", "", "Hostname of the service.")
	pversion := flagset.String("version", "", "Version of the service.")
	pnamespace := flagset.String("namespace", "", "Prefix for all urls.")
	pdb := flagset.String("db", "", "Database connection")
	pcertchainfile := flagset.String("cert", "", "Certificate chain (host cert + all sigining CAs)")
	pkeyfile := flagset.String("key", "", "Private key file.")
	pcafile := flagset.String("ca", "", "CA chains.")
	pconfig := flagset.String("config", "", "Configuration file.")
	// pmaxTcpConnections := flagset.Int("maxtcpconnections", -1, "Maximum of parallel connections to accept on TCP level.")
	pmaxConnections := flagset.Int("maxconnections", -1, "Maximum of parallel connections to accept.")
	pdelayReply := flagset.Int("delayreply", -1, "Slow down replying by this amount of ms.")
	pclientTimeout := flagset.Int("clienttimeout", 1500, "Timeout of HTTP client in ms.")
	plogfile := flagset.String("logfile", "", "Logfile (empty=stdout).")
	pnometrics := flagset.Bool("nometrics", false, "Don't report metrics..")
	ppasswordfile := flagset.String("passwordfile", "", "User/password list.")

	flagset.Parse(os.Args[1:])

	if len(responseHeaders) > 0 {
		cfg.ResponseHeaderStrings = responseHeaders
	}

	if len(requestHeaders) > 0 {
		cfg.RequestHeaderStrings = requestHeaders
	}

	if len(copyHeaders) > 0 {
		cfg.CopyHeaderStrings = copyHeaders
	}

	if *pport >= 0 {
		cfg.Port = *pport
	} else {
		ev := os.Getenv("MS_PORT")
		if len(ev) > 0 {
			cfg.Port, _ = strconv.Atoi(ev)
		}
	}

	if len(*phost) > 0 {
		cfg.Host = *phost
	}
	if len(cfg.Host) == 0 {
		cfg.Host = os.Getenv("MS_HOST")
	}

	if len(*pname) > 0 {
		cfg.Name = *pname
	}
	if len(cfg.Name) == 0 {
		cfg.Name = os.Getenv("MS_NAME")
	}

	if len(*phostname) > 0 {
		cfg.Hostname = *phostname
	}
	if len(cfg.Hostname) == 0 {
		cfg.Hostname = os.Getenv("MS_HOSTNAME")
	}

	if len(*pversion) > 0 {
		cfg.Version = *pversion
	}
	if len(cfg.Version) == 0 {
		cfg.Version = os.Getenv("MS_VERSION")
	}

	if len(*pnamespace) > 0 {
		cfg.Namespace = *pnamespace
	}
	if len(cfg.Namespace) == 0 {
		cfg.Namespace = os.Getenv("MS_NAMESPACE")
	}

	if len(*pdb) > 0 {
		cfg.DBName = *pdb
	}
	if len(cfg.DBName) == 0 {
		cfg.DBName = os.Getenv("MS_DATABASE")
	}

	if len(*pcertchainfile) > 0 {
		cfg.CertChainFile = *pcertchainfile
	}
	if len(cfg.CertChainFile) == 0 {
		cfg.CertChainFile = os.Getenv("MS_CERTCHAINFILE")
	}

	if len(*pkeyfile) > 0 {
		cfg.KeyFile = *pkeyfile
	}
	if len(cfg.KeyFile) == 0 {
		cfg.KeyFile = os.Getenv("MS_KEYFILE")
	}

	if len(*pcafile) > 0 {
		cfg.CAFile = *pcafile
	}
	if len(cfg.CAFile) == 0 {
		cfg.CAFile = os.Getenv("MS_CAFILE")
	}

	if len(*pconfig) > 0 {
		configurationFile = *pconfig
	}
	if len(configurationFile) == 0 {
		configurationFile = os.Getenv("MS_CONFIG")
	}

	// if *pmaxTcpConnections > 0 {
	// 	cfg.MaxTcpConnections = *pmaxTcpConnections
	// } else {
	// 	ev := os.Getenv("MS_MAXTCPCONNECTIONS")
	// 	if len(ev) > 0 {
	// 		cfg.MaxTcpConnections, _ = strconv.Atoi(ev)
	// 	}
	// }

	if *pmaxConnections > 0 {
		cfg.MaxConnections = *pmaxConnections
	} else {
		ev := os.Getenv("MS_MAXCONNECTIONS")
		if len(ev) > 0 {
			cfg.MaxConnections, _ = strconv.Atoi(ev)
		}
	}

	if *pdelayReply > 0 {
		cfg.DelayReply = *pdelayReply
	} else {
		ev := os.Getenv("MS_DELAYREPLY")
		if len(ev) > 0 {
			cfg.DelayReply, _ = strconv.Atoi(ev)
		}
	}

	if *pclientTimeout > 0 {
		cfg.ClientTimeout = *pclientTimeout
	} else {
		ev := os.Getenv("MS_CLIENTTIMEOUT")
		if len(ev) > 0 {
			cfg.ClientTimeout, _ = strconv.Atoi(ev)
		}
	}

	if len(*plogfile) > 0 {
		cfg.Logfile = *plogfile
	}
	if len(cfg.Logfile) == 0 {
		cfg.Logfile = os.Getenv("MS_LOGFILE")
	}

	if *pnometrics {
		cfg.NoMetrics = true
	} else {
		ev := os.Getenv("MS_NO_METRICS")
		if len(ev) > 0 {
			cfg.NoMetrics = true
		}
	}

	if len(*ppasswordfile) > 0 {
		cfg.Passwordfile = *ppasswordfile
	}
	if len(cfg.Passwordfile) == 0 {
		cfg.Passwordfile = os.Getenv("MS_PASSWORDFILE")
	}

	if len(configurationFile) > 0 {

		cfg.ConfigurationFile = configurationFile
		file, err := os.Open(configurationFile)

		if err != nil {
			flagset.Usage()
			panic(fmt.Sprintf(fmt.Sprintf("Error: Failed to open onfiguration file '%s'. Error was %s!\n", configurationFile, err.Error())))
		}

		defer file.Close()
		decoder := json.NewDecoder(file)
		cfgFile := Configuration{}
		err = decoder.Decode(&cfgFile)
		if err != nil {
			flagset.Usage()
			panic(fmt.Sprintf(fmt.Sprintf("Error: Failed to parse onfiguration file '%s'. Error was %s!\n", configurationFile, err.Error())))
		}

		if cfg.Port < 0 {
			cfg.Port = cfgFile.Port
		}

		if len(cfg.Host) == 0 {
			cfg.Host = cfgFile.Host
		}

		if len(cfg.Name) == 0 {
			cfg.Name = cfgFile.Name
		}

		if len(cfg.Hostname) == 0 {
			cfg.Hostname = cfgFile.Hostname
		}

		if len(cfg.Version) == 0 {
			cfg.Version = cfgFile.Version
		}

		if len(cfg.Namespace) == 0 {
			cfg.Namespace = cfgFile.Namespace
		}

		if len(cfg.DBName) == 0 {
			cfg.DBName = cfgFile.DBName
		}

		if len(cfg.CertChainFile) == 0 {
			cfg.CertChainFile = cfgFile.CertChainFile
		}

		if len(cfg.KeyFile) == 0 {
			cfg.KeyFile = cfgFile.KeyFile
		}

		if len(cfg.CAFile) == 0 {
			cfg.CAFile = cfgFile.CAFile
		}

		// if cfg.MaxTcpConnections < 0 {
		// 	cfg.MaxTcpConnections = cfgFile.MaxConnections
		// }

		if cfg.MaxConnections < 0 {
			cfg.MaxConnections = cfgFile.MaxConnections
		}

		if cfg.DelayReply < 0 {
			cfg.DelayReply = cfgFile.DelayReply
		}

		if cfg.ClientTimeout < 0 {
			cfg.ClientTimeout = cfgFile.ClientTimeout
		}

		if len(cfg.Logfile) == 0 {
			cfg.Logfile = cfgFile.Logfile
		}

		cfg.NoMetrics = cfgFile.NoMetrics
	}

	if len(cfg.Name) == 0 {
		cfg.Name = os.Getenv("HOSTNAME")
	}

	if len(cfg.Name) == 0 {
		var seededRand *rand.Rand = rand.New(
			rand.NewSource(time.Now().UnixNano()))
		charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
		b := make([]byte, 10)
		for i := range b {
			b[i] = charset[seededRand.Intn(len(charset))]
		}
		cfg.Name = string(b)
	}

	if len(cfg.Hostname) == 0 {
		cfg.Hostname = os.Getenv("HOSTNAME")
	}

	if len(cfg.Hostname) == 0 {
		cfg.Hostname = cfg.Name
	}

	if len(cfg.Host) == 0 {
		cfg.Host = DefaultHost
	}

	if cfg.Port <= 0 {
		cfg.Port = DefaultPort
	}
}
