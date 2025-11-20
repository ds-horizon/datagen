package main

import (
	"encoding/json"
	"image"
	"math"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type SVGOpts struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Type   string `json:"type"`
}

type PasswordOptions struct {
	Lower   bool
	Upper   bool
	Numeric bool
	Special bool
	Space   bool
	Length  int
}

type ParagraphOptions struct {
	ParagraphCount int
	SentenceCount  int
	WordCount      int
	Separator      string
}

func Name() string {
	return gofakeit.Name()
}

func Email() string {
	return gofakeit.Email()
}

func Phone() string {
	return gofakeit.Phone()
}

func Boolean() bool {
	return gofakeit.Bool()
}

func Url() string {
	return gofakeit.URL()
}

func Company() string {
	return gofakeit.Company()
}

func HexColor() string {
	return gofakeit.HexColor()
}

func DateBetweenStr(startStr, endStr string) string {
	const layout = "2006-01-02 15:04:05"

	startDate, err := time.Parse(layout, startStr)
	if err != nil {
		startDate = time.Now()
	}

	endDate, err := time.Parse(layout, endStr)
	if err != nil {
		endDate = time.Now()
	}

	if startDate.After(endDate) {
		startDate, endDate = endDate, startDate
	}

	diff := endDate.Unix() - startDate.Unix()
	randomOffset := rand.Int63n(diff + 1)
	randomTime := startDate.Add(time.Duration(randomOffset) * time.Second)

	return randomTime.Format(layout)
}

func Text(len int) string {
	return gofakeit.LetterN(uint(len))
}

func ToJSON(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func IntBetween(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}

func DateBetween(startDate, endDate time.Time) time.Time {
	if startDate.After(endDate) {
		startDate, endDate = endDate, startDate
	}

	diff := endDate.Unix() - startDate.Unix()
	randomOffset := rand.Int63n(diff + 1)
	randomTime := startDate.Add(time.Duration(randomOffset) * time.Second)
	return randomTime
}

func RandomFrom[T any](values ...T) T {
	if len(values) == 0 {
		var zero T
		return zero
	}
	return values[rand.Intn(len(values))]
}

func UUID() string {
	return gofakeit.UUID()
}

func FloatBetween(min, max float64) float64 {
	value := min + rand.Float64()*(max-min)
	return math.Round(value*100) / 100
}

func Float32Between(min, max float32) float32 {
	value := float64(min) + rand.Float64()*(float64(max)-float64(min))
	rounded := math.Round(value*100) / 100
	return float32(rounded)
}

func FirstName() string {
	return gofakeit.FirstName()
}

func LastName() string {
	return gofakeit.LastName()
}

func NamePrefix() string {
	return gofakeit.NamePrefix()
}

func NameSuffix() string {
	return gofakeit.NameSuffix()
}

func MiddleName() string {
	return gofakeit.MiddleName()
}

func Gender() string {
	return gofakeit.Gender()
}

func SSN() string {
	return gofakeit.SSN()
}

func PhoneFormatted() string {
	return gofakeit.PhoneFormatted()
}

func Username() string {
	return gofakeit.Username()
}


func Password(opts PasswordOptions) string {
	return gofakeit.Password(
		opts.Lower,
		opts.Upper,
		opts.Numeric,
		opts.Special,
		opts.Space,
		opts.Length,
	)
}

func AddressInfo() string {
	return gofakeit.Address().Address
}

func Street() string {
	return gofakeit.Street()
}

func StreetName() string {
	return gofakeit.StreetName()
}

func StreetNumber() string {
	return gofakeit.StreetNumber()
}

func StreetPrefix() string {
	return gofakeit.StreetPrefix()
}

func StreetSuffix() string {
	return gofakeit.StreetSuffix()
}

func City() string {
	return gofakeit.City()
}

func State() string {
	return gofakeit.State()
}

func StateAbr() string {
	return gofakeit.StateAbr()
}

func Zip() string {
	return gofakeit.Zip()
}

func Country() string {
	return gofakeit.Country()
}

func CountryAbr() string {
	return gofakeit.CountryAbr()
}

func Latitude() float64 {
	return gofakeit.Latitude()
}

func Longitude() float64 {
	return gofakeit.Longitude()
}

func Date() time.Time {
	return gofakeit.Date()
}

func FutureDate() time.Time {
	return gofakeit.FutureDate()
}

func PastDate() time.Time {
	return gofakeit.PastDate()
}

func NanoSecond() int {
	return gofakeit.NanoSecond()
}

func Second() int {
	return gofakeit.Second()
}

func Minute() int {
	return gofakeit.Minute()
}

func Hour() int {
	return gofakeit.Hour()
}

func Month() int {
	return gofakeit.Month()
}

func MonthString() string {
	return gofakeit.MonthString()
}

func Day() int {
	return gofakeit.Day()
}

func WeekDay() string {
	return gofakeit.WeekDay()
}

func Year() int {
	return gofakeit.Year()
}

func TimeZone() string {
	return gofakeit.TimeZone()
}

func TimeZoneAbv() string {
	return gofakeit.TimeZoneAbv()
}

func TimeZoneFull() string {
	return gofakeit.TimeZoneFull()
}

func TimeZoneOffset() float32 {
	return gofakeit.TimeZoneOffset()
}

func TimeZoneRegion() string {
	return gofakeit.TimeZoneRegion()
}

func DomainName() string {
	return gofakeit.DomainName()
}

func DomainSuffix() string {
	return gofakeit.DomainSuffix()
}

func IPv4Address() string {
	return gofakeit.IPv4Address()
}

func IPv6Address() string {
	return gofakeit.IPv6Address()
}

func MacAddress() string {
	return gofakeit.MacAddress()
}

func HTTPMethod() string {
	return gofakeit.HTTPMethod()
}

func HTTPStatusCode() int {
	return gofakeit.HTTPStatusCode()
}

func UserAgent() string {
	return gofakeit.UserAgent()
}

func ChromeUserAgent() string {
	return gofakeit.ChromeUserAgent()
}

func FirefoxUserAgent() string {
	return gofakeit.FirefoxUserAgent()
}

func SafariUserAgent() string {
	return gofakeit.SafariUserAgent()
}

func OperaUserAgent() string {
	return gofakeit.OperaUserAgent()
}

func Int8() int8 {
	return gofakeit.Int8()
}

func Int16() int16 {
	return gofakeit.Int16()
}

func Int32() int32 {
	return gofakeit.Int32()
}

func Int64() int64 {
	return gofakeit.Int64()
}

func Uint8() uint8 {
	return gofakeit.Uint8()
}

func Uint16() uint16 {
	return gofakeit.Uint16()
}

func Uint32() uint32 {
	return gofakeit.Uint32()
}

func Uint64() uint64 {
	return gofakeit.Uint64()
}

func Float32() float32 {
	return gofakeit.Float32()
}

func Float64() float64 {
	return gofakeit.Float64()
}

func Color() string {
	return gofakeit.Color()
}

func SafeColor() string {
	return gofakeit.SafeColor()
}

func RGBColor() []int {
	return gofakeit.RGBColor()
}

func FileExtension() string {
	return gofakeit.FileExtension()
}

func CreditCardNumber() string {
	cc := gofakeit.CreditCard()
	return cc.Number
}

func CreditCardType() string {
	cc := gofakeit.CreditCard()
	return cc.Type
}

func CreditCardExp() string {
	cc := gofakeit.CreditCard()
	return cc.Exp
}

func CreditCardCvv() string {
	cc := gofakeit.CreditCard()
	return cc.Cvv
}

func Currency() string {
	curr := gofakeit.Currency()
	return curr.Short
}

func CurrencyLong() string {
	curr := gofakeit.Currency()
	return curr.Long
}

func JobTitle() string {
	return gofakeit.JobTitle()
}

func JobDescriptor() string {
	return gofakeit.JobDescriptor()
}

func JobLevel() string {
	return gofakeit.JobLevel()
}

func Emoji() string {
	return gofakeit.Emoji()
}

func Language() string {
	return gofakeit.Language()
}

func ProgrammingLanguage() string {
	return gofakeit.ProgrammingLanguage()
}

func CompanySuffix() string {
	return gofakeit.CompanySuffix()
}

func AddressFull() string {
	addr := gofakeit.Address()
	return addr.Street + ", " + addr.City + ", " + addr.State + " " + addr.Zip
}

func HTTPVersion() string {
	return gofakeit.HTTPVersion()
}

func LogLevel() string {
	return gofakeit.LogLevel("general")
}

func Digit() string {
	return gofakeit.Digit()
}

func DigitN(n uint) string {
	return gofakeit.DigitN(n)
}

func Letter() string {
	return gofakeit.Letter()
}

func Numerify(str string) string {
	return gofakeit.Numerify(str)
}

func Word() string {
	return gofakeit.LoremIpsumWord()
}

func Sentence(wordCount int) string {
	return gofakeit.LoremIpsumSentence(wordCount)
}

func Paragraph(lp ParagraphOptions) string {
	return gofakeit.LoremIpsumParagraph(lp.ParagraphCount, lp.SentenceCount, lp.WordCount, lp.Separator)
}

func Regex(regexStr string) string {
	return gofakeit.Regex(regexStr)
}

func Map() map[string]interface{} {
	return gofakeit.Map()
}

func Slice(v interface{}) {
	gofakeit.Slice(v)
}

func Struct(v interface{}) error {
	return gofakeit.Struct(v)
}

func ErrorObject() error {
	return gofakeit.ErrorObject()
}

func ErrorDatabase() error {
	return gofakeit.ErrorDatabase()
}

func ErrorGRPC() error {
	return gofakeit.ErrorGRPC()
}

func ErrorHTTP() error {
	return gofakeit.ErrorHTTP()
}

func ErrorRuntime() error {
	return gofakeit.ErrorRuntime()
}

func ErrorValidation() error {
	return gofakeit.ErrorValidation()
}

func Fruit() string {
	return gofakeit.Fruit()
}

func Vegetable() string {
	return gofakeit.Vegetable()
}

func ShuffleInts(a []int) {
	gofakeit.ShuffleInts(a)
}

func ShuffleStrings(a []string) {
	gofakeit.ShuffleStrings(a)
}

func FileMimeType() string {
	return gofakeit.FileMimeType()
}
func Image(width, height int) *image.RGBA {
	return gofakeit.Image(width, height)
}
func ImageJpeg(width, height int) []byte {
	return gofakeit.ImageJpeg(width, height)
}
func ImagePng(width, height int) []byte {
	return gofakeit.ImagePng(width, height)
}

func InputName() string {
	return gofakeit.InputName()
}
func HTTPStatusCodeSimple() int {
	return gofakeit.HTTPStatusCodeSimple()
}

func Int() int {
	return gofakeit.Int()
}
func IntN(n int) int {
	return gofakeit.IntN(n)
}
func Uint() uint {
	return gofakeit.Uint()
}
func UintN(n uint) uint {
	return gofakeit.UintN(n)
}
func UintRange(min, max uint) uint {
	return gofakeit.UintRange(min, max)
}
func HexUint(bitSize int) string {
	return gofakeit.HexUint(bitSize)
}
func Unit() string {
	return gofakeit.Unit()
}
func Weighted(options []any, weights []float32) (any, error) {
	return gofakeit.Weighted(options, weights)
}
func ShuffleAnySlice(v any) {
	gofakeit.ShuffleAnySlice(v)
}
func LatitudeInRange(min, max float64) (float64, error) {
	return gofakeit.LatitudeInRange(min, max)
}
func LongitudeInRange(min, max float64) (float64, error) {
	return gofakeit.LongitudeInRange(min, max)
}
func NiceColors() []string {
	return gofakeit.NiceColors()
}
func LanguageBCP() string {
	return gofakeit.LanguageBCP()
}
func SvgString(opts SVGOpts) string {
	goOpts := &gofakeit.SVGOptions{Width: opts.Width, Height: opts.Height, Type: opts.Type}
	return gofakeit.Svg(goOpts)
}

func MarkdownSimple() (string, error) {
	return gofakeit.Markdown(nil)
}

func Error() error {
	return gofakeit.Error()
}
func ErrorHTTPClient() error {
	return gofakeit.ErrorHTTPClient()
}
func ErrorHTTPServer() error {
	return gofakeit.ErrorHTTPServer()
}

var _ = json.Marshal
