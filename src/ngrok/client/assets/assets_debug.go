// Code generated by go-bindata.
// sources:
// assets/client/page.html
// assets/client/static/css/bootstrap.min.css
// assets/client/static/css/highlight.min.css
// assets/client/static/img/glyphicons-halflings.png
// assets/client/static/js/angular-sanitize.min.js
// assets/client/static/js/angular.js
// assets/client/static/js/base64.js
// assets/client/static/js/highlight.min.js
// assets/client/static/js/jquery-1.9.1.min.js
// assets/client/static/js/jquery.timeago.js
// assets/client/static/js/ngrok.js
// assets/client/static/js/vkbeautify.js
// assets/client/tls/ngrokroot.crt
// assets/client/tls/snakeoilca.crt
// DO NOT EDIT!

// +build !release

package assets

import (
	"fmt"
	"io/ioutil"
	"os"
	p "path"
	"path/filepath"
	"strings"
)

const assets_path = "D:\\projects\\mgrok\\assets\\"

// bindataRead reads the given file from disk. It returns an error on failure.
func bindataRead(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

func assetsPath() string {
	return assets_path
}

// assetsClientPageHtml reads file data from disk. It returns an error on failure.
func assetsClientPageHtml() (*asset, error) {
	path := p.Join(assetsPath(), "client\\page.html") //"D:\\projects\\mgrok\\assets\\client\\page.html"
	name := "assets/client/page.html"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticCssBootstrapMinCss reads file data from disk. It returns an error on failure.
func assetsClientStaticCssBootstrapMinCss() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/css/bootstrap.min.css"
	path := p.Join(assetsPath(), "client/static/css/bootstrap.min.css")
	name := "assets/client/static/css/bootstrap.min.css"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticCssHighlightMinCss reads file data from disk. It returns an error on failure.
func assetsClientStaticCssHighlightMinCss() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/css/highlight.min.css"
	path := p.Join(assetsPath(), "client/static/css/highlight.min.css")
	name := "assets/client/static/css/highlight.min.css"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticImgGlyphiconsHalflingsPng reads file data from disk. It returns an error on failure.
func assetsClientStaticImgGlyphiconsHalflingsPng() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/img/glyphicons-halflings.png"
	path := p.Join(assetsPath(), "client/static/img/glyphicons-halflings.png")
	name := "assets/client/static/img/glyphicons-halflings.png"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticJsAngularSanitizeMinJs reads file data from disk. It returns an error on failure.
func assetsClientStaticJsAngularSanitizeMinJs() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/js/angular-sanitize.min.js"
	path := p.Join(assetsPath(), "client/static/js/angular-sanitize.min.js")
	name := "assets/client/static/js/angular-sanitize.min.js"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticJsAngularJs reads file data from disk. It returns an error on failure.
func assetsClientStaticJsAngularJs() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/js/angular.js"
	path := p.Join(assetsPath(), "client/static/js/angular.js")
	name := "assets/client/static/js/angular.js"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticJsBase64Js reads file data from disk. It returns an error on failure.
func assetsClientStaticJsBase64Js() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/js/base64.js"
	path := p.Join(assetsPath(), "client/static/js/base64.js")
	name := "assets/client/static/js/base64.js"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticJsHighlightMinJs reads file data from disk. It returns an error on failure.
func assetsClientStaticJsHighlightMinJs() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/js/highlight.min.js"
	path := p.Join(assetsPath(), "client/static/js/highlight.min.js")
	name := "assets/client/static/js/highlight.min.js"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticJsJquery191MinJs reads file data from disk. It returns an error on failure.
func assetsClientStaticJsJquery191MinJs() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/js/jquery-1.9.1.min.js"
	path := p.Join(assetsPath(), "client/static/js/jquery-1.9.1.min.js")
	name := "assets/client/static/js/jquery-1.9.1.min.js"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticJsJqueryTimeagoJs reads file data from disk. It returns an error on failure.
func assetsClientStaticJsJqueryTimeagoJs() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/js/jquery.timeago.js"
	path := p.Join(assetsPath(), "client/static/js/jquery.timeago.js")
	name := "assets/client/static/js/jquery.timeago.js"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticJsNgrokJs reads file data from disk. It returns an error on failure.
func assetsClientStaticJsNgrokJs() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/js/ngrok.js"
	path := p.Join(assetsPath(), "client/static/js/ngrok.js")
	name := "assets/client/static/js/ngrok.js"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientStaticJsVkbeautifyJs reads file data from disk. It returns an error on failure.
func assetsClientStaticJsVkbeautifyJs() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/static/js/vkbeautify.js"
	path := p.Join(assetsPath(), "client/static/js/vkbeautify.js")
	name := "assets/client/static/js/vkbeautify.js"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientTlsNgrokrootCrt reads file data from disk. It returns an error on failure.
func assetsClientTlsNgrokrootCrt() (*asset, error) {
	// path := "D:/projects/mgrok/assets/client/tls/ngrokroot.crt"
	path := p.Join(assetsPath(), "client/tls/ngrokroot.crt")
	name := "assets/client/tls/ngrokroot.crt"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// assetsClientTlsSnakeoilcaCrt reads file data from disk. It returns an error on failure.
func assetsClientTlsSnakeoilcaCrt() (*asset, error) {
	// path := "D:\\projects\\mgrok\\assets\\client\\tls\\snakeoilca.crt"
	path := p.Join(assetsPath(), "client/tls/snakeoilca.crt")
	name := "assets/client/tls/snakeoilca.crt"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"assets/client/page.html":                           assetsClientPageHtml,
	"assets/client/static/css/bootstrap.min.css":        assetsClientStaticCssBootstrapMinCss,
	"assets/client/static/css/highlight.min.css":        assetsClientStaticCssHighlightMinCss,
	"assets/client/static/img/glyphicons-halflings.png": assetsClientStaticImgGlyphiconsHalflingsPng,
	"assets/client/static/js/angular-sanitize.min.js":   assetsClientStaticJsAngularSanitizeMinJs,
	"assets/client/static/js/angular.js":                assetsClientStaticJsAngularJs,
	"assets/client/static/js/base64.js":                 assetsClientStaticJsBase64Js,
	"assets/client/static/js/highlight.min.js":          assetsClientStaticJsHighlightMinJs,
	"assets/client/static/js/jquery-1.9.1.min.js":       assetsClientStaticJsJquery191MinJs,
	"assets/client/static/js/jquery.timeago.js":         assetsClientStaticJsJqueryTimeagoJs,
	"assets/client/static/js/ngrok.js":                  assetsClientStaticJsNgrokJs,
	"assets/client/static/js/vkbeautify.js":             assetsClientStaticJsVkbeautifyJs,
	"assets/client/tls/ngrokroot.crt":                   assetsClientTlsNgrokrootCrt,
	"assets/client/tls/snakeoilca.crt":                  assetsClientTlsSnakeoilcaCrt,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"assets": &bintree{nil, map[string]*bintree{
		"client": &bintree{nil, map[string]*bintree{
			"page.html": &bintree{assetsClientPageHtml, map[string]*bintree{}},
			"static": &bintree{nil, map[string]*bintree{
				"css": &bintree{nil, map[string]*bintree{
					"bootstrap.min.css": &bintree{assetsClientStaticCssBootstrapMinCss, map[string]*bintree{}},
					"highlight.min.css": &bintree{assetsClientStaticCssHighlightMinCss, map[string]*bintree{}},
				}},
				"img": &bintree{nil, map[string]*bintree{
					"glyphicons-halflings.png": &bintree{assetsClientStaticImgGlyphiconsHalflingsPng, map[string]*bintree{}},
				}},
				"js": &bintree{nil, map[string]*bintree{
					"angular-sanitize.min.js": &bintree{assetsClientStaticJsAngularSanitizeMinJs, map[string]*bintree{}},
					"angular.js":              &bintree{assetsClientStaticJsAngularJs, map[string]*bintree{}},
					"base64.js":               &bintree{assetsClientStaticJsBase64Js, map[string]*bintree{}},
					"highlight.min.js":        &bintree{assetsClientStaticJsHighlightMinJs, map[string]*bintree{}},
					"jquery-1.9.1.min.js":     &bintree{assetsClientStaticJsJquery191MinJs, map[string]*bintree{}},
					"jquery.timeago.js":       &bintree{assetsClientStaticJsJqueryTimeagoJs, map[string]*bintree{}},
					"ngrok.js":                &bintree{assetsClientStaticJsNgrokJs, map[string]*bintree{}},
					"vkbeautify.js":           &bintree{assetsClientStaticJsVkbeautifyJs, map[string]*bintree{}},
				}},
			}},
			"tls": &bintree{nil, map[string]*bintree{
				"ngrokroot.crt":  &bintree{assetsClientTlsNgrokrootCrt, map[string]*bintree{}},
				"snakeoilca.crt": &bintree{assetsClientTlsSnakeoilcaCrt, map[string]*bintree{}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
