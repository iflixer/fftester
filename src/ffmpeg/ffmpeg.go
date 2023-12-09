package ffmpeg

import (
	"encoding/json"
	"fftester/executor"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Service struct {
}

func (s *Service) Probe(source string) (probe string, probeScore int, timeMs int, err error) {

	type probeJson struct {
		Format struct {
			Filename       string `json:"filename"`
			NbStreams      int    `json:"nb_streams"`
			NbPrograms     int    `json:"nb_programs"`
			FormatName     string `json:"format_name"`
			FormatLongName string `json:"format_long_name"`
			StartTime      string `json:"start_time"`
			Duration       string `json:"duration"`
			Size           string `json:"size"`
			BitRate        string `json:"bit_rate"`
			ProbeScore     int    `json:"probe_score"`
			Tags           struct {
				MajorBrand       string    `json:"major_brand"`
				MinorVersion     string    `json:"minor_version"`
				CompatibleBrands string    `json:"compatible_brands"`
				CreationTime     time.Time `json:"creation_time"`
			} `json:"tags"`
		} `json:"format"`
	}

	args := fmt.Sprintf("-v error -print_format json -show_format -show_streams %s", source)
	if o, e, err := executor.Run("ffprobe", args); err != nil {
		return e.String(), 0, 0, err
	} else {
		str := o.String()
		f := &probeJson{}
		if err := json.Unmarshal([]byte(str), f); err != nil {
			return str, 0, 0, err
		}
		duration, _ := strconv.Atoi(strings.ReplaceAll(f.Format.Duration, ".", ""))
		return str, f.Format.ProbeScore, duration, nil
	}
}

func (s *Service) Convert(taskId int, source, target, pass1in, pass2in string) (pass1out, pass2out, pass1log, pass2log string, err error) {
	log.Printf("convert %s to %s\n", source, target)
	// tpl1 := `-i %s -movflags +faststart -progress http://127.0.0.1:38021/progress -y -vcodec libx264 -filter:v scale=-2:%d -b:v %dM -acodec aac -pix_fmt yuv420p -preset slow -crf 22 -maxrate %dM -bufsize %dM -pass 1 -f %s /dev/null`
	// tpl2 := `-i %s -movflags +faststart -progress http://127.0.0.1:38021/progress -y -vcodec libx264 -filter:v scale=-2:%d -b:v %dM -acodec aac -pix_fmt yuv420p -preset slow -crf 22 -maxrate %dM -bufsize %dM -pass 2 -f %s %s`

	// ffmpeg -i filename.mp4 -codec: copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls filename.m3u8

	if pass1in != "" {
		log.Println("pass1:", pass1in)
		progressKey := fmt.Sprintf(`-progress http://127.0.0.1:38021/progress/1/%d`, taskId)
		pass1out = pass1in
		pass1out = strings.ReplaceAll(pass1out, "[PROGRESS_KEY]", progressKey)
		pass1out = strings.ReplaceAll(pass1out, "[SOURCE]", source)
		if _, e, err1 := executor.Run("ffmpeg", pass1out); err1 != nil {
			return pass1out, pass2out, e.String(), "", err1
		}
	}

	if pass2in != "" {
		log.Println("pass2:", pass2in)
		progressKey := fmt.Sprintf(`-progress http://127.0.0.1:38021/progress/2/%d`, taskId)
		pass2out = pass2in
		pass2out = strings.ReplaceAll(pass2out, "[PROGRESS_KEY]", progressKey)
		pass2out = strings.ReplaceAll(pass2out, "[SOURCE]", source)
		pass2out = strings.ReplaceAll(pass2out, "[TARGET]", target)
		if _, e, err1 := executor.Run("ffmpeg", pass2out); err1 != nil {
			return pass1out, pass2out, "", e.String(), err1
		}
	}
	fmt.Printf("convert %s to %s DONE\n", source, target)
	return
}

/*
	func (f *Ffmpeg) convert(cmd string) error {
		out, errout, err := shell.Shellout(cmd)
		if err != nil {
			log.Printf("convert error: %s, errour: %s\n", err, errout)
			return err
		}
		f.writeLogs(out, errout)
		return nil
	}

	func (f *Ffmpeg) Probe(full bool) (string, error) {
		cmd := f.cmdProbe(full)
		log.Println("Ffmpeg.probe ", cmd)
		stdout, stderr, err := shell.Shellout(cmd)
		if err != nil {
			res := fmt.Sprintf("probe error: %s\n %s\n %s", stdout, stderr, err)
			log.Println(res)
			return res, err
		}
		return fmt.Sprintf("%s\n%s", stdout, stderr), nil
	}

	func (f *Ffmpeg) writeLogs(out, errout string) error {
		fmt.Println("--- stderr ---")
		fmt.Println(errout)
		if err := os.WriteFile("/files/"+f.FileNameLog, []byte(errout), 0644); err != nil {
			log.Printf("error saving errlog: %v\n", err)
			return err
		}
		return nil
	}

	func (f *Ffmpeg) cmdFfmpeg(quality string) string {
		size := ""
		fileResult := f.FileName
		switch quality {
		case "":
			size = "1280:720"
		case "sd":
			size = "512:288"
			fileResult += "_sd"
		case "hd":
			size = "1920:1080"
			fileResult += "_hd"
		}

		p := `ffmpeg \
		-y \
		-hide_banner \
		-fflags +discardcorrupt
		-i %s \
		-preset medium \
		-movflags faststart \
		-c:v libx264 \
		-vf scale=%s \
		-c:a aac \
		-f mp4 \
		%s.mp4`
		return fmt.Sprintf(p, "/files/"+f.FileName, size, "/files/"+fileResult)
	}

	func (f *Ffmpeg) cmdProbe(full bool) string {
		if full {
			return f.cmdProbeFull()
		}
		p := `ffprobe \
		-hide_banner \
		-v error \
		%s`
		return fmt.Sprintf(p, "/files/"+f.FileName)
	}

	func (f *Ffmpeg) cmdProbeFull() string {
		p := `ffprobe \
		-hide_banner \
		%s`
		return fmt.Sprintf(p, "/files/"+f.FileName)
	}
*/
func NewService() (s *Service, err error) {
	// res := strings.Replace(fileName, ".", "_out.", 1)
	s = &Service{}
	// test ffmpeg

	return
}
