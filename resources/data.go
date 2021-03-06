// Code generated by go-bindata.
// sources:
// ../build/configs/default_config.yaml
// ../build/configs/testing_config.yaml
// DO NOT EDIT!

package resources

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _goCentrifugeBuildConfigsDefault_configYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x56\x49\x77\xdb\x3a\x0f\xdd\xeb\x57\xe0\x38\x9b\xef\x5b\xd4\xd6\x60\xc9\xb6\x76\x19\x3b\x39\x79\xb6\xe3\x36\x4d\x76\x14\x09\x49\xac\x25\x52\x21\x29\xcb\xee\xaf\x7f\x87\x1a\xd2\x0c\x1d\x5e\x56\x3a\x01\x70\x09\x5c\x5c\x00\x3e\x81\x0b\x4c\x49\x5d\x18\x60\xb8\xc7\x42\x56\x25\x0a\x03\x06\xb5\x11\x68\x80\x64\x84\x0b\x6d\x40\x71\xb1\xc3\xe4\xe8\x50\x14\x46\xf1\xb4\xce\xf0\x06\x4d\x23\xd5\x2e\x06\x55\x6b\xcd\x89\xc8\x79\x51\x38\x2d\x18\x17\x08\x26\x47\x60\x3d\xae\xe8\x3c\x35\x98\x9c\x18\x38\x7f\x42\x80\x92\x70\x61\x2c\xbe\x33\xb8\xc4\x0e\xc0\x09\x2c\x25\x25\x45\x9b\x02\x17\x19\x50\x29\x8c\x22\xd4\x00\x61\x4c\xa1\xd6\xa8\x41\x20\x32\x30\x12\x12\x04\x8d\x06\x1a\x6e\x72\x40\xb1\x87\x3d\x51\x9c\x24\x05\xea\xb1\x03\x43\xbc\x85\x04\xe0\x2c\x86\x20\x08\xda\x6f\x34\x39\x2a\xac\xcb\xbe\x82\x8f\x2c\x86\x79\x30\xef\x6c\x89\x94\x46\x1b\x45\xaa\x15\xa2\xd2\x5d\x2c\xc0\x3b\x18\x4d\x78\x35\x9d\x78\xfe\x6c\xec\x8e\xdd\xb1\x37\x31\xb4\x9a\x04\x73\xdf\xf5\x27\xbc\x4a\xf5\x64\x5d\x6e\xd7\x87\xa4\xd9\xd5\x0f\xf7\xf7\x17\x69\xfd\x63\x9b\x1c\x2e\x4f\x37\xb8\xbd\x39\x5f\xca\x1f\xc7\x63\x18\xce\xf7\x6b\x91\x7d\xdd\xaf\xae\xbf\x2f\xef\x77\xa3\xbf\xc2\x06\x03\xec\xd7\x34\xba\xbc\x89\xca\xdd\xe3\x1d\x7e\xbf\xfb\x7c\xe7\x3f\xae\x6a\x2f\xfa\x56\xb1\xf7\xc1\xee\x93\xf4\xb6\x41\x99\x93\x7c\x75\x16\xde\x62\x28\xbc\x0e\x76\xa0\xeb\x74\x60\x6b\x28\x82\x33\x14\x86\x9b\xe3\x15\xa1\x46\xaa\x63\x0c\xa3\xd1\x2b\xcb\x06\x33\xae\xcd\x0b\x13\x11\x34\x97\x6a\x83\x95\xd4\xfc\x55\x54\x45\x8e\x56\x2a\xff\x24\x05\xcf\x88\xe1\x52\xb4\xb6\xb6\x81\xd7\x84\x8b\x5f\xca\xa9\xef\xb3\x03\xcf\x55\xd3\x25\x78\x02\x37\x75\x89\x8a\x53\xf8\x78\x01\x32\x6d\x15\xf4\x4c\x2b\x3f\x23\xbb\x66\x86\x5e\x1f\x75\x36\x74\x0c\x0a\xae\x8d\x8d\x14\x92\xe1\x5b\xb1\x55\x4a\xee\x79\x6b\x90\x2d\xf6\xb3\x04\x86\xf4\xfe\x83\x02\x82\x70\xec\xfb\xe1\xd8\x77\xdd\xf1\xd4\x7f\xad\x02\xcf\xbf\x08\x3e\x4b\x79\xb7\xe4\x9c\xae\xbf\x36\xdb\x7c\x7b\x76\x1f\x1d\x3e\xd3\x95\x5c\xa6\xd1\x66\x7d\xff\xe9\xaa\x6a\x52\x4f\xcd\xc2\x66\x79\xf0\x1f\x36\x41\x75\xce\xbc\xd7\x5a\xe8\x1f\x98\x47\x63\xdf\x73\x7f\xf7\xc0\xfa\xe1\xfa\x74\xfe\x7e\xf5\x41\xed\x2f\x1f\xce\x16\x0d\xdb\xc9\x2f\xf4\xf4\xb4\x3c\x7f\xf8\x50\x2d\xf0\x78\x7c\x98\xde\x5e\xce\xb3\x2b\x15\xe4\xdb\x9b\x6f\xa3\x9e\xa7\xcb\x5e\xf5\x03\x93\x96\xe6\x77\xb0\xe9\xe7\xfa\x37\x73\x31\xed\x83\x97\xc4\x52\x04\x0c\xab\x42\x1e\x91\xc1\x6d\x49\x94\x81\xf3\x5e\x6a\x1a\x52\xa9\x5a\x52\x33\xbe\x47\xf1\x82\xce\xb7\x72\x84\xdf\xea\xd1\x3d\x2c\x5c\xe6\x2f\xa6\xe1\xcc\xc3\x59\x30\x9f\xfa\xd1\x62\x46\xa2\x28\x99\x91\xc5\x82\xb8\x0b\xc6\x22\x3a\x0b\x58\x10\x46\xec\x0f\xca\x75\x0f\x8b\x28\x72\xa9\x1b\x2c\x58\xe0\x79\xd3\x30\x20\xa9\xcb\xc2\x39\x0d\xa3\x28\x9a\xf9\x01\x5b\x50\x3f\x25\x33\x16\x21\xfd\x83\xc6\xdd\xc3\x2c\x9d\x87\x53\x96\x92\xc5\xdc\xf5\x7c\x36\x4b\x49\x18\xd2\xb9\x1b\x24\x09\xf1\xfd\xc8\x4d\x28\x43\x9c\x26\x21\xb2\xbf\x4c\xc3\x09\xac\x6b\xac\xd1\xd2\x90\xf2\xac\x56\xad\xa9\xe3\x8b\xe8\xa3\xa0\xb9\x92\x42\xd6\xda\xaa\x93\xa2\xd6\x5c\x64\xce\xa3\x0d\xe8\x16\x61\xb7\x4c\x75\x4b\xad\xa8\xcb\x04\x95\xd5\xb7\x6d\x0e\x2a\x3d\xa1\x52\x68\x3b\x32\xbd\xd6\x1b\xab\xe4\x04\x81\x14\x85\xa4\xc4\x20\x03\x62\x40\x1b\xa2\x4c\x5d\x39\x60\xe3\xef\xba\xc0\x18\xfc\x16\xfd\x4a\x21\x6a\xa8\x2b\x38\x5f\x7d\x01\x7a\xa4\x05\x6a\x68\x72\x14\xfd\x03\xc0\x35\x34\x84\xb7\x3b\xd8\xe6\x8b\x7b\x14\x46\x3b\xd0\x9b\xef\x08\x37\x5b\x5e\xe2\xf5\x6d\x0c\x9e\x2d\xf4\x49\x60\xba\x42\xca\x53\x4e\x5f\x16\xed\x0c\xf2\xea\x4a\xbb\xc5\x02\xad\x72\x9a\x9c\xd3\xfc\x49\x7a\x40\x28\x95\xb5\x5d\x18\x12\x6a\x8d\xc3\x1e\x90\x96\x84\x7e\x80\x19\x70\xd1\xfe\x93\xd6\xda\xc8\xb2\x7f\x04\x52\x5e\xa0\x03\xc3\xcd\x39\xed\x60\x6e\x48\x89\x31\x8c\xec\x9d\x19\x3d\x5d\x16\x9b\xcc\x00\xfc\xf4\x2e\x2d\xb8\xdd\x53\x76\x75\xc0\xff\x1a\x04\x85\x8f\x35\x57\x08\x8d\x06\xa9\x80\x57\xb4\x3f\x37\xf6\xba\xd8\x4f\x4a\x8c\x4d\xbb\xa5\xe4\xff\x96\x5d\xc9\xf0\xcb\x66\x19\x43\xa3\xe3\xc9\xc4\x36\xa0\xc8\xa5\x36\xf1\x22\x9c\x46\x43\x2b\xdb\x63\x98\x11\x5b\x09\xa7\x36\xd9\x8c\xe8\x95\xfd\x8c\xc1\x73\x87\xbf\x37\xce\x05\x2f\xb9\xe9\x9c\x97\xf6\x33\x86\xe9\xcc\xf3\x83\xf9\xbc\xf5\xb4\x1d\x90\x75\xcb\x97\xed\x55\x27\x2c\xf1\xb3\x2e\xa3\x88\xd0\x84\xb6\x35\xf7\x15\x30\xd6\x1d\x4f\x02\x49\x21\xe9\x0e\x88\x60\x7d\x21\x60\x14\xcf\x32\x54\xc8\x9c\x6e\x70\xf1\x60\x86\x36\xcb\xda\xc4\x30\x8a\x5c\x57\x77\x4c\xde\x58\xa6\x9e\xa3\x57\x52\x16\x50\x92\x03\x28\x34\x8a\x77\x7b\x56\xa3\x60\x40\x5e\xb8\xc9\x3d\x2a\x07\xac\xe3\xa6\xf3\x8b\xc1\xef\xab\xfe\x35\x24\x17\x06\xd5\x9e\x14\x2d\xee\xb1\x13\x28\xb1\xd9\xd1\x5a\xa9\xf6\xb6\x3c\x8b\xc8\x89\x86\x04\xd1\x1e\x1f\x83\xd4\xb4\x85\x0c\x00\xf6\x3d\x3b\xdb\x7e\x5f\xc1\x05\xd7\x6d\x37\x5b\x44\x2d\xcb\x37\x6a\xd0\xc0\x24\x08\x69\x40\xd7\x55\x25\x95\x01\x73\x68\x33\x22\x15\xb7\x3f\x2f\x0e\x2b\x29\x8b\x53\x6a\xc7\xf6\x52\x58\x24\x16\x83\x51\x35\x3a\xce\xbf\x01\x00\x00\xff\xff\x75\xa8\x4f\xcf\x53\x09\x00\x00")

func goCentrifugeBuildConfigsDefault_configYamlBytes() ([]byte, error) {
	return bindataRead(
		_goCentrifugeBuildConfigsDefault_configYaml,
		"go-centrifuge/build/configs/default_config.yaml",
	)
}

func goCentrifugeBuildConfigsDefault_configYaml() (*asset, error) {
	bytes, err := goCentrifugeBuildConfigsDefault_configYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "go-centrifuge/build/configs/default_config.yaml", size: 2387, mode: os.FileMode(420), modTime: time.Unix(1540471826, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _goCentrifugeBuildConfigsTesting_configYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\xc9\x6e\xe3\x48\x0c\x86\xef\x7a\x0a\x81\x97\x5c\xbc\xd4\xbe\xbd\xc1\x20\x98\xd3\x0c\x90\x33\xab\xc8\x8a\x05\xdb\xb2\x5a\x4b\x12\x23\xc8\xbb\x37\xe4\x38\x9d\x6b\x1a\xba\x90\x04\x7f\xfe\xa4\xea\xe3\xf9\xc0\x23\x2f\xe7\xd4\xb4\x2d\x96\x72\x59\xfa\x79\x5a\xe3\xb6\x3d\x63\xd7\xa7\xf6\x16\xb6\xed\x91\xaf\xa9\x7d\x78\x07\x24\x1a\x79\x9a\x20\x41\x88\x59\x60\x70\x36\xe8\x62\x8c\x31\x58\x2a\x79\x99\x8d\xd3\x2c\x48\x17\x6b\x91\xa5\x91\x0a\x2d\x6c\xa0\x8c\xd7\x61\xbe\x40\x7a\x87\xd2\x0d\x07\x1e\x21\x01\xf2\xb4\x95\x2a\x6c\xcb\x3c\xae\x0d\xb7\xf2\xcc\x6f\x33\x24\x28\xde\xc7\x1a\xb4\x8f\xe4\xbd\xa0\xa8\x4a\x2d\x92\x88\x0c\x86\xaa\x25\x59\x14\x48\x25\x54\x85\x22\x2b\x94\x46\x48\xed\x05\x69\xa7\x45\xd5\xa1\x88\x12\xf0\xcf\xbc\x01\x47\x3c\x4f\xab\x6d\xf7\x02\x09\xb4\x2b\xd2\x05\xf6\x3a\xd7\x18\x44\x65\x6f\xb3\xf0\xca\xd7\x10\x05\x7a\x89\x04\x1f\x1b\x38\x52\x85\x04\xd3\x6d\x61\xb8\xa5\xdf\x43\xe8\x78\xe2\x1e\x92\x56\x1b\xe8\x21\x29\xa7\xa4\x31\x1b\x18\x20\xc9\x0d\x8c\x90\xc2\x06\x26\x3c\xad\x07\x10\xcb\xcc\xd2\xb1\x2e\x31\xc8\x68\x0c\x49\x2e\xa8\x72\xc8\xca\xb3\x61\xc7\x22\xdb\x5c\xb3\xd1\x99\x85\xf6\x0e\x2d\x85\x10\x62\x45\xe7\x23\xaa\x20\x95\x5a\x17\x39\x63\x59\x7f\x45\x91\x2a\xe4\x20\xad\xb5\x36\xa3\x64\x24\x5f\x90\xa3\x70\x82\x43\x30\x0a\x6b\xc1\xa0\xad\x23\xe1\x8c\xb5\x99\x22\x5a\x6f\x55\x46\x57\x4b\x11\x51\x71\x5d\x27\x75\x04\x09\x8c\x65\xe1\x04\xba\x2d\x29\xe4\xad\xd1\x39\x6c\xa3\x52\x75\x6b\x4c\x50\xd1\xc4\x48\xda\x13\x6c\xe0\x85\xc7\xa9\xbb\xac\x47\x7e\x3c\xdc\x1f\x7e\xc0\x69\x7a\xbd\x8c\x94\xda\x87\xaf\xd2\x9d\x81\xd4\xfe\x14\x81\xa6\xe9\x88\xfb\xb9\x9b\xaf\xff\x50\x6a\x41\xbc\x09\xf9\xfd\x41\xd3\xfc\x5a\x78\xe1\x15\xba\x7e\x39\x3f\x5d\xc6\x23\x8f\x53\x6a\x55\xd3\xb6\xaf\xb7\xe4\x09\xbb\xf9\xff\xee\xcc\xff\xfe\x97\x5a\xd9\x34\x47\xbe\xde\x08\x9d\xba\xe7\xbe\xeb\x9f\x3f\x61\x1d\x96\x7c\xea\xca\xe3\x4a\xe9\x6e\xb7\xdf\xed\xf6\x79\xe9\x4e\xb4\x1f\x79\xba\x2c\x63\xe1\x69\x7f\xef\x7e\xe4\xeb\x6e\x58\xf2\x6e\xe0\xf3\xa7\x6e\xec\x5e\x70\xe6\x9f\x09\x8f\xab\xf8\x26\xe4\xf9\x80\xcb\x7c\xf8\xa1\xf7\xbd\xfb\x2f\x8d\xbf\x54\x5f\xae\xbf\x03\x00\x00\xff\xff\xb0\x1c\xaf\x3f\xaa\x03\x00\x00")

func goCentrifugeBuildConfigsTesting_configYamlBytes() ([]byte, error) {
	return bindataRead(
		_goCentrifugeBuildConfigsTesting_configYaml,
		"go-centrifuge/build/configs/testing_config.yaml",
	)
}

func goCentrifugeBuildConfigsTesting_configYaml() (*asset, error) {
	bytes, err := goCentrifugeBuildConfigsTesting_configYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "go-centrifuge/build/configs/testing_config.yaml", size: 938, mode: os.FileMode(420), modTime: time.Unix(1540471826, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
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
	"go-centrifuge/build/configs/default_config.yaml": goCentrifugeBuildConfigsDefault_configYaml,
	"go-centrifuge/build/configs/testing_config.yaml": goCentrifugeBuildConfigsTesting_configYaml,
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
	"go-centrifuge": &bintree{nil, map[string]*bintree{
		"build": &bintree{nil, map[string]*bintree{
			"configs": &bintree{nil, map[string]*bintree{
				"default_config.yaml": &bintree{goCentrifugeBuildConfigsDefault_configYaml, map[string]*bintree{}},
				"testing_config.yaml": &bintree{goCentrifugeBuildConfigsTesting_configYaml, map[string]*bintree{}},
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
