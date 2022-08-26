
build:
				mkdir -pv bin
				go build -o bin/monarch ./core/

clean:
				rm -rf bin

install:
				cp -v sd/monarch_* "${HOME}/.config/systemd/user"
				systemctl --user enable monarch_arecord.timer
				systemctl --user enable monarch_compress.timer
				systemctl --user enable monarch_v4l2.timer
				systemctl --user enable monarch_x11grab.timer
				systemctl --user daemon-reload

uninstall:
				systemctl --user stop monarch_*
				systemctl --user disable monarch_arecord.timer
				systemctl --user disable monarch_compress.timer
				systemctl --user disable monarch_v4l2.timer
				systemctl --user disable monarch_x11grab.timer
				rm -v "${HOME}/.config/systemd/user/monarch_"*
				systemctl --user daemon-reload

start:
				systemctl --user start monarch_arecord.timer
				systemctl --user start monarch_v4l2.timer
				systemctl --user start monarch_x11grab.timer
				systemctl --user start monarch_compress.timer

stop:
				systemctl --user stop monarch_*

status:
				systemctl status --user monarch_*
