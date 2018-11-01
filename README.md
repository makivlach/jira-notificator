Jira notificator
=================

![logo](doc/logo.png)

An notification application for cloud based Jira workspace. Notifications are directly fired inside the desktop environment with a custom audio tone. Especially useful in case of disabled email notifications.

The notificator consists of CLI and GUI version. Both are cross-platform for Linux, Mac and Windows.

_Note: the Linux version might require additional libraries to install. More about it down below._

# Installation

Follow the instructions down below.

## Windows + Mac:

Just [download](https://github.com/vlachmilan/jira-notificator/releases) the latest version  and open the app. There should be no other requirement.

## Linux

_The installation process described below has been designed for Ubuntu based distributions only. Installation for other distros might differ - if so, **please open new [issue](https://github.com/vlachmilan/jira-notificator/issues)**._

1. `sudo apt install libasound2-dev`
2. `sudo apt install libgtk-3-dev`
3. `sudo apt install notify-osd`
4. [download](https://github.com/vlachmilan/jira-notificator/releases) and run the app 

# Screenshots

### UI:
![screenshot1](doc/screenshot1.png)

# Known issues & todos
- [ ] "crackly" tone issues
- [ ] ugly GUI
- [ ] ability to minimize the app to system tray
- [ ] code cleanup and documentation
- [ ] better GUI code organization

# Licence

MIT