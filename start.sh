#!/bin/bash

function start() {
  systemctl --user restart monarch_arecord
  systemctl --user restart monarch_v4l2
  systemctl --user restart monarch_x11grab
}

function stop() {
  systemctl --user stop monarch_arecord
  systemctl --user stop monarch_v4l2
  systemctl --user stop monarch_x11grab
}

function deploy() {
  cp -v sd/* "${HOME}/.config/systemd/user/"
  systemctl --user daemon-reload
}

