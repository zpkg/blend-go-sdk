package template

import (
	"crypto/hmac"
	"crypto/md5"
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
	"math"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	texttemplate "text/template"
	"time"

	jmespath "github.com/blend/go-jmespath"
	"github.com/blend/go-sdk/util"

	"github.com/blend/go-sdk/semver"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/yaml"
)

// Funcs are common view func helpers.
var Funcs ViewFuncs

// ViewFuncs is the type stub for view functions.
type ViewFuncs struct{}

// FileExists returns if the file at a given path exists.
func (vf ViewFuncs) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// File reads the contents of a file path as a string.
func (vf ViewFuncs) File(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	return string(contents), err
}

// ExpandEnv replaces ${var} or $var in the string according
// to the current environment variables. References to undefined
// variables are replaced with the empty string.
func (vf ViewFuncs) ExpandEnv(corpus string) string {
	return os.ExpandEnv(corpus)
}

// ToString attempts to return a string representation of a value.
func (vf ViewFuncs) ToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// ToBytes attempts to return a bytes representation of a value.
func (vf ViewFuncs) ToBytes(v interface{}) []byte {
	return []byte(fmt.Sprintf("%v", v))
}

// ToInt parses a value as an integer.
func (vf ViewFuncs) ToInt(v interface{}) (int, error) {
	return strconv.Atoi(fmt.Sprintf("%v", v))
}

// ToInt64 parses a value as an int64.
func (vf ViewFuncs) ToInt64(v interface{}) (int64, error) {
	return strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
}

// ToFloat64 parses a value as a float64.
func (vf ViewFuncs) ToFloat64(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

// Unix returns the unix format for a timestamp.
func (vf ViewFuncs) Unix(t time.Time) string {
	return fmt.Sprintf("%d", t.Unix())
}

// RFC3339 returns the RFC3339 format for a timestamp.
func (vf ViewFuncs) RFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// TimeShort returns the short format for a timestamp.
// The format string is "1/02/2006 3:04:05 PM".
func (vf ViewFuncs) TimeShort(t time.Time) string {
	return t.Format("1/02/2006 3:04:05 PM")
}

// DateShort returns the short date for a timestamp.
// The format string is "1/02/2006"
func (vf ViewFuncs) DateShort(t time.Time) string {
	return t.Format("1/02/2006")
}

// TimeMedium returns the medium format for a timestamp.
// The format string is "1/02/2006 3:04:05 PM".
func (vf ViewFuncs) TimeMedium(t time.Time) string {
	return t.Format("Jan 02, 2006 3:04:05 PM")
}

// TimeKitchen returns the kitchen format for a timestamp.
// The format string is "3:04PM".
func (vf ViewFuncs) TimeKitchen(t time.Time) string {
	return t.Format(time.Kitchen)
}

// DateMonthDay returns the month dat format for a timestamp.
// The format string is "1/2".
func (vf ViewFuncs) DateMonthDay(t time.Time) string {
	return t.Format("1/2")
}

// TimeIn returns the time in a given location by string.
// If the location is invalid, this will error.
func (vf ViewFuncs) TimeIn(loc string, t time.Time) (time.Time, error) {
	location, err := time.LoadLocation(loc)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(location), err
}

// ParseTime parses a time string with a given format.
func (vf ViewFuncs) ParseTime(format, v string) (time.Time, error) {
	return time.Parse(format, v)
}

// ParseUnix returns a timestamp from a unix format.
func (vf ViewFuncs) ParseUnix(v int64) time.Time {
	return time.Unix(v, 0)
}

// Year returns the year component of a timestamp.
func (vf ViewFuncs) Year(t time.Time) int {
	return t.Year()
}

// Month returns the month component of a timestamp.
func (vf ViewFuncs) Month(t time.Time) int {
	return int(t.Month())
}

// Day returns the day component of a timestamp.
func (vf ViewFuncs) Day(t time.Time) int {
	return t.Day()
}

// Hour returns the hour component of a timestamp.
func (vf ViewFuncs) Hour(t time.Time) int {
	return t.Hour()
}

// Minute returns the minute component of a timestamp.
func (vf ViewFuncs) Minute(t time.Time) int {
	return t.Minute()
}

// Second returns the seconds component of a timestamp.
func (vf ViewFuncs) Second(t time.Time) int {
	return t.Second()
}

// Millisecond returns the millisecond component of a timestamp.
func (vf ViewFuncs) Millisecond(t time.Time) int {
	return int(time.Duration(t.Nanosecond()) / time.Millisecond)
}

// Since returns the duration since a given timestamp.
// It is relative, meaning the value returned can be negative.
func (vf ViewFuncs) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// SinceUTC returns the duration since a given timestamp in UTC.
// It is relative, meaning the value returned can be negative.
func (vf ViewFuncs) SinceUTC(t time.Time) time.Duration {
	return time.Now().UTC().Sub(t.UTC())
}

// ParseBool attempts to parse a value as a bool.
// "truthy" values include "true", "1", "yes".
// "falsey" values include "false", "0", "no".
func (vf ViewFuncs) ParseBool(raw interface{}) (bool, error) {
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

// Round returns the value rounded to a given set of places.
// It uses midpoint rounding.
func (vf ViewFuncs) Round(d float64, places int) float64 {
	return util.Math.Round(d, places)
}

// Ceil returns the value rounded up to the nearest integer.
func (vf ViewFuncs) Ceil(d float64) float64 {
	return math.Ceil(d)
}

// Floor returns the value rounded down to zero.
func (vf ViewFuncs) Floor(d float64) float64 {
	return math.Floor(d)
}

// FormatMoney returns a float as a formatted string rounded to two decimal places.
func (vf ViewFuncs) FormatMoney(d float64) string {
	return fmt.Sprintf("$%0.2f", util.Math.Round(d, 2))
}

// FormatPct formats a float as a percentage (it is multiplied by 100,
// then suffixed with '%')
func (vf ViewFuncs) FormatPct(d float64) string {
	return fmt.Sprintf("%0.2f%%", d*100)
}

// Base64 encodes data as a string as a base6 string.
func (vf ViewFuncs) Base64(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

//Base64Decode decodes a base 64 string.
func (vf ViewFuncs) Base64Decode(v string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// UUIDv4 generates a uuid v4.
func (vf ViewFuncs) UUIDv4() uuid.UUID {
	return uuid.V4()
}

// ToUpper returns a string case shifted to upper case.
func (vf ViewFuncs) ToUpper(v string) string {
	return strings.ToUpper(v)
}

// ToLower returns a string case shifted to lower case.
func (vf ViewFuncs) ToLower(v string) string {
	return strings.ToLower(v)
}

// ToTitle returns a title cased string.
func (vf ViewFuncs) ToTitle(v string) string {
	return strings.ToTitle(v)
}

// TrimSpace trims whitespace from the beginning and end of a string.
func (vf ViewFuncs) TrimSpace(v string) string {
	return strings.TrimSpace(v)
}

// Prefix appends a given string to a prefix.
func (vf ViewFuncs) Prefix(pref, v string) string {
	return pref + v
}

// Suffix appends a given prefix to a string.
func (vf ViewFuncs) Suffix(suf, v string) string {
	return v + suf
}

// Split splits a string by a separator.
func (vf ViewFuncs) Split(sep, v string) []string {
	return strings.Split(v, sep)
}

// SplitN splits a string by a separator a given number of times.
func (vf ViewFuncs) SplitN(sep, v string, n int) []string {
	return strings.SplitN(v, sep, n)
}

// RandomLetters returns a string of random letters.
func (vf ViewFuncs) RandomLetters(length int) string {
	return util.String.RandomLetters(length)
}

// RandomLettersWithNumbers returns a string of random letters.
func (vf ViewFuncs) RandomLettersWithNumbers(count int) string {
	return util.String.RandomStringWithNumbers(count)
}

// RandomLettersWithNumbersAndSymbols returns a string of random letters.
func (vf ViewFuncs) RandomLettersWithNumbersAndSymbols(count int) string {
	return util.String.RandomStringWithNumbersAndSymbols(count)
}

//
// array functions
//

// Slice returns a subrange of a collection.
func (vf ViewFuncs) Slice(from, to int, collection interface{}) (interface{}, error) {
	value := reflect.ValueOf(collection)

	if value.Type().Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}

	return value.Slice(from, to).Interface(), nil
}

// First returns the first element of a collection.
func (vf ViewFuncs) First(collection interface{}) (interface{}, error) {
	value := reflect.ValueOf(collection)
	if value.Type().Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}
	if value.Len() == 0 {
		return nil, nil
	}
	return value.Index(0).Interface(), nil
}

// Index returns an element at a given index.
func (vf ViewFuncs) Index(index int, collection interface{}) (interface{}, error) {
	value := reflect.ValueOf(collection)
	if value.Type().Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}
	if value.Len() == 0 {
		return nil, nil
	}
	return value.Index(index).Interface(), nil
}

// Last returns the last element of a collection.
func (vf ViewFuncs) Last(collection interface{}) (interface{}, error) {
	value := reflect.ValueOf(collection)
	if value.Type().Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}
	if value.Len() == 0 {
		return nil, nil
	}
	return value.Index(value.Len() - 1).Interface(), nil
}

// Join creates a string joined with a given separator.
func (vf ViewFuncs) Join(sep string, collection interface{}) (string, error) {
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

// HasSuffix returns if a string has a given suffix.
func (vf ViewFuncs) HasSuffix(suffix, v string) bool {
	return strings.HasSuffix(v, suffix)
}

// HasPrefix returns if a string has a given prefix.
func (vf ViewFuncs) HasPrefix(prefix, v string) bool {
	return strings.HasPrefix(v, prefix)
}

// Contains returns if a string contains a given substring.
func (vf ViewFuncs) Contains(substr, v string) bool {
	return strings.Contains(v, substr)
}

// Matches returns if a string matches a given regular expression.
func (vf ViewFuncs) Matches(expr, v string) (bool, error) {
	return regexp.MatchString(expr, v)
}

// ParseURL parses a url.
func (vf ViewFuncs) ParseURL(v string) (*url.URL, error) {
	return url.Parse(v)
}

// URLScheme returns the scheme of a url.
func (vf ViewFuncs) URLScheme(v *url.URL) string {
	return v.Scheme
}

// URLHost returns the host of a url.
func (vf ViewFuncs) URLHost(v *url.URL) string {
	return v.Host
}

// URLPort returns the url port or a default.
func (vf ViewFuncs) URLPort(v *url.URL) string {
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

// URLPath returns the url path.
func (vf ViewFuncs) URLPath(v *url.URL) string {
	return v.Path
}

// URLRawQuery returns the url raw query.
func (vf ViewFuncs) URLRawQuery(v *url.URL) string {
	return v.RawQuery
}

// URLQuery returns a url query param.
func (vf ViewFuncs) URLQuery(name string, v *url.URL) string {
	return v.Query().Get(name)
}

// MD5 returns the md5 sum of a string.
func (vf ViewFuncs) MD5(v string) string {
	h := md5.New()
	io.WriteString(h, v)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA1 returns the sha1 sum of a string.
func (vf ViewFuncs) SHA1(v string) string {
	h := sha1.New()
	io.WriteString(h, v)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 returns the sha256 sum of a string.
func (vf ViewFuncs) SHA256(v string) string {
	h := sha256.New()
	io.WriteString(h, v)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA512 returns the sha512 sum of a string.
func (vf ViewFuncs) SHA512(v string) string {
	h := sha512.New()
	io.WriteString(h, v)
	return hex.EncodeToString(h.Sum(nil))
}

// HMAC512 returns the hmac signed sha 512 sum of a string.
func (vf ViewFuncs) HMAC512(key, v string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha512.New, keyBytes)
	io.WriteString(h, v)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// ParseSemver parses a semantic version string.
func (vf ViewFuncs) ParseSemver(v string) (*semver.Version, error) {
	return semver.NewVersion(v)
}

// SemverMajor returns the major component of a semver.
func (vf ViewFuncs) SemverMajor(v *semver.Version) int {
	return int(v.Major())
}

// SemverBumpMajor returns a semver with an incremented major version.
func (vf ViewFuncs) SemverBumpMajor(v *semver.Version) *semver.Version {
	v.BumpMajor()
	return v
}

// SemverMinor returns the minor component of a semver.
func (vf ViewFuncs) SemverMinor(v *semver.Version) int {
	return int(v.Minor())
}

// SemverBumpMinor returns a semver with an incremented minor version.
func (vf ViewFuncs) SemverBumpMinor(v *semver.Version) *semver.Version {
	v.BumpMinor()
	return v
}

// SemverPatch returns the patch component of a semver.
func (vf ViewFuncs) SemverPatch(v *semver.Version) int {
	return int(v.Patch())
}

// SemverBumpPatch returns a semver with an incremented patch version.
func (vf ViewFuncs) SemverBumpPatch(v *semver.Version) *semver.Version {
	v.BumpPatch()
	return v
}

// IndentTabs indents a string with a given number of tabs.
func (vf ViewFuncs) IndentTabs(tabCount int, v interface{}) string {
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

// IndentSpaces indents a string by a given set of spaces.
func (vf ViewFuncs) IndentSpaces(spaceCount int, v interface{}) string {
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

// GenerateOrdinalNames generates ordinal names by passing the index to a given formatter.
// The formatter should be in Sprintf format (i.e. using a '%d' token for where the index should go).
/*
Example:
    {{ generate_ordinal_names "worker-%d" 3 }} // [ worker-0, worker-1, worker-2 ]
*/
func (vf ViewFuncs) GenerateOrdinalNames(format string, replicas int) []string {
	output := make([]string, replicas)
	for index := 0; index < replicas; index++ {
		output[index] = fmt.Sprintf(format, index)
	}
	return output
}

// GenerateKey generates a key of a given size base 64 encoded.
func (vf ViewFuncs) GenerateKey(keySize int) string {
	key := make([]byte, keySize)
	io.ReadFull(rand.Reader, key)
	return base64.StdEncoding.EncodeToString(key)
}

// YAMLEncode returns an object encoded as yaml.
func (vf ViewFuncs) YAMLEncode(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	return string(data), err
}

// JSONEncode returns an object encoded as json.
func (vf ViewFuncs) JSONEncode(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	return string(data), err
}

// ParseYAML decodes a corups as yaml.
func (vf ViewFuncs) ParseYAML(v string) (interface{}, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(v), &data)
	return data, err
}

// ParseJSON returns an object encoded as json.
func (vf ViewFuncs) ParseJSON(v string) (interface{}, error) {
	var data interface{}
	err := json.Unmarshal([]byte(v), &data)
	return data, err
}

// YAMLPath parses a yaml string and returns the value at a given path.
// It uses JMESPath syntax.
func (vf ViewFuncs) YAMLPath(path string, v string) (interface{}, error) {
	var data interface{}
	if err := yaml.Unmarshal([]byte(v), &data); err != nil {
		return nil, err
	}
	return jmespath.Search(path, data)
}

// JSONPath parses a json string and returns the value at a given path.
// It uses JMESPath syntax.
func (vf ViewFuncs) JSONPath(path string, v string) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(v), &data); err != nil {
		return nil, err
	}
	return jmespath.Search(path, data)
}

// JMESPath applies the search query to the object.
// More information can be found here: http://jmespath.org/specification.html
func (vf ViewFuncs) JMESPath(path string, v interface{}) (interface{}, error) {
	return jmespath.Search(path, v)
}

// FuncMap returns the name => func mapping.
func (vf ViewFuncs) FuncMap() texttemplate.FuncMap {
	return texttemplate.FuncMap{
		"file_exists":                 vf.FileExists,
		"file":                        vf.File,
		"expand_env":                  vf.ExpandEnv,
		"to_string":                   vf.ToString,
		"to_bytes":                    vf.ToBytes,
		"to_int":                      vf.ToInt,
		"to_int64":                    vf.ToInt64,
		"to_float64":                  vf.ToFloat64,
		"date_short":                  vf.DateShort,
		"date_month_day":              vf.DateMonthDay,
		"unix":                        vf.Unix,
		"rfc3339":                     vf.RFC3339,
		"time_short":                  vf.TimeShort,
		"time_medium":                 vf.TimeMedium,
		"time_kitchen":                vf.TimeKitchen,
		"time_in":                     vf.TimeIn,
		"parse_time":                  vf.ParseTime,
		"parse_unix":                  vf.ParseUnix,
		"since":                       vf.Since,
		"since_utc":                   vf.SinceUTC,
		"year":                        vf.Year,
		"month":                       vf.Month,
		"day":                         vf.Day,
		"hour":                        vf.Hour,
		"minute":                      vf.Minute,
		"second":                      vf.Second,
		"millisecond":                 vf.Millisecond,
		"parse_bool":                  vf.ParseBool,
		"format_money":                vf.FormatMoney,
		"format_pct":                  vf.FormatPct,
		"round":                       vf.Round,
		"ceil":                        vf.Ceil,
		"floor":                       vf.Floor,
		"base64":                      vf.Base64,
		"base64decode":                vf.Base64Decode,
		"uuidv4":                      vf.UUIDv4,
		"to_upper":                    vf.ToUpper,
		"to_lower":                    vf.ToLower,
		"to_title":                    vf.ToTitle,
		"random_letters":              vf.RandomLetters,
		"random_letters_with_numbers": vf.RandomLettersWithNumbers,
		"trim_space":                  vf.TrimSpace,
		"prefix":                      vf.Prefix,
		"suffix":                      vf.Suffix,
		"split":                       vf.Split,
		"splitn":                      vf.SplitN,
		"slice":                       vf.Slice,
		"first":                       vf.First,
		"index":                       vf.Index,
		"last":                        vf.Last,
		"join":                        vf.Join,
		"has_suffix":                  vf.HasSuffix,
		"has_prefix":                  vf.HasPrefix,
		"contains":                    vf.Contains,
		"matches":                     vf.Matches,
		"parse_url":                   vf.ParseURL,
		"url_scheme":                  vf.URLScheme,
		"url_host":                    vf.URLHost,
		"url_port":                    vf.URLPort,
		"url_path":                    vf.URLPath,
		"url_rawquery":                vf.URLRawQuery,
		"url_query":                   vf.URLQuery,
		"md5":                         vf.MD5,
		"sha1":                        vf.SHA1,
		"sha256":                      vf.SHA256,
		"sha512":                      vf.SHA512,
		"hmac":                        vf.HMAC512,
		"parse_semver":                vf.ParseSemver,
		"semver_major":                vf.SemverMajor,
		"semver_bump_major":           vf.SemverBumpMajor,
		"semver_minor":                vf.SemverMinor,
		"semver_bump_minor":           vf.SemverBumpMinor,
		"semver_patch":                vf.SemverPatch,
		"semver_bump_patch":           vf.SemverBumpPatch,
		"generate_ordinal_names":      vf.GenerateOrdinalNames,
		"generate_password":           vf.RandomLettersWithNumbersAndSymbols,
		"generate_key":                vf.GenerateKey,
		"to_json":                     vf.JSONEncode,
		"to_yaml":                     vf.YAMLEncode,
		"parse_json":                  vf.ParseJSON,
		"parse_yaml":                  vf.ParseYAML,
		"json_path":                   vf.JSONPath,
		"yaml_path":                   vf.YAMLPath,
		"indent_tabs":                 vf.IndentTabs,
		"indent_spaces":               vf.IndentSpaces,
	}
}
