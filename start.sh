#!/bin/bash

systemctl --user restart monarch_arecord
systemctl --user restart monarch_v4l2
systemctl --user restart monarch_x11grab
