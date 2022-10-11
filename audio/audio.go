package audio

import (
	"fmt"
	"log"
	"os"

	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	"github.com/gen2brain/malgo"
	wav "github.com/youpy/go-wav"
)


func PlayMp3(fileName string){
	file , err := os.Open(fileName);
	if err != nil {
		log.Fatal(err);
	}
	defer file.Close();

	streamer , format , err := mp3.Decode(file);
	if err != nil {
		log.Fatal(err);
	}
	defer streamer.Close();
	
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10));

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}


type Recorder struct {
	IsRecording 		bool
	FileName 			string

	ctx 				*malgo.AllocatedContext
	device 				*malgo.Device

	Samples 			[]byte
	SampleCount 		uint32

	SampleRate 			uint32 // can modify this to your liking
	Channels 			uint32 // can modify
	Format 				malgo.FormatType // can modify

	ByteSize 			uint32
	BitSize 			uint32
}

func (r *Recorder) Init(fileName string){
	r.Samples = make([]byte , 0);

	r.FileName = fileName;
	r.SampleRate = 16000;
	r.Channels = 1;
	r.Format = malgo.FormatS16;
}

func (r *Recorder) Start() error {
	var err error
	r.ctx, err = malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		return err;
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = r.Format
	deviceConfig.Capture.Channels = r.Channels
	deviceConfig.SampleRate = r.SampleRate
	deviceConfig.Alsa.NoMMap = 1



	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))
	r.ByteSize = sizeInBytes;
	r.BitSize = 8 * sizeInBytes;
	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {

		sampleCount := framecount * deviceConfig.Capture.Channels * sizeInBytes

		newCapturedSampleCount := r.SampleCount + sampleCount

		r.Samples = append(r.Samples , pSample...)

		r.SampleCount = newCapturedSampleCount

	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}
	device, err := malgo.InitDevice(r.ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		return err;
	}

	r.device = device;

	err = device.Start()
	if err != nil {
		return err;
	}

	return nil;
}

func (r *Recorder) Stop() error {

	r.device.Uninit();
	_ = r.ctx.Uninit();
	r.ctx.Free();

	file , err := os.Create(r.FileName);
	if err != nil {
		return err;
	}
	defer file.Close();

	w := wav.NewWriter(file , r.SampleCount , uint16(r.Channels) , r.SampleRate , uint16(r.BitSize));
	
	_ , err = w.Write(r.Samples);
	if err != nil {return err}

	return nil;
}
