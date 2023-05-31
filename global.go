package cfg

import "time"

var c = New(nil)

func Get(key string) any                                      { return c.Get(key) }
func Bool(key string) bool                                    { return c.Bool(key) }
func Int(key string, def ...int) int                          { return c.Int(key, def...) }
func Float(key string, def ...float64) float64                { return c.Float(key, def...) }
func String(key string, def ...string) string                 { return c.String(key, def...) }
func Strings(key string, def ...[]string) []string            { return c.Strings(key, def...) }
func Duration(key string, def ...time.Duration) time.Duration { return c.Duration(key, def...) }

func Time(key string, def ...time.Time) time.Time     { return c.Time(key, def...) }
func DateTime(key string, def ...time.Time) time.Time { return c.DateTime(key, def...) }
func DateOnly(key string, def ...time.Time) time.Time { return c.DateOnly(key, def...) }
func TimeOnly(key string, def ...time.Time) time.Time { return c.TimeOnly(key, def...) }
func TimeLayout(key string, layout string, def ...time.Time) time.Time {
	return c.TimeLayout(key, layout, def...)
}

func SetString(key string, value string) { c.Set(key, value) }
func Clone() *Env                        { return c.Clone() }
func Merge(src *Env)                     { c.Merge(src) }

func SetLogger(logger Logger)               { c.SetLogger(logger) }
func SetFileSystem(fs FileSystem)           { c.SetFileSystem(fs) }
func SetFilePaths(filePaths []string)       { c.SetFilePaths(filePaths) }
func SetFileExt(ext string, fn UnmarshalFn) { c.SetFileExt(ext, fn) }
func SetProfileKey(profileKey string)       { c.SetProfileKey(profileKey) }

func Load() error                  { return c.Load() }
func LoadOsArgs(args []string)     { c.LoadOsArgs(args) }
func LoadOsEnv()                   { c.LoadOsEnv() }
func LoadDotEnv() error            { return c.LoadDotEnv() }
func LoadEnviron(environ []string) { c.LoadEnviron(environ) }
func LoadObject(config O)          { c.LoadObject(config) }
func LoadFiles() error             { return c.LoadFiles() }
func LoadProfiles() error          { return c.LoadProfiles() }
