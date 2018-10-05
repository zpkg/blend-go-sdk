package template

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	texttemplate "text/template"
	"time"

	"github.com/blend/go-sdk/semver"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/yaml"
)

// ViewFuncs are common view func helpers.
var ViewFuncs viewFuncs

type viewFuncs struct{}

func (vf viewFuncs) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (vf viewFuncs) File(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	return string(contents), err
}

func (vf viewFuncs) ExpandEnv(corpus string) string {
	return os.ExpandEnv(corpus)
}

func (vf viewFuncs) ToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func (vf viewFuncs) Unix(t time.Time) string {
	return fmt.Sprintf("%d", t.Unix())
}

func (vf viewFuncs) RFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

func (vf viewFuncs) Short(t time.Time) string {
	return t.Format("1/02/2006 3:04:05 PM")
}

func (vf viewFuncs) ShortDate(t time.Time) string {
	return t.Format("1/02/2006")
}

func (vf viewFuncs) Medium(t time.Time) string {
	return t.Format("Jan 02, 2006 3:04:05 PM")
}

func (vf viewFuncs) Kitchen(t time.Time) string {
	return t.Format(time.Kitchen)
}

func (vf viewFuncs) MonthDay(t time.Time) string {
	return t.Format("1/2")
}

func (vf viewFuncs) InLoc(loc string, t time.Time) (time.Time, error) {
	location, err := time.LoadLocation(loc)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(location), err
}

func (vf viewFuncs) Time(format, v string) (time.Time, error) {
	return time.Parse(format, v)
}

func (vf viewFuncs) TimeUnix(v int64) time.Time {
	return time.Unix(v, 0)
}

func (vf viewFuncs) Year(t time.Time) int {
	return t.Year()
}

func (vf viewFuncs) Month(t time.Time) int {
	return int(t.Month())
}

func (vf viewFuncs) Day(t time.Time) int {
	return t.Day()
}

func (vf viewFuncs) Hour(t time.Time) int {
	return t.Hour()
}

func (vf viewFuncs) Minute(t time.Time) int {
	return t.Minute()
}

func (vf viewFuncs) Second(t time.Time) int {
	return t.Second()
}

func (vf viewFuncs) Millisecond(t time.Time) int {
	return int(time.Duration(t.Nanosecond()) / time.Millisecond)
}

func (vf viewFuncs) Bool(raw interface{}) (bool, error) {
	v := fmt.Sprintf("%v", raw)
	if len(v) == 0 {
		return false, nil
	}
	switch strings.ToLower(v) {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value `%s`", v)
	}
}

func (vf viewFuncs) Int(v interface{}) (int, error) {
	return strconv.Atoi(fmt.Sprintf("%v", v))
}

func (vf viewFuncs) Int64(v interface{}) (int64, error) {
	return strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
}

func (vf viewFuncs) Float64(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

func (vf viewFuncs) Money(d float64) string {
	return fmt.Sprintf("$%0.2f", d)
}

func (vf viewFuncs) Pct(d float64) string {
	return fmt.Sprintf("%0.2f%%", d*100)
}

func (vf viewFuncs) Base64(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func (vf viewFuncs) Base64Decode(v string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (vf viewFuncs) CreateKey(keySize int) string {
	key := make([]byte, keySize)
	io.ReadFull(rand.Reader, key)
	return base64.StdEncoding.EncodeToString(key)
}

func (vf viewFuncs) UUIDv4() string {
	return uuid.V4().String()
}

func (vf viewFuncs) ToUpper(v string) string {
	return strings.ToUpper(v)
}

func (vf viewFuncs) ToLower(v string) string {
	return strings.ToLower(v)
}

func (vf viewFuncs) ToTitle(v string) string {
	return strings.ToTitle(v)
}

func (vf viewFuncs) TrimSpace(v string) string {
	return strings.TrimSpace(v)
}

func (vf viewFuncs) Prefix(pref, v string) string {
	return pref + v
}

func (vf viewFuncs) Suffix(suf, v string) string {
	return v + suf
}

func (vf viewFuncs) Split(sep, v string) []string {
	return strings.Split(v, sep)
}

func (vf viewFuncs) SplitN(sep, v string, n int) []string {
	return strings.SplitN(v, sep, n)
}

func (vf viewFuncs) Slice(from, to int, collection interface{}) (interface{}, error) {
	value := reflect.ValueOf(collection)

	if value.Type().Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}

	return value.Slice(from, to).Interface(), nil
}

func (vf viewFuncs) First(collection interface{}) (interface{}, error) {
	value := reflect.ValueOf(collection)
	if value.Type().Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}
	if value.Len() == 0 {
		return nil, nil
	}
	return value.Index(0).Interface(), nil
}

func (vf viewFuncs) Index(index int, collection interface{}) (interface{}, error) {
	value := reflect.ValueOf(collection)
	if value.Type().Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}
	if value.Len() == 0 {
		return nil, nil
	}
	return value.Index(index).Interface(), nil
}

func (vf viewFuncs) Last(collection interface{}) (interface{}, error) {
	value := reflect.ValueOf(collection)
	if value.Type().Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}
	if value.Len() == 0 {
		return nil, nil
	}
	return value.Index(value.Len() - 1).Interface(), nil
}

func (vf viewFuncs) Join(sep string, collection interface{}) (string, error) {
	value := reflect.ValueOf(collection)
	if value.Type().Kind() != reflect.Slice {
		return "", fmt.Errorf("input must be a slice")
	}
	if value.Len() == 0 {
		return "", nil
	}
	values := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		values[i] = fmt.Sprintf("%v", value.Index(i).Interface())
	}
	return strings.Join(values, sep), nil
}

func (vf viewFuncs) HasSuffix(suffix, v string) bool {
	return strings.HasSuffix(v, suffix)
}

func (vf viewFuncs) HasPrefix(prefix, v string) bool {
	return strings.HasPrefix(v, prefix)
}

func (vf viewFuncs) Contains(substr, v string) bool {
	return strings.Contains(v, substr)
}

func (vf viewFuncs) Matches(expr, v string) (bool, error) {
	return regexp.MatchString(expr, v)
}

func (vf viewFuncs) ParseURL(v string) (*url.URL, error) {
	return url.Parse(v)
}

func (vf viewFuncs) URLScheme(v *url.URL) string {
	return v.Scheme
}

func (vf viewFuncs) URLHost(v *url.URL) string {
	return v.Host
}

func (vf viewFuncs) URLPort(v *url.URL) string {
	portValue := v.Port()
	if len(portValue) > 0 {
		return portValue
	}
	switch strings.ToLower(v.Scheme) {
	case "http":
		return "80"
	case "https":
		return "443"
	case "ssh":
		return "22"
	case "ftp":
		return "21"
	case "sftp":
		return "22"
	}
	return ""
}

func (vf viewFuncs) URLPath(v *url.URL) string {
	return v.Path
}

func (vf viewFuncs) URLRawQuery(v *url.URL) string {
	return v.RawQuery
}

func (vf viewFuncs) URLQuery(name string, v *url.URL) string {
	return v.Query().Get(name)
}

func (vf viewFuncs) SHA1(v string) string {
	h := sha1.New()
	io.WriteString(h, v)
	return hex.EncodeToString(h.Sum(nil))
}

func (vf viewFuncs) SHA256(v string) string {
	h := sha256.New()
	io.WriteString(h, v)
	return hex.EncodeToString(h.Sum(nil))
}

func (vf viewFuncs) SHA512(v string) string {
	h := sha512.New()
	io.WriteString(h, v)
	return hex.EncodeToString(h.Sum(nil))
}

func (vf viewFuncs) ParseSemver(v string) (*semver.Version, error) {
	return semver.NewVersion(v)
}

func (vf viewFuncs) SemverMajor(v *semver.Version) int {
	return int(v.Major())
}

func (vf viewFuncs) SemverBumpMajor(v *semver.Version) *semver.Version {
	v.BumpMajor()
	return v
}

func (vf viewFuncs) SemverMinor(v *semver.Version) int {
	return int(v.Minor())
}

func (vf viewFuncs) SemverBumpMinor(v *semver.Version) *semver.Version {
	v.BumpMinor()
	return v
}

func (vf viewFuncs) SemverPatch(v *semver.Version) int {
	return int(v.Patch())
}

func (vf viewFuncs) SemverBumpPatch(v *semver.Version) *semver.Version {
	v.BumpPatch()
	return v
}

func (vf viewFuncs) YAML(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	return string(data), err
}

func (vf viewFuncs) JSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	return string(data), err
}

func (vf viewFuncs) IndentTabs(tabCount int, v interface{}) string {
	lines := strings.Split(fmt.Sprintf("%v", v), "\n")
	outputLines := make([]string, len(lines))

	var tabs string
	for i := 0; i < tabCount; i++ {
		tabs = tabs + "\t"
	}

	for i := 0; i < len(lines); i++ {
		outputLines[i] = tabs + lines[i]
	}
	return strings.Join(outputLines, "\n")
}

func (vf viewFuncs) IndentSpaces(spaceCount int, v interface{}) string {
	lines := strings.Split(fmt.Sprintf("%v", v), "\n")
	outputLines := make([]string, len(lines))

	var spaces string
	for i := 0; i < spaceCount; i++ {
		spaces = spaces + " "
	}

	for i := 0; i < len(lines); i++ {
		outputLines[i] = spaces + lines[i]
	}
	return strings.Join(outputLines, "\n")
}

func (vf viewFuncs) FuncMap() texttemplate.FuncMap {
	return texttemplate.FuncMap{
		"file_exists":       ViewFuncs.FileExists,
		"file":              ViewFuncs.File,
		"expand_env":        ViewFuncs.ExpandEnv,
		"to_string":         ViewFuncs.ToString,
		"unix":              ViewFuncs.Unix,
		"rfc3339":           ViewFuncs.RFC3339,
		"short":             ViewFuncs.Short,
		"short_date":        ViewFuncs.ShortDate,
		"medium":            ViewFuncs.Medium,
		"kitchen":           ViewFuncs.Kitchen,
		"month_day":         ViewFuncs.MonthDay,
		"in_loc":            ViewFuncs.InLoc,
		"time":              ViewFuncs.Time,
		"time_unix":         ViewFuncs.TimeUnix,
		"year":              ViewFuncs.Year,
		"month":             ViewFuncs.Month,
		"day":               ViewFuncs.Day,
		"hour":              ViewFuncs.Hour,
		"minute":            ViewFuncs.Minute,
		"second":            ViewFuncs.Second,
		"millisecond":       ViewFuncs.Millisecond,
		"bool":              ViewFuncs.Bool,
		"int":               ViewFuncs.Int,
		"int64":             ViewFuncs.Int64,
		"float64":           ViewFuncs.Float64,
		"money":             ViewFuncs.Money,
		"pct":               ViewFuncs.Pct,
		"base64":            ViewFuncs.Base64,
		"base64decode":      ViewFuncs.Base64Decode,
		"createKey":         ViewFuncs.CreateKey,
		"uuidv4":            ViewFuncs.UUIDv4,
		"to_upper":          ViewFuncs.ToUpper,
		"to_lower":          ViewFuncs.ToLower,
		"to_title":          ViewFuncs.ToTitle,
		"trim_space":        ViewFuncs.TrimSpace,
		"prefix":            ViewFuncs.Prefix,
		"suffix":            ViewFuncs.Suffix,
		"split":             ViewFuncs.Split,
		"splitn":            ViewFuncs.SplitN,
		"slice":             ViewFuncs.Slice,
		"first":             ViewFuncs.First,
		"index":             ViewFuncs.Index,
		"last":              ViewFuncs.Last,
		"join":              ViewFuncs.Join,
		"has_suffix":        ViewFuncs.HasSuffix,
		"has_prefix":        ViewFuncs.HasPrefix,
		"contains":          ViewFuncs.Contains,
		"matches":           ViewFuncs.Matches,
		"parse_url":         ViewFuncs.ParseURL,
		"url_scheme":        ViewFuncs.URLScheme,
		"url_host":          ViewFuncs.URLHost,
		"url_port":          ViewFuncs.URLPort,
		"url_path":          ViewFuncs.URLPath,
		"url_rawquery":      ViewFuncs.URLRawQuery,
		"url_query":         ViewFuncs.URLQuery,
		"sha1":              ViewFuncs.SHA1,
		"sha256":            ViewFuncs.SHA256,
		"sha512":            ViewFuncs.SHA512,
		"parse_semver":      ViewFuncs.ParseSemver,
		"semver_major":      ViewFuncs.SemverMajor,
		"semver_bump_major": ViewFuncs.SemverBumpMajor,
		"semver_minor":      ViewFuncs.SemverMinor,
		"semver_bump_minor": ViewFuncs.SemverBumpMinor,
		"semver_patch":      ViewFuncs.SemverPatch,
		"semver_bump_patch": ViewFuncs.SemverBumpPatch,
		"yaml":              ViewFuncs.YAML,
		"json":              ViewFuncs.JSON,
		"indent_tabs":       ViewFuncs.IndentTabs,
		"indent_spaces":     ViewFuncs.IndentSpaces,
	}
}
