package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/elias-gill/walldo-in-go/globals"
)

// Resize the image to create a thumbnail.
// If a thumbnail already exists just do nothing
func resizeImage(i int) {
	destino := globals.Resized_images[i]
	image := globals.Original_images[i]

	// if the thumnail does not exists
	if _, err := os.Stat(destino); err != nil {
		src, _ := imaging.Open(image)
		src = imaging.Thumbnail(src, 200, 150, imaging.Box)
		// save the thumbnail on a folder
		// TODO  make this folder into .cache or /temp
		imaging.Save(src, destino)
	}
}

// TODO  this need a redesign
// Update the resized_images list
func getResizedImages() {
	var res []string
	sys_os := runtime.GOOS
	path, _ := os.UserHomeDir() // home folder

	// set the path depending on the current OS
	if sys_os == "windows" {
		path += "/AppData/Local/walldo/resized_images/"
	} else {
		// Unix (Mac y Linux)
		path += "/.config/walldo/resized_images/"
	}

	// set a new entry for the resized_images list with a "unique" name
	for _, image := range globals.Original_images {
		destino := path + aislarNombreImagenReescalada(image) + ".jpg"
		res = append(res, destino) // guardar la nueva direccion
	}
	globals.Resized_images = res // guardar la imagenes
}

// Goes trought the configured folders recursivelly and list all the supported image files
func listImagesRecursivelly() {
	// get configured folders from the config file
	globals.Original_images = []string{}
	folders := ConfiguredPaths()

	// loop trought the folder recursivelly
	for _, folder := range folders {
		err := filepath.Walk(folder, func(file string, info os.FileInfo, err error) error {
			if err != nil {
				log.Print(err)
				return err
			}

			// ignore .git files
			if strings.Contains(file, ".git") {
				return filepath.SkipDir
			}
			// TODO  I have a good idea for filters here
			// ignore directories
			if !info.IsDir() && extensionIsValid(file) {
				globals.Original_images = append(globals.Original_images, file)
			}
			return nil
		})
		// TODO  display a dialog error on invalid folders
		if err != nil {
			log.Print(err)
		}
	}
}

// sort images by name
func sortImages(metodo string) {
	// TODO  agregar mas metodos de ordenamiento
	if metodo == "default" {
		sort.Strings(globals.Original_images)
	}
}

// Determine if the file has a valid extension.
// It can be jpg, jpeg and png.
func extensionIsValid(file string) bool {
	// aislar la extension
	aux := strings.Split(file, ".")
	file = aux[len(aux)-1]

	validos := map[string]int{"jpg": 1, "jpeg": 1, "png": 1}
	_, res := validos[file]
	return res
}

// TODO  change the caption size dependending on the grid size
// Returns the first 12 letters of the name of a image. This is for fitting into the captions
func isolateImageName(name string) string {
	// Change backslashes to normal ones
	name = strings.ReplaceAll(name, `\`, `/`)
	res := strings.Split(name, "/")

	largo := len(res) - 1
	aux := res[largo]
	if len(res[largo]) > 12 {
		aux = res[largo][0:12]
		aux = aux + " ..."
	}
	return aux
}

// Returns a new name for the resized image.
// this name has the format parent+file
func aislarNombreImagenReescalada(name string) string {
	name = strings.ReplaceAll(name, `\`, `/`)
	res := strings.Split(name, "/")

	largo := len(res) - 1
	aux := res[largo] + res[largo-1]
	return aux
}
