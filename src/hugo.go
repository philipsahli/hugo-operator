package hugo

import (
	// "os"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type Repository struct {
	URL       string
	Name      string
	Directory string
}

type Branch struct {
	Name       string
	Hash       string
	Repository *Repository
}

func build(dir string) (string, error) {

	cmd := exec.Cmd{
		Path: "/usr/local/bin/hugo",
		Dir:  dir,
	}
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// fmt.Printf("The date is %s\n", output)
	return string(output), nil
}

func (r *Repository) buildAndUpload() error {
	branches, err := r.GetBranches()
	if err != nil {
		return err
	}
	for _, b := range branches {
		b.buildAndUpload()
	}
	return nil
}

func (b *Branch) buildAndUpload() (string, error) {
	_, err := b.build()
	if err != nil {
		return "", err
	}
	bucketName, err := b.upload()
	if err != nil {
		return "", err
	}
	return bucketName, nil
}

func (b *Branch) build() (string, error) {
	dir, err := b.checkout()

	if err != nil {
		return "", err
	}

	// defer os.RemoveAll(dir)

	b.Repository.Directory = fmt.Sprintf("%s/%s", dir, "examples/blog")
	output, err := build(b.Repository.Directory)
	if err != nil {
		return "", err
	}
	fmt.Println(output)

	return output, nil
}

func (b *Branch) upload() (string, error) {
	awsProfileName := "default"
	// bucketName := fmt.Sprintf("%s-%s", "hugo-operator", b.Name)
	bucketName := fmt.Sprintf("%s-%s", "hugo-operator", b.Repository.Name)
	directory := fmt.Sprintf("%s/%s", b.Repository.Directory, "public")
	prefix := b.Name
	sess := makeSession(awsProfileName)
	// sess *session.Session, bucketName string, bucketPrefix string, dirPath string) (string, error) {
	bucketName, err := uploadDirToS3(sess, bucketName, prefix, directory)
	if err != nil {
		return "", err
	}
	return bucketName, nil
}

func (b *Branch) checkout() (string, error) {
	dir, err := ioutil.TempDir("/tmp", "prefix")
	b.Repository.Directory = dir
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	r, _ := git.PlainClone(dir, false, &git.CloneOptions{
		URL: b.Repository.URL,
	})
	w, _ := r.Worktree()

	err = r.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		fmt.Println(err)
		return dir, err
	}

	bN := plumbing.NewBranchReferenceName(b.Name)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: bN,
		Force:  true,
	})
	if err != nil {
		fmt.Printf("Error on Checkout() ref %s: %s", bN, err)
		return dir, err
	}

	return dir, err
}

func (r *Repository) GetBranches() ([]Branch, error) {

	bList := make([]Branch, 0)
	bListPtr := &bList

	gr, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: r.URL,
	})
	opts := &git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	}

	if err := gr.Fetch(opts); err != nil {
		return bList, err
	}
	if err != nil {
		return bList, fmt.Errorf("Error in Fetch((): %s", err)
	}
	bIter, err := gr.Branches()
	if err != nil {
		return bList, fmt.Errorf("Error in Branches(): %s", err)
	}
	bIter.ForEach(func(c *plumbing.Reference) error {
		fmt.Println(c)

		n := c.Name()
		b := Branch{
			Name:       n.String(),
			Hash:       c.Hash().String(),
			Repository: r,
		}
		*bListPtr = append(*bListPtr, b)
		return nil
	})
	return *bListPtr, nil
}
