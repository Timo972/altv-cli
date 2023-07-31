package altcdn

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type ModuleMap map[string]string

type altCDN struct {
	BaseURL  string
	includes ModuleMap
}

var BaseURL = "https://cdn.alt-mp.com"
var DefaultModules = ModuleMap{
	"server":             "server/%s/%s",
	"data-files":         "data/%s",
	"js-module":          "js-module/%s/%s",
	"csharp-module":      `coreclr-module/%s/%s`,
	"js-bytecode-module": `js-bytecode-module/%s/%s`,
}

var Default cdn.CDN = &altCDN{
	BaseURL:  BaseURL,
	includes: DefaultModules,
}

func New(baseURL string, modules ModuleMap) *altCDN {
	return &altCDN{
		BaseURL:  BaseURL,
		includes: modules,
	}
}

func (c *altCDN) SetBaseURL(baseURL string) {
	c.BaseURL = baseURL
}

func (c *altCDN) SetModules(modules ModuleMap) {
	c.includes = modules
}

func (c *altCDN) Has(module string) bool {
	_, ok := c.includes[module]
	return ok
}

func (c *altCDN) Manifest(branch version.Branch, arch platform.Arch, module string) (*cdn.Manifest, error) {
	manUrl := c.fileURL(branch, arch, module, "update.json")
	logging.DebugLogger.Printf("Fetching manifest from %s", manUrl)

	resp, err := http.DefaultClient.Get(manUrl)
	logging.DebugLogger.Printf("Got response: %s", resp.Status)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logging.DebugLogger.Printf("Failed to fetch manifest from %s: %s", manUrl, resp.Status)
		return nil, fmt.Errorf("failed to fetch manifest from %s: %s", manUrl, resp.Status)
	}

	var manifest cdn.Manifest
	if err = json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %w", err)
	}

	logging.DebugLogger.Printf("Manifest: %+v", manifest)

	return &manifest, nil
}

func (c *altCDN) fileURL(branch version.Branch, arch platform.Arch, module string, name string) string {
	if module == "data-files" {
		return fmt.Sprintf("%s/"+c.includes[module]+"/%s", c.BaseURL, branch, name)
	}
	return fmt.Sprintf("%s/"+c.includes[module]+"/%s", c.BaseURL, branch, arch, name)
}

func (c *altCDN) Files(branch version.Branch, arch platform.Arch, module string, manifest bool) ([]*cdn.File, error) {
	man, err := c.Manifest(branch, arch, module)
	if err != nil {
		return nil, fmt.Errorf("unable to gather module files: %w", err)
	}

	logging.DebugLogger.Println("got manifest")

	fileCount := len(man.HashList)
	if manifest {
		fileCount++
	}

	files := make([]*cdn.File, fileCount)
	i := 0
	if manifest {
		files[i] = &cdn.File{
			Name: "update.json",
			Hash: "",
			Url:  c.fileURL(branch, arch, module, "update.json"),
		}
		i++
	}

	for name, hash := range man.HashList {
		logging.DebugLogger.Printf("adding file %s", name)
		files[i] = &cdn.File{
			Name: name,
			Hash: hash,
			Url:  c.fileURL(branch, arch, module, name),
		}
		i++
	}

	logging.DebugLogger.Println("extracted files")

	return files, nil
}
