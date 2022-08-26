package compress

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Find & delete lockfiles in Indir
func (p *Pipeline) ReleaseLocks() error {
	// Release all locks in Indir
	files, err := ioutil.ReadDir(p.Indir)
	if err != nil {
		return err
	}
	for _, fi := range files {
		// Filter .lock files
		if filepath.Ext(fi.Name()) != ".lock" {
			continue
		}
		lockfile := p.Indir + "/" + fi.Name()
		infile := strings.TrimSuffix(lockfile, ".lock")

		fp := p.getFilePathsByInfile(infile)

		if isFileBusy(fp.In) {
			fmt.Printf("ReleaseLocks(): Skip: %q is busy, leaving lockfile\n", fp.In)
			continue
		}
		fmt.Printf("ReleaseLocks(): Removing infile lock %q\n", fp.InLock)
		fp.InLock.Release()
	}

	// Release all locks in Outdir
	files, err = ioutil.ReadDir(p.Outdir)
	if err != nil {
		return err
	}
	for _, fi := range files {
		// Filter out non lock files
		if filepath.Ext(fi.Name()) != ".lock" {
			continue
		}
		lockfile := p.Outdir + "/" + fi.Name()
		outfile := strings.TrimSuffix(lockfile, ".lock")

		fp := p.getFilePathsByOutfile(outfile)

		if isFileBusy(fp.Out) {
			fmt.Printf("ReleaseLocks(): Skip: %q is busy, leaving lockfile\n", fp.Out)
			continue
		}

		fmt.Printf("ReleaseLocks(): Removing lockfile %q\n", fp.OutLock)
		fp.OutLock.Release()
	}

	// Release done locks having no outfile
	files, err = ioutil.ReadDir(p.Outdir)
	if err != nil {
		return err
	}
	for _, fi := range files {
		// Filter done locks
		if filepath.Ext(fi.Name()) != ".done" {
			continue
		}
		lockfile := p.Outdir + "/" + fi.Name()
		outfile := strings.TrimSuffix(lockfile, ".done")

		fp := p.getFilePathsByOutfile(outfile)

		if isFileBusy(fp.Out) {
			continue
		}

		// Check if outfile exists
		_, err := os.Stat(fp.Out)
		if err == nil {
			continue
		}
		fmt.Printf("ReleaseLocks(): Removing done lock %q\n", fp.Done)

		fp.Done.Release()
	}
	return nil
}

// When ffmpeg completed successfully {outfile}.done is created
// Remove non-locked files in Outdir not having .done file
func (p *Pipeline) RemovePartialCompressions() error {
	fps, err := p.getFilePathsByOutdir()
	if err != nil {
		return err
	}

	for _, fp := range fps {
		// Don't touch locked files
		if fp.InLock.IsLocked() || fp.OutLock.IsLocked() {
			continue
		}

		// Check .done lock
		if fp.Done.IsLocked() {
			continue
		}

		// Check if busy
		if isFileBusy(fp.In) || isFileBusy(fp.Out) {
			continue
		}
		// Remove out file
		fmt.Printf("RemovePartialCompressions(): rm %q\n", fp.Out)
		os.Remove(fp.Out)
	}
	return nil
}

// Find & rm files in indir having done lock
// These files are processed successfully and are no longer required
func (p *Pipeline) RemoveProcessedFiles() error {
	fps, err := p.getFilePathsByOutdir()
	if err != nil {
		return err
	}

	for _, fp := range fps {
		// Only process files having {name}.done
		if !fp.Done.IsLocked() {
			continue
		}

		// Check if busy
		if isFileBusy(fp.In) || isFileBusy(fp.Out) {
			continue
		}

		// Verify outfile exists
		outFi, err := os.Stat(fp.Out)
		if err != nil {
			// fp.Out has .done lock but fp.Out does not exist
			fp.Done.Release()
			continue
		}

		// Don't rm infile when outfile is empty
		if outFi.Size() == 0 {
			fmt.Printf("RemoveProcessedFiles(): Skip: outfile %q was empty", fp.Out)
			continue
		}

		// If infile does not exist, continue
		inFi, err := os.Stat(fp.In)
		if err != nil {
			// When the infile does not exist anymore we removed it previously
			continue
		}

		fmt.Printf("RemoveProcessedFiles(): rm %s; sizes: in=%dMB; out=%dMB\n", fp.In, inFi.Size()/1024/1024, outFi.Size()/1024/1024)

		os.Remove(fp.In)
	}

	return nil
}

func (p *Pipeline) RemoveEmptyFiles() error {
	// Check Outdir
	fps, err := p.getFilePathsByOutdir()
	if err != nil {
		return err
	}

	for _, fp := range fps {
		if fp.OutLock.IsLocked() {
			continue
		}

		if isFileBusy(fp.Out) {
			continue
		}

		// Check if empty
		fi, err := os.Stat(fp.Out)
		if err != nil {
			continue
		}

		if fi.Size() != 0 {
			continue
		}

		fmt.Printf("RemoveEmptyFiles(): removing %s\n", fp.Out)
		os.Remove(fp.Out)
	}

	// Check Indir
	fps, err = p.getFilePathsByIndir()
	if err != nil {
		return err
	}

	for _, fp := range fps {
		if fp.InLock.IsLocked() {
			continue
		}

		if isFileBusy(fp.In) {
			continue
		}

		// Check if empty
		fi, err := os.Stat(fp.In)
		if err != nil {
			continue
		}

		if fi.Size() != 0 {
			continue
		}

		fmt.Printf("RemoveEmptyFiles(): removing %s\n", fp.In)
		os.Remove(fp.In)
	}
	return nil
}
