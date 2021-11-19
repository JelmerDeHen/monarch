
```
ffmpeg -f v4l2 -framerate 60 -video_size 1920x1080 -input_format mjpeg -i /dev/video0 -c copy mjpeg.mkv

ffmpeg -f v4l2 -framerate 60 -video_size 1920x1080 -input_format mjpeg -i /dev/video0 -preset faster -pix_fmt yuv420p mjpeg.mkv





arecord -D sysdefault:CARD=NTUSB -t wav -f S24_3LE -r 192000 out.wav
ffmpeg -f v4l2 -framerate 60 -video_size 1920x1080 -input_format mjpeg -i /dev/video0 -preset faster -pix_fmt yuv420p mjpeg.mkv
```
