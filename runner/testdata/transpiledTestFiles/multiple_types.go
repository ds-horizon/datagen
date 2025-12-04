package main

import (
	// archive
	"archive/tar"
	"archive/zip"

	// buf / bytes
	"bufio"
	"bytes"

	// compress
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"

	// container
	"container/heap"
	"container/list"
	"container/ring"

	// context
	"context"

	// crypto (selected; many more below)
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"

	// database
	"database/sql"
	"database/sql/driver"

	// embed (package name is embed; blank ref below)
	_ "embed"

	// encoding
	"encoding"
	"encoding/ascii85"
	"encoding/asn1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/csv"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"

	// errors & expvar
	"errors"
	"expvar"

	// flag, fmt
	"flag"
	"fmt"

	// hash
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"

	// html
	"html"
	htmltmpl "html/template"

	// image
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"

	// index
	"index/suffixarray"

	// io
	"io"
	"io/fs"
	"io/ioutil"

	// log
	"log"
	"log/slog"

	"cmp"

	// math
	"math"
	"math/big"
	"math/bits"
	"math/cmplx"
	mrand "math/rand"

	// mime
	"mime"
	"mime/multipart"
	"mime/quotedprintable"

	// net
	"net"
	"net/http"
	"net/http/cgi"
	"net/http/cookiejar"
	"net/http/fcgi"
	"net/http/httptest"
	"net/http/httptrace"
	"net/http/httputil"
	"net/mail"
	"net/netip"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net/smtp"
	"net/textproto"
	"net/url"

	// os
	"os"
	"os/exec"
	"os/signal"
	"os/user"

	// path
	"path"
	"path/filepath"

	// reflect/regexp
	"reflect"
	"regexp"
	"regexp/syntax"

	// sort/strconv/strings
	"sort"
	"strconv"
	"strings"

	// sync
	"sync"
	"sync/atomic"

	// syscall (portable API only here)
	"syscall"

	// text
	textscanner "text/scanner"
	"text/tabwriter"
	texttmpl "text/template"
	"text/template/parse"

	// time
	"time"

	// unicode
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	// unsafe
	"unsafe"
)

var (
	// archive
	_ = tar.Header{}
	_ = zip.File{}

	// buf / bytes
	_ = bufio.Reader{}
	_ = bytes.Buffer{}

	// compress
	_ = bzip2.NewReader
	_ = flate.NewReader
	_ = gzip.Writer{}
	_ = lzw.NewReader
	_ = zlib.NewReader

	// container
	_ = heap.Init
	_ = list.List{}
	_ = ring.Ring{}

	// context
	_ = context.Background

	// crypto
	_ crypto.Hash
	_ = crand.Reader
	_ = aes.BlockSize
	_ = cipher.NewGCM
	_ = des.BlockSize
	_ = dsa.Parameters{}
	_ = ecdsa.PublicKey{}
	_ = ed25519.PrivateKey{}
	_ = elliptic.P256
	_ = hmac.New
	_ = md5.New
	_ = rsa.GenerateKey
	_ = sha1.New
	_ = sha256.New
	_ = sha512.New
	_ = subtle.ConstantTimeCompare
	_ = tls.VersionTLS13
	_ = x509.Certificate{}
	_ = pkix.Name{}

	// database
	_               = sql.ErrNoRows
	_ driver.Valuer = nil

	// embed
	// (embed has no exported identifiers intended for direct use; the blank import above makes the compiler link it
	// but it has no side effects. Keeping no var-ref here is fine.)

	// encoding
	_ = encoding.BinaryMarshaler(nil)
	_ = ascii85.Encode
	_ = asn1.Marshal
	_ = base32.StdEncoding
	_ = base64.StdEncoding
	_ = binary.BigEndian
	_ = csv.Reader{}
	_ = gob.NewEncoder
	_ = hex.EncodeToString
	_ = json.Marshal
	_ = pem.Encode
	_ = xml.Marshal

	// errors & expvar
	_ = errors.New
	_ = expvar.NewInt

	// flag, fmt
	_ = flag.String
	_ = fmt.Println

	// hash
	_ = hash.Hash(nil)
	_ = adler32.New
	_ = crc32.New
	_ = crc64.New
	_ = fnv.New32

	// html
	_ = html.EscapeString
	_ = htmltmpl.Template{}

	// image
	_ = image.NewRGBA
	_ = color.RGBA{}
	_ = palette.Plan9
	_ = draw.Draw
	_ = gif.Decode
	_ = jpeg.Encode
	_ = png.Decode

	// index
	_ = suffixarray.New

	// io
	_       = io.Copy
	_ fs.FS = nil
	_       = ioutil.ReadFile

	// log
	_ = log.Println
	_ = slog.Any // Go 1.21 structured logging

	// maps/slices/cmp
	_ = cmp.Compare[int]

	// math
	_ = math.Pi
	_ = big.Int{}
	_ = bits.LeadingZeros
	_ = cmplx.Abs
	_ = mrand.Int

	// mime
	_ = mime.TypeByExtension
	_ = multipart.Writer{}
	_ = quotedprintable.NewReader

	// net
	_ = net.Dial
	_ = http.ListenAndServe
	_ = cgi.Handler{}
	_ = cookiejar.New
	_ = fcgi.Serve
	_ = httptest.NewServer
	_ = httptrace.WithClientTrace
	_ = httputil.DumpRequest
	_ = mail.ReadMessage
	_ = netip.Addr{}
	_ = rpc.NewServer
	_ = jsonrpc.NewServerCodec
	_ = smtp.SendMail
	_ = textproto.NewReader
	_ = url.Parse

	// os
	_ = os.Open
	_ = exec.Command
	_ = signal.Notify
	_ = user.Current

	// path
	_ = path.Join
	_ = filepath.Abs

	// reflect/regexp
	_ = reflect.TypeOf
	_ = regexp.MustCompile
	_ = syntax.Op(0)

	// sort/strconv/strings
	_ = sort.Sort
	_ = strconv.Itoa
	_ = strings.Split

	// sync
	_ = sync.Mutex{}
	_ = atomic.AddInt32

	// syscall
	_ = syscall.Getpid

	// text
	_ = textscanner.Scanner{}
	_ = tabwriter.NewWriter
	_ = texttmpl.Must
	_ = parse.Tree{}

	// time
	_ = time.Now

	// unicode
	_ = unicode.IsLetter
	_ = utf16.Encode
	_ = utf8.RuneCountInString

	// unsafe
	_ = unsafe.Sizeof(0)
)

var multiple_typesMetadata __dgi_Metadata = __dgi_Metadata{
	Count: 1,
	Tags:  map[string]string{},
}

func (cg *__datagen_multiple_typesGenerator) Metadata() __dgi_Metadata {
	return multiple_typesMetadata
}

type __datagen_multiple_types struct {
	id     int
	score  float64
	name   string
	active bool
}

type __datagen_multiple_typesGenerator struct {
	id      func(iter int) int
	score   func(iter int) float64
	name    func(iter int) string
	active  func(iter int) bool
	all     *__datagen_multiple_typesDataHolder
	datagen *__dgi_DataGenGenerators
}

type __datagen_multiple_typesDataHolder struct {
	id     []int
	score  []float64
	name   []string
	active []bool
}

func (cg *__datagen_multiple_typesGenerator) __gen_wrapper_active() func(iter int) bool {
	return func(iter int) bool {
		if iter < len(cg.all.active) {
			return cg.all.active[iter]
		}

		for i := len(cg.all.active); i <= iter; i++ {
			val := cg.__gen_active(i)
			cg.all.active = append(cg.all.active, val)
		}

		return cg.all.active[iter]
	}
}

func (self *__datagen_multiple_typesGenerator) __gen_active(iter int) bool {
	return iter%2 == 0
}

func (cg *__datagen_multiple_typesGenerator) __gen_wrapper_name() func(iter int) string {
	return func(iter int) string {
		if iter < len(cg.all.name) {
			return cg.all.name[iter]
		}

		for i := len(cg.all.name); i <= iter; i++ {
			val := cg.__gen_name(i)
			cg.all.name = append(cg.all.name, val)
		}

		return cg.all.name[iter]
	}
}

func (self *__datagen_multiple_typesGenerator) __gen_name(iter int) string {
	return fmt.Sprintf("user_%d", iter)
}

func (cg *__datagen_multiple_typesGenerator) __gen_wrapper_score() func(iter int) float64 {
	return func(iter int) float64 {
		if iter < len(cg.all.score) {
			return cg.all.score[iter]
		}

		for i := len(cg.all.score); i <= iter; i++ {
			val := cg.__gen_score(i)
			cg.all.score = append(cg.all.score, val)
		}

		return cg.all.score[iter]
	}
}

func (self *__datagen_multiple_typesGenerator) __gen_score(iter int) float64 {
	return FloatBetween(0.0, 100.0)
}

func (cg *__datagen_multiple_typesGenerator) __gen_wrapper_id() func(iter int) int {
	return func(iter int) int {
		if iter < len(cg.all.id) {
			return cg.all.id[iter]
		}

		for i := len(cg.all.id); i <= iter; i++ {
			val := cg.__gen_id(i)
			cg.all.id = append(cg.all.id, val)
		}

		return cg.all.id[iter]
	}
}

func (self *__datagen_multiple_typesGenerator) __gen_id(iter int) int {
	return iter
}

func (cg *__datagen_multiple_typesGenerator) Gen(iter int) __dgi_Record {
	return &__datagen_multiple_types{
		id:     cg.id(iter),
		score:  cg.score(iter),
		name:   cg.name(iter),
		active: cg.active(iter),
	}
}

func __init___datagen_multiple_typesGenerator() *__datagen_multiple_typesGenerator {
	all := &__datagen_multiple_typesDataHolder{}
	cg := &__datagen_multiple_typesGenerator{all: all}
	cg.id = cg.__gen_wrapper_id()
	cg.score = cg.__gen_wrapper_score()
	cg.name = cg.__gen_wrapper_name()
	cg.active = cg.__gen_wrapper_active()
	return cg
}

func (e *__datagen_multiple_types) ToCSV() []string {
	return []string{
		fmt.Sprintf("%v", e.id),
		fmt.Sprintf("%v", e.score),
		fmt.Sprintf("%v", e.name),
		fmt.Sprintf("%v", e.active),
	}
}

func (e *__datagen_multiple_types) CSVHeaders() []string {
	return []string{
		"id",
		"score",
		"name",
		"active",
	}
}

func (e *__datagen_multiple_types) ToJSON() string {
	data, err := json.Marshal(map[string]interface{}{
		"id":     e.id,
		"score":  e.score,
		"name":   e.name,
		"active": e.active,
	})
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func (e *__datagen_multiple_types) ToXML() string {
	type __dgi_xmlAlias struct {
		XMLName    xml.Name `xml:"multiple_types"`
		Xml_id     int      `xml:"id"`
		Xml_score  float64  `xml:"score"`
		Xml_name   string   `xml:"name"`
		Xml_active bool     `xml:"active"`
	}

	data := __dgi_xmlAlias{
		Xml_id:     e.id,
		Xml_score:  e.score,
		Xml_name:   e.name,
		Xml_active: e.active,
	}

	xmlData, err := xml.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(xmlData)
}

func (self *__datagen_multiple_types) __dgi_Serialise() []byte {
	return []byte{}
}
