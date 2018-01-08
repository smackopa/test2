package main

import 	("fmt"
	"strings"
	"strconv"
	"time"
	"serial"
//	"os"
	"log"
	)


var pubx_00 string ="$PUBX,00,081350.00,4717.113210,N,00833.915187,E,546.589,G3,2.1,2.0,100.0,77.52,0.007,,0.92,1.19,0.77,9,0,0*5F"
var gprmc string ="$GPRMC,083559.00,A,4717.11437,N,00833.91522,E,0.004,77.52,091203,,,A*57"
var pubx_04 string ="$PUBX,04,073731.00,091202,113851.00,1196,15D,1930035,-2660.664,43,*3C"
var gpsTime time.Time;
var p = fmt.Println


//$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A


func setNMEAChecksum(s string) (string) {
	// strip leading "$" and split at "*"
	s_split := strings.Split(strings.TrimPrefix(s, "$"), "*")
	s_out := s_split[0]

	cs_calc := byte(0)
	for i := range s_out {
		cs_calc = cs_calc ^ byte(s_out[i])
	}

	cs_s :=fmt.Sprintf("%2X",cs_calc)

	
	return "$"+s_out+"*"+cs_s
}


func initSerial() (*serial.Port ,bool) {

/*	baudrate := int(4800)
	device := os.Getenv("COM1")
	if device == "" {
		log.Printf("Skipping test because COM1 environment variable is not set")
	}	// Open port at default baud for config.*/

	serialConfig := &serial.Config{Name: "COM1", Baud: 4800}
	s, err := serial.OpenPort(serialConfig)
        if err != nil {
                log.Fatal(err)
		return nil, true
        }

	return s, false
}

func sendNMEAserial(l string)(bool){
	var gprmc string

	s,err := initSerial()
	if !err {
	
	x := strings.Split(l, ",")


	if x[0] == "$PUBX" {
		if x[1] == "00" {
			// Speed
			gs, _ := strconv.ParseFloat(x[11], 32)
			gs = gs * 0.540003 // convert to knots
			gs_s :=fmt.Sprintf("%05.1f",gs) 
			// track
			track, _ := strconv.ParseFloat(x[12], 32)
			track_s :=fmt.Sprintf("%05.1f",track)
			// date
			date_s:=gpsTime.Format("020106")

			gprmc="$GPRMC,"+x[2][0:6]+",A,"+x[3]+","+x[4]+","+x[5]+","+x[6]+","+gs_s+","+track_s+","+date_s+",,*XX"
			gprmc = setNMEAChecksum(gprmc)
		} else if x[1] == "04" {
			gpsTime, _ = time.Parse("020106", x[3])
			return true
		}
	} else if (x[0] == "$GPRMC") {
		gpsTime, _ = time.Parse("020106", x[9])
		return true
	}	
	

        n,_ := s.Write([]byte(gprmc))
        
        p("Sended bytes%d", n)
	s.Close()

	p(gprmc)	
	return false
	} else {
		return true
	}

}

func main() {
	sendNMEAserial(pubx_04)
	sendNMEAserial(pubx_00)
	sendNMEAserial(gprmc)
	sendNMEAserial(pubx_00)
}