package analyzer

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
)

var goRoot string

type PackageIndex struct {
	Packages []*Package
	nameMap  map[string]*Package
	charMap  map[string][]*CompletionItem
}

func (pi PackageIndex) PackageByName(name string) (*Package, bool) {
	pkg, ok := pi.nameMap[name]
	return pkg, ok
}

func (pi PackageIndex) Match(char string) []*CompletionItem {
	return pi.charMap[char]
}

func (pi PackageIndex) Len() int {
	return len(pi.Packages)
}

func BuildPackageIndex(pkgs []*Package) PackageIndex {
	idx := PackageIndex{
		Packages: pkgs,
		nameMap:  make(map[string]*Package, len(pkgs)),
		charMap:  make(map[string][]*CompletionItem),
	}

	for _, pkg := range pkgs {
		addPackageToIndex(pkg, &idx)
	}

	return idx
}

func addPackageToIndex(pkg *Package, idx *PackageIndex) {
	firstChar := pkg.Name[:1]
	idx.nameMap[pkg.Name] = pkg
	idx.charMap[firstChar] = append(idx.charMap[firstChar], pkg.GetCompletionItem())

	if len(pkg.Children) > 0 {
		for _, child := range pkg.Children {
			addPackageToIndex(child, idx)
		}
	}
}

func SetRoot(root string) {
	goRoot = root
}

func SetLogger(l *zap.Logger) {
	log = l.Named("analyzer").Sugar()
}

func ReadPackagesFile(f string) ([]*Package, error) {
	r, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	defer r.Close()
	var pkgs []*Package
	err = json.NewDecoder(r).Decode(&pkgs)
	return pkgs, err
}
