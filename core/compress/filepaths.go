package compress

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type FilePaths struct {
	In  string
	Out string

	InLock *Lock

	Done    *Lock
	OutLock *Lock
}

func (fp *FilePaths) String() string {
	out := "&compress.FilePaths{\n"
	out += fmt.Sprintf("\t%s:%q,\n", "In", fp.In)
	out += fmt.Sprintf("\t%s:%q,\n", "Out", fp.Out)
	out += fmt.Sprintf("\t%s:%q,\n", "InLock", fp.InLock)
	out += fmt.Sprintf("\t%s:%q,\n", "Done", fp.Done)
	out += fmt.Sprintf("\t%s:%q,\n", "OutLock", fp.OutLock)
	out += "}"
	return out
}

func (p *Pipeline) getFilePathsByInfile(name string) *FilePaths {
	fp := &FilePaths{
		In: name,
	}

	fn := filepath.Base(name)

	fp.Out = fmt.Sprintf("%s/%s.%s", p.Outdir, strings.TrimSuffix(fn, filepath.Ext(fn)), p.OutfileExt)

	fp.InLock = NewLock(fp.In + ".lock")

	fp.Done = NewLock(fp.Out + ".done")
	fp.OutLock = NewLock(fp.Out + ".lock")

	return fp
}

func (p *Pipeline) getFilePathsByOutfile(name string) *FilePaths {
	fp := &FilePaths{
		Out: name,
	}

	fn := filepath.Base(name)
	fp.In = fmt.Sprintf("%s/%s.%s", p.Indir, strings.TrimSuffix(fn, filepath.Ext(fn)), p.InfileExt)
	fp.InLock = NewLock(fp.In + ".lock")
	fp.Done = NewLock(fp.Out + ".done")
	fp.OutLock = NewLock(fp.Out + ".lock")

	return fp
}

func (p *Pipeline) getFilePathsByIndir() ([]*FilePaths, error) {
	files, err := ioutil.ReadDir(p.Indir)
	if err != nil {
		return nil, err
	}

	var fps []*FilePaths

	for _, fi := range files {
		// Filter out any named files not ending with InfileExt
		if filepath.Ext(fi.Name()) != "."+p.InfileExt {
			continue
		}
		fp := p.getFilePathsByInfile(fmt.Sprintf("%s/%s", p.Indir, fi.Name()))

		fps = append(fps, fp)
	}
	return fps, nil
}

func (p *Pipeline) getFilePathsByOutdir() ([]*FilePaths, error) {
	files, err := ioutil.ReadDir(p.Outdir)
	if err != nil {
		return nil, err
	}

	var fps []*FilePaths

	for _, fi := range files {
		// Filter out any named files not ending with OutfileExt
		if filepath.Ext(fi.Name()) != "."+p.OutfileExt {
			continue
		}
		fp := p.getFilePathsByOutfile(fmt.Sprintf("%s/%s", p.Outdir, fi.Name()))

		fps = append(fps, fp)
	}
	return fps, nil
}
