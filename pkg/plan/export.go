package plan

import (
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
)

// Export allows save plan to file
func (p *Plan) Export() error {
	if err := os.RemoveAll(p.dir); err != nil {
		return err
	}

	if err := p.exportManifest(); err != nil {
		return err
	}

	if err := p.exportValues(); err != nil {
		return err
	}

	if err := p.exportGraphMD(); err != nil {
		return err
	}

	return helper.SaveInterface(p.fullPath, p.body)
}

func (p *Plan) exportManifest() error {
	if len(p.body.Releases) == 0 {
		return nil
	}

	for k, v := range p.manifests {
		m := filepath.Join(p.dir, Manifest, string(k))

		f, err := helper.CreateFile(m)
		if err != nil {
			return err
		}

		_, err = f.WriteString(v)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Plan) exportGraphMD() error {
	f, err := helper.CreateFile(filepath.Join(p.dir, "graph.md"))
	if err != nil {
		return err
	}

	_, err = f.WriteString(p.graphMD)
	if err != nil {
		return err
	}

	return f.Close()
}

func (p *Plan) exportValues() error {
	if len(p.body.Releases) == 0 {
		return nil
	}

	found := false
	for i := 0; i < len(p.body.Releases)-1 && !found; i++ {
		for range p.body.Releases[i].Values {
			found = true
			break
		}
	}

	if !found {
		return nil
	}

	return os.Rename(
		filepath.Join(p.tmpDir, Values),
		filepath.Join(p.dir, Values),
	)
}

// IsExist returns true if planfile exists
func (p *Plan) IsExist() bool {
	return helper.IsExists(p.fullPath)
}

// IsManifestExist returns true if planfile exists
func (p *Plan) IsManifestExist() bool {
	return helper.IsExists(filepath.Join(p.dir, Manifest))
}
