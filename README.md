# GO AUDIO LIBRARY

## extremely simple and easy to use 

### can play MP3 files
### can capture audio from the mic to a WAV file

## Examples

### Audio Capture

```golang
 
r := &audio.Recorder{}

r.Init("./test.wav");
err := r.Start();
if err != nil {
    fmt.Println(err);
    return;
}

fmt.Println("Press ENTER to stop recording");

fmt.Scanln();

fmt.Println("done recording!");

err = r.Stop();
if err != nil {
    fmt.Println(err);
    return;
}
```
