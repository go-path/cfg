# cfg

Simplified configuration for Go projects.

Automatically loads config files, `.env`, environment variables and command line arguments.

Allows you to define profiles, automatically loading files with the defined profile (e.g. `config.yaml`, `config-dev.yaml`, `config-prod.yaml`).


Docs: https://pkg.go.dev/github.com/go-path/cfg

## How to use?

If you have a config file named `config.json` or `confgi.yaml` in the project root directory. Do.

```go
var config = cfg.Global()

func init() {

    if err := config.Load(); err != nil {
		slog.Error(err.Error(), slog.Any("error", err))
		os.Exit(1)
	}

    slog.Info("My db connection:" + config.String("app.db.host"))
}
```

`config.Load()` initialize default settings. A variable is obtained respecting the order below.

1. command line arguments (starting with "--", e.g. `--server.port=9000`)
2. DotEnv file variables `.env`
3. Operating system variables
4. Profile specific configuration (`config-{dev|prod|test}.{json,yaml,yml}`)
5. global config (`config.{json,yaml,yml}`)
6. Default config (cfg.New(DefaultConfig))


> If you want to change the configuration file name, use the `SetFilePaths` method (default is `"config"`)


## Methods

### Get
- config.Get(key string) any
- config.Bool(key string) bool
- config.Int(key string, def ...int) int
- config.Float(key string, def ...float64) float64
- config.String(key string, def ...string) string
- config.Strings(key string, def ...[]string) []string
- config.Duration(key string, def ...time.Duration) time.Duration
- config.Time(key string, def ...time.Time) time.Time
- config.DateTime(key string, def ...time.Time) time.Time
- config.DateOnly(key string, def ...time.Time) time.Time
- config.TimeOnly(key string, def ...time.Time) time.Time
- config.TimeLayout(key string, layout string, def ...time.Time) time.Time
- config.Keys(key string) []string

### Set
- config.Set(key string, value any)
- config.SetString(key string, value string)

### FileSystem
- config.SetFileSystem(fs FileSystem)
- config.SetFilePaths(filePaths ...string)
- config.SetFileExt(ext string, fn UnmarshalFn)
- config.SetProfileKey(profileKey string)

## Loading
- config.Load() error
- config.LoadOsArgs(args []string)
- config.LoadOsEnv()
- config.LoadDotEnv() error
- config.LoadEnviron(environ []string)
- config.LoadObject(config O)
- config.LoadFiles() error
- config.LoadProfiles() error

## Utils
- config.Clone() *Env
- config.Merge(src *Env)