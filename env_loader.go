package cfg

import (
	"log/slog"
	"os"
	"strings"
)

type UnmarshalFn func(content []byte) (map[string]any, error)

// SetFileSystem define a instância do FileSystem que será usado para carregamento
func (c *Env) SetFileSystem(fs FileSystem) {
	c.fs = fs
}

// SetFilePaths define o caminho dos arquivos de configuração.
func (c *Env) SetFilePaths(filePaths ...string) {
	c.filePaths = filePaths
}

// SetFileExt define o processador para essa extensão de arquivo. Usado para suportar .yaml, .toml, .xml
func (c *Env) SetFileExt(ext string, fn UnmarshalFn) {
	if fn == nil {
		delete(c.fileExts, ext)
	} else {
		c.fileExts[ext] = fn
	}
}

// SetProfileKey define a key que identifica os arquivos de perfil de configuração.
func (c *Env) SetProfileKey(profileKey string) {
	c.profileKey = profileKey
}

// Load initialize default settings. A variable is obtained respecting the order below.
//
// 1) command line arguments (starting with "--", e.g. --server.port=9000)
// 2) DotEnv file variables ".env"
// 3) Operating system variables
// 4) Profile specific configuration (config-{dev|prod|test}.json)
// 5) global config (config.json)
// 6) Default config (cfg.New(DefaultConfig))
func (c *Env) Load() error {

	// 5) global config (config.json)
	if err := c.LoadFiles(); err != nil {
		return err
	}

	profiles := c.String(c.profileKey)

	// settings with priority over profiles
	h := New()
	h.fs = c.fs
	h.fileExts = c.fileExts
	h.filePaths = c.filePaths
	h.profileKey = c.profileKey

	// 3) Operating system variables
	h.LoadOsEnv()

	// 2) DotEnv file variables ".env"
	if err := h.LoadDotEnv(); err != nil {
		return err
	}

	// 1) command line arguments (starting with "--", e.g. --server.port=9000)
	h.LoadOsArgs(os.Args[1:])

	// 4) Profile specific configuration (config-{dev|prod|test}.json)
	newProfile := h.String(c.profileKey)
	if newProfile != "" && newProfile != profiles {
		c.LoadObject(O{c.profileKey: newProfile})
	}
	if err := c.LoadProfiles(); err != nil {
		return err
	}
	c.Merge(h)
	return nil
}

// LoadOsArgs will convert any command line option arguments (starting with ‘--’, e.g. --server.port=9000) to a
// property and add it to the Env.
//
// Command line properties always take precedence over other property sources.
func (c *Env) LoadOsArgs(args []string) {
	//args := os.Args[1:]

	var environ []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") && strings.IndexByte(arg, '=') > 1 {
			// --server.port=9000
			environ = append(environ, strings.TrimPrefix(arg, "--"))
		}
	}
	if len(environ) > 0 {
		c.LoadEnviron(environ)
	}
}

// LoadOsEnv obtém todas as configurações do ambiente
func (c *Env) LoadOsEnv() {
	c.LoadEnviron(os.Environ())
}

// LoadDotEnv from https://github.com/joho/godotenv
func (c *Env) LoadDotEnv() error {
	if content, err := c.loadFile(".env"); err != nil {
		return err
	} else if content != nil {
		c.LoadEnviron(strings.Split(string(content), "\n"))
	}
	return nil
}

func (c *Env) LoadEnviron(environ []string) {
	config := map[string]any{}
	for _, env := range environ {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if strings.IndexByte(key, '[') >= 0 {
				continue
			}
			if strings.IndexByte(key, '.') < 0 {
				config[key] = value
			} else {
				// @TODO: AQUI, GERAR O CONTEÚDO INTERNO USANDO A ESTRUTURA DA KEY
				var lastKey string
				parent := config
				var current map[string]any
				for _, k := range Segments(key) {
					lastKey = k
					if child, exist := parent[lastKey]; !exist {
						parent[lastKey] = map[string]any{}
					} else if _, isString := child.(string); isString {
						parent[lastKey] = map[string]any{}
					}
					current = parent
					parent = parent[lastKey].(map[string]any)
				}
				current[lastKey] = value
			}
		}
	}
	c.LoadObject(config)
}

// LoadObject obtém as configurações a partir de um mapa em memória
func (c *Env) LoadObject(config O) {
	if config == nil {
		return
	}
	entries := &Entry{}
	parseEntryMap(config, entries)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.root.Merge(entries)
	c.cache = map[string]*cacheEntry{}
}

// LoadFiles processa arquivos de configuração
func (c *Env) LoadFiles() error {
	for _, filepath := range c.filePaths {
		for ext, fn := range c.fileExts {
			if err := c.processFile(filepath+"."+ext, fn); err != nil {
				return err
			}
		}
	}
	return nil
}

// Profiles get active profiles
func (c *Env) Profiles() []string {
	profiles := c.String(c.profileKey)
	return strings.Split(profiles, ",")
}

// LoadProfiles processa arquivos de configuração (config.json)
func (c *Env) LoadProfiles() error {
	profiles := c.String(c.profileKey)
	for _, profile := range strings.Split(profiles, ",") {
		profile = strings.TrimSpace(profile)
		for _, filepath := range c.filePaths {
			for ext, fn := range c.fileExts {
				// load additional resources
				if err := c.processFile(filepath+"-"+profile+"."+ext, fn); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// processFile faz o carregamento de arquivos config.json e suas variantes (config-{profile}.json)
func (c *Env) processFile(filepath string, unmarshal UnmarshalFn) error {
	if content, err := c.loadFile(filepath); err != nil {
		return err
	} else if content == nil {
		return nil
	} else if config, errUnmarshal := unmarshal(content); errUnmarshal != nil {
		slog.Error(
			"[cfg] error processing file.",
			slog.Any("error", errUnmarshal),
			slog.String("filepath", filepath),
		)
		return errUnmarshal
	} else {
		c.LoadObject(config)
	}
	return nil
}

func (c *Env) loadFile(filepath string) ([]byte, error) {
	if c.fs == nil {
		if fs, err := defaultFileSystem(); err != nil {
			slog.Error(
				"[cfg] could not create default FileSystem.",
				slog.Any("error", err),
				slog.String("filepath", filepath),
			)
			return nil, err
		} else {
			c.fs = fs
		}
	}

	if exist, errFsExists := c.fs.Exists(filepath); errFsExists != nil {
		slog.Error(
			"[cfg] file cannot be loaded.",
			slog.Any("error", errFsExists),
			slog.String("filepath", filepath),
		)
		return nil, errFsExists
	} else if exist {
		slog.Info("[cfg] loading file.", slog.String("filepath", filepath))
		if content, errFsRead := c.fs.Read(filepath); errFsRead != nil {
			slog.Error(
				"[cfg] file cannot be loaded.",
				slog.Any("error", errFsRead),
				slog.String("filepath", filepath),
			)
			return nil, errFsRead
		} else {
			return content, nil
		}
	}
	return nil, nil
}
