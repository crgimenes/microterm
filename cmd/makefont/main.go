package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fontFile := ""

	flag.StringVar(&fontFile, "font", "", "file with the font model")
	flag.Parse()

	if fontFile == "" {
		fmt.Println("Font file is required")
		return
	}

	// Read fonte file
	b, err := os.ReadFile(fontFile)
	if err != nil {
		fmt.Println("Error reading font file", err)
		return
	}

	fontFile = filepath.Base(fontFile)
	fontName := fontFile[:len(fontFile)-len(filepath.Ext(fontFile))]

	// Create finename.cpp with fontfile
	cFile := fmt.Sprintf("%s.cpp", fontName)

	header := fmt.Sprintf(`static const unsigned char %s[] = {
	`, fontName)

	footer := `
};
`
	f, err := os.Create(cFile)
	if err != nil {
		fmt.Println("Error creating file", err)
		return
	}

	defer f.Close()

	_, err = f.WriteString(header)
	if err != nil {
		fmt.Println("Error writing header", err)
		return
	}

	/*

	   font file 4x6
	   "*** "
	   "* * "
	   "*** "
	   "* * "
	   "* * "
	   ""
	   "*** "
	   "* * "
	   "**  "
	   "* * "
	   "*** "
	   ""
	*/

	col := 0
	aux := 0
	c := 0
	fmtAux := "0x%02X"
	for _, v := range b {
		switch v {
		case '*':
			aux = aux | (1 << (7 - col))
		case ' ':
			aux = aux | (0 << (7 - col))
		case '\n':
			continue
		}
		col++
		if col == 8 {
			col = 0
			_, err = f.WriteString(fmt.Sprintf(fmtAux, aux))
			if err != nil {
				fmt.Println("Error writing data", err)
				return
			}
			fmtAux = ", 0x%02X"
			c++
			if c > 8 {
				fmtAux = ",\n\t0x%02X"
				c = 0
			}
			aux = 0
		}
	}

	_, err = f.WriteString(footer)
	if err != nil {
		fmt.Println("Error writing footer", err)
		return
	}
}
