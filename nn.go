package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	_ "image/jpeg"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/wcharczuk/go-chart"
)

var input [1024]float64
var wih [1024][16]float64
var dwih [1024][16]float64
var biasi [16]float64
var hidden [16]float64
var who [16][3]float64
var dwho [16][3]float64
var biash [3]float64
var output [3]float64
var target [5][3]float64
var errn [3]float64
var errcy [50]float64
var arrmse [50]float64
var lerate float64
var mseMin float64
var iterMax int = 50
var itterCo int
var imgList = list.New()
var pxbuff, _ = gdk.PixbufNewFromFileAtScale("./assets/img/chart.png", 400, 400, true)

// var imgpx, _ = gtk.ImageNew()

func actfun(value float64) float64 {
	var aktivasi float64
	if value > 0 {
		aktivasi = 0.0
	} else {
		aktivasi = 1.0
	}
	// aktivasi = (1 / (1 - math.Exp(value)))
	// return aktivasi

	return aktivasi

}
func genrand() float64 {
	rand.Seed(time.Now().UnixNano())
	a := (rand.Float64() * 2) - 1
	return a
}
func MainBox() *gtk.Box {
	MainBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	return MainBox
}
func get_buffer_from_tview(tv *gtk.TextView) *gtk.TextBuffer {
	buffer, err := tv.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	return buffer
}
func get_text_from_tview(tv *gtk.TextView) string {
	buffer := get_buffer_from_tview(tv)
	start, end := buffer.GetBounds()

	text, err := buffer.GetText(start, end, true)
	if err != nil {
		log.Fatal("Unable to get text:", err)
	}
	return text
}
func set_text_in_tview(tv *gtk.TextView, text string) {
	buffer := get_buffer_from_tview(tv)
	buffer.SetText(text)
}
func containimg() (*gtk.Grid, *gdk.Pixbuf, *gtk.Image) {
	//pixbuf for image
	containb, _ := gtk.GridNew()
	containimg, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	containimg.SetOrientation(gtk.ORIENTATION_VERTICAL)
	limga := pxbuff
	tv1, _ := gtk.TextViewNew()
	tv2, _ := gtk.TextViewNew()
	tv3, _ := gtk.TextViewNew()

	label1, _ := gtk.LabelNew("Learning Rate")
	label2, _ := gtk.LabelNew("Min MSE")
	label3, _ := gtk.LabelNew("Iter Count")

	//image from pixbuf
	limg, _ := gtk.ImageNew()
	startButton, _ := gtk.ButtonNewWithLabel("start learn")
	conti, _ := gtk.ButtonNewWithLabel("continue learn")
	saveb, _ := gtk.ButtonNewWithLabel("save w")
	loadb, _ := gtk.ButtonNewWithLabel("load w")

	set_text_in_tview(tv1, "0.1")
	set_text_in_tview(tv2, "0.01")
	set_text_in_tview(tv3, "0")

	saveb.Connect("clicked", func() {
		saveWeightIH()
		saveWeightHO()
	})
	loadb.Connect("clicked", func() {
		loadWeightIH()
		loadWeightHO()
	})

	conti.Connect("clicked", func() {
		letsIter()
		pxbuff, _ = gdk.PixbufNewFromFileAtScale("./assets/img/chart.png", 400, 400, true)
		limg.SetFromPixbuf(pxbuff)
		itterCo = itterCo + 50
		set_text_in_tview(tv3, strconv.Itoa(itterCo))
	})

	startButton.Connect("clicked", func() {
		lr, _ := strconv.ParseFloat(get_text_from_tview(tv1), 64)
		mse, _ := strconv.ParseFloat(get_text_from_tview(tv2), 64)
		// miniter, _ := strconv.Atoi(get_text_from_tview(tv3))
		lerate = lr
		mseMin = mse
		// iterMax = miniter
		fmt.Println(lerate)
		fmt.Println(mseMin)
		fmt.Println(iterMax)
		//randoming weight
		for i := 0; i < len(wih); i++ {
			for j := 0; j < len(wih[i]); j++ {
				wih[i][j] = genrand()
			}
		}
		for i := 0; i < len(who); i++ {
			for j := 0; j < len(who[i]); j++ {
				who[i][j] = genrand()
			}
		}
		for i := 0; i < len(biasi); i++ {
			biasi[i] = genrand()
		}
		for i := 0; i < len(biash); i++ {
			biash[i] = genrand()
		}
		letsIter()
		pxbuff, _ = gdk.PixbufNewFromFileAtScale("./assets/img/chart.png", 400, 400, true)
		limg.SetFromPixbuf(pxbuff)

		itterCo = 50
		set_text_in_tview(tv3, strconv.Itoa(itterCo))

	})

	imgList.PushBack(limg)
	containb.Attach(saveb, 0, 0, 1, 1)
	containb.Attach(loadb, 0, 1, 1, 1)

	containimg.Attach(startButton, 0, 0, 1, 1)
	containimg.Attach(conti, 1, 0, 1, 1)
	containimg.Attach(containb, 2, 0, 1, 1)
	containimg.Attach(label1, 0, 1, 1, 1)
	containimg.Attach(label2, 1, 1, 1, 1)
	containimg.Attach(label3, 2, 1, 1, 1)
	containimg.Attach(tv1, 0, 2, 1, 1)
	containimg.Attach(tv2, 1, 2, 1, 1)
	containimg.Attach(tv3, 2, 2, 1, 1)
	containimg.Attach(limg, 0, 3, 3, 1)

	// miniter, _ := strconv.Atoi(get_text_from_tview(tv3))
	lr, _ := strconv.ParseFloat(get_text_from_tview(tv1), 64)
	mse, _ := strconv.ParseFloat(get_text_from_tview(tv2), 64)
	// fmt.Println(miniter)
	fmt.Println(lr)
	fmt.Println(mse)
	// iterMax = miniter
	mseMin = mse
	lerate = lr

	return containimg, limga, limg
}
func BuildUI() *gtk.Box {

	MainBox := MainBox()
	// Add the label to the window.
	// MainBox.Add(l)
	// MainBox.Add(m)
	// MainBox.Add(limg)
	containimg, _, _ := containimg()
	MainBox.Add(containimg)
	return MainBox

}
func Graphs() {
	var loops [50]float64
	for i := 0; i < len(loops); i++ {
		loops[i] = float64(i)
	}
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(9).WithAlpha(64),
					StrokeWidth: 8,
				},
				XValues: loops[:],
				YValues: arrmse[:],
			},
		},
	}

	f, _ := os.Create("./assets/img/chart.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
func spagethynn(idx int, key string, nkey int) {
	fName := "sample/" + key + strconv.Itoa(idx) + ".jpg"
	f, err := os.Open(fName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// http://www.dcode.fr/binary-image

	// store input value
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			bin := 1
			if float64((r+g+b))/65535/3 > 0.7 {
				bin = 0
			}
			input[(y*32)+x] = float64(bin) //stored at array
			// fmt.Print(int(input[(y*32)+x]))
		}
		// fmt.Println("")
	}

	//
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(hidden); j++ {
			hidden[j] = hidden[j] + (input[i] * wih[i][j])
		}
	}

	for i := 0; i < len(biasi); i++ {
		hidden[i] = hidden[i] + biasi[i]
		hidden[i] = actfun(hidden[i])
	}

	for i := 0; i < len(hidden); i++ {
		for j := 0; j < len(output); j++ {
			output[j] = output[j] + (hidden[i] * who[i][j])
		}
	}

	//test
	var error1 float64
	var err3 float64

	for i := 0; i < len(biash); i++ {
		output[i] = output[i] + biash[i]
		output[i] = actfun(output[i])
		error1 = math.Abs(output[i] - target[nkey][i])
		errn[i] = error1
		// fmt.Print(output[i], " ", target[nkey][i], "   ")
	}
	// fmt.Println("    ")
	for i := 0; i < len(output); i++ {
		err3 = err3 + (0.5 * math.Pow(errn[i], 2))

	}
	// fmt.Println(err3)
	//store error for MSE
	errcy[idx-1+(10*nkey)] = err3

	// calculate delta of all weight
	for i := 0; i < len(hidden); i++ {
		for j := 0; j < len(output); j++ {
			dwho[i][j] = (output[j] - target[nkey][j]) * 1 * hidden[i] * lerate
		}
	}
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(hidden); j++ {
			var sumerr float64
			for k := 0; k < len(output); k++ {
				sumerr = sumerr + ((output[k] - target[nkey][k]) * 1 * who[j][k])
			}
			dwih[i][j] = sumerr * 1 * input[i] * lerate
		}
	}
	//adding delta to wighth
	for i := 0; i < len(hidden); i++ {
		for j := 0; j < len(output); j++ {
			who[i][j] = who[i][j] + dwho[i][j]
		}
	}
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(hidden); j++ {
			wih[i][j] = wih[i][j] + dwih[i][j]
		}
	}

}
func phaseLoop() {
	keysheet := [5]string{"a", "b", "c", "d", "e"}
	for nkey, value := range keysheet {
		for i := 1; i <= 10; i++ {
			spagethynn(i, value, nkey)

		}
	}
}
func updateImg() {

	pixb, _ := gdk.PixbufNewFromFileAtScale("./assets/img/chart.png", 400, 400, true)
	imgtk, _ := gtk.ImageNew()
	imgtk.SetFromPixbuf(pixb)
	imgtk.ShowAll()
	// box.Add(imgtk)
}
func letsIter() {
	for i := 0; i < iterMax; i++ {
		var me float64
		var mse float64
		phaseLoop()
		for j := 0; j < len(errcy); j++ {
			me = me + errcy[j]
		}
		mse = math.Pow((me / float64(len(errcy))), 2)
		arrmse[i] = mse
		if mse < mseMin {
			for k := i; k < iterMax; k++ {
				arrmse[k] = 0.0
			}
			Graphs()
			updateImg()
			fmt.Println(arrmse)
			break
		}

	}
	Graphs()
	updateImg()
	fmt.Println(arrmse)
}
func saveWeightIH() {
	var txt bytes.Buffer
	file, err := os.Create("./assets/weight/bobotih.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	for i := 0; i < len(wih); i++ {
		for j := 0; j < len(wih[i]); j++ {
			s := strconv.FormatFloat(wih[i][j], 'E', -1, 64)
			txt.WriteString(s)
			if j == len(wih[i])-1 {
				continue
			}
			txt.WriteString(",")
		}
		txt.WriteString("\n")
	}
	file.WriteString(txt.String())
	// fmt.Println(txt.String())
}
func loadWeightIH() {

	var a string
	file, err := os.Open("./assets/weight/bobotih.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		a = scanner.Text()
		splita := strings.SplitAfter(a, ",")
		for j := 0; j < len(splita); j++ {
			num, _ := strconv.ParseFloat(splita[j], 64)
			wih[i][j] = num
		}
		i++
	}

}
func saveWeightHO() {
	var txt bytes.Buffer
	// os.Remove("test.txt")
	file, err := os.Create("./assets/weight/bobotho.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	for i := 0; i < len(who); i++ {
		for j := 0; j < len(who[i]); j++ {
			s := strconv.FormatFloat(who[i][j], 'E', -1, 64)
			txt.WriteString(s)
			if j == len(who[i])-1 {
				continue
			}
			txt.WriteString(",")
		}
		txt.WriteString("\n")
	}
	file.WriteString(txt.String())
	// fmt.Println(txt.String())
}
func loadWeightHO() {

	var a string
	file, err := os.Open("./assets/weight/bobotho.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		a = scanner.Text()
		splita := strings.SplitAfter(a, ",")
		for j := 0; j < len(splita); j++ {
			num, _ := strconv.ParseFloat(splita[j], 64)
			who[i][j] = num
		}
		i++
	}

}
func main() {
	targetf := [5][3]float64{
		{0.0, 1.0, 0.0},
		{0.0, 1.0, 1.0},
		{1.0, 0.0, 0.0},
		{1.0, 0.0, 1.0},
		{1.0, 1.0, 0.0},
	}
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("freakin noodle (nn-learn-handwrite)")
	// win.SetDefaultSize(800, 600)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	// BuildUI()
	MainBox := BuildUI()
	win.Add(MainBox)
	target = targetf
	win.ShowAll()

	// MainBox.ShowAll()

	gtk.Main()
}
