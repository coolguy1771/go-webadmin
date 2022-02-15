package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	flag.StringVar(&username, "u", "", "Specify username.")
	flag.StringVar(&password, "p", "", "Specify password.")
	flag.StringVar(&path, "path", "", "Specify path for SteamCMD.")
	flag.Parse()
	downloadSteamCMD(path)
	downloadArma(path, username, password)
	startArmaServer()
	stopArmaServer()

}

var path string
var username string
var password string

func downloadSteamCMD(path string) {

	//Change Path to user specified path
	home, _ := os.UserHomeDir()
	err := os.Chdir(filepath.Join(home, path))
	if err != nil {
		log.Println(err)
	}

	//Print out path to confirm correct path
	log.Println("Downloading to " + path)
	workPath, _ := os.Getwd()
	log.Println("Now in " + workPath)

	//OS Detection to download correct SteamCMD Version
	log.Println("Detecting OS")
	osType := detectOS()
	log.Println(strings.ToUpper(osType) + " Detected")

	//Download and install SteamCMD for Windows
	if osType == "windows" {
		url := "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"
		filename := "steamcmd.zip"
		err := DownloadFile(filename, url)
		if err != nil {
			log.Println(err)
		}

		log.Println("Downloaded SteamCMD Windows Version")

		zipfile, err := os.Open(filename)
		if err != nil {
			log.Println(err)
		}

		err = Untar("", zipfile)
		if err != nil {
			log.Println(err)
		}
		log.Println("Installing SteamCMD")

		install := exec.Command("./steamcmd.sh")
		install.Stdin = strings.NewReader("quit")
		var out bytes.Buffer
		install.Stdout = &out
		err = install.Run()
		if err != nil {
			log.Println(err)
		}
		log.Println("SteamCMD Installed")
	}

	//Download and install SteamCMD for Linux
	if osType == "linux" {
		url := "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz"
		filename := "steamcmd_linux.tar.gz"
		err := DownloadFile(filename, url)
		if err != nil {
			log.Println(err)
		}

		log.Println("Downloaded SteamCMD Linux Version")

		gzipfile, err := os.Open(filename)
		if err != nil {
			log.Println(err)
		}

		err = Untar("", gzipfile)
		if err != nil {
			log.Println(err)
		}
		log.Println("Installing SteamCMD")

		install := exec.Command("./steamcmd.sh")
		install.Stdin = strings.NewReader("quit")
		var out bytes.Buffer
		install.Stdout = &out
		err = install.Run()
		if err != nil {
			log.Println(err)
		}
		log.Println("SteamCMD Installed")
	}

	//Download and install SteamCMD for MacOS
	if osType == "darwin" {
		url := "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_osx.tar.gz"
		filename := "steamcmd_osx.tar.gz"
		err := DownloadFile(filename, url)
		if err != nil {
			log.Println(err)
		}

		log.Println("Downloaded SteamCMD MacOS Version")

		gzipfile, err := os.Open(filename)
		if err != nil {
			log.Println(err)
		}

		err = Untar("", gzipfile)
		if err != nil {
			log.Println(err)
		}
		log.Println("Installing SteamCMD")

		install := exec.Command("./steamcmd.sh")
		install.Stdin = strings.NewReader("quit")
		var out bytes.Buffer
		install.Stdout = &out
		err = install.Run()
		if err != nil {
			log.Println(err)
		}
		log.Println("SteamCMD Installed")
	}
}

func downloadArma(path string, username string, password string) {

	//Change Workdir to Download Path
	home, _ := os.UserHomeDir()
	err := os.Chdir(filepath.Join(home, path))
	if err != nil {
		log.Println(err)
	}
	//Change dir to SteamCMD path
	log.Println("Arma 3 Server Installing")

	installArma := exec.Command("./steamcmd.sh")
	installArma.Stdin = strings.NewReader("+force_install_dir ./arma3 +login " + username + password + "+app_update 233780 validate +quit")
	var out bytes.Buffer
	installArma.Stdout = &out
	err = installArma.Run()
	if err != nil {
		log.Println(err)
	}
	log.Println("Arma 3 Server Installed")

}

func startArmaServer() {

	//Start Arma 3 Server
	home, _ := os.UserHomeDir()
	err := os.Chdir(filepath.Join(home, path+"/arma3"))
	if err != nil {
		log.Println(err)
	}
	log.Println("Starting Arma 3 Server")
	installArma := exec.Command("./arma3server_x64")
	var out = bytes.Buffer{}
	installArma.Stdout = &out
	err = installArma.Run()
	if err != nil {
		log.Println(err)
	}
	log.Println("Arma 3 Server Running")
}

func stopArmaServer() {

	//Stop Arma 3 Server
}

func detectOS() string {

	osType := runtime.GOOS

	return osType
}

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func Untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0750); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
