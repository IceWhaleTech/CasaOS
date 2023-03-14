# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed


### Removed

### Security

## [0.4.3]

### Added

- [Disk] Now usb also supports merging to


### Changed

- [File] Solve the installation dependency problem, make the installation more smoothly
- [File] Change the default permissions of the sharing folder

### Fixed

- [System] Fixed  not see wlan iface ([#909](https://github.com/IceWhaleTech/CasaOS/issues/909))
- [System] Terminal font issue fix ([#929](https://github.com/IceWhaleTech/CasaOS/issues/929))
- [File] Fixed the problem of not being able to launch after mounting

### Removed


## [0.4.2]

### Added

- [App] Increase the display of progress during the installation process
- [App] Label whether the current app supports x86 or Pi devices
- [App] Support single app version upgrade
- [File] Support mounting of Google Drive and Dropbox cloud drives
- [System] Support Mint Linux

### Changed

- [File] Optimize the download speed of a single file

### Fixed

- [Share] Fix the samba permission issue 
- [Disk] Fix the problem of disk mount point plus 1 after upgrade ([#770](https://github.com/IceWhaleTech/CasaOS/issues/770))
- [File] Fix the problem of file permission change caused by modifying files in casaos ([#829](https://github.com/IceWhaleTech/CasaOS/issues/829))
- [Share] Fix the problem of files being deleted due to samba uninstallation failure ([#843](https://github.com/IceWhaleTech/CasaOS/issues/843))



## [0.4.1] - 2023-1-19


### Added
- [Disk] Added disk merging feature in storage management (beta) that allows for multiple disks to be merged into a single storage space
- [System] Added option for startpage.com search engine
- [APP] Added app cloning feature in the app's context menu.
### Changed
- [APP] Improved app installation process, including display of the installation process, checks for successful installation, and prompts
- [System] Binary sizes are 40%~60% smaller (thanks to upx)
- [App] Optimization of install and update for certain country.
- [All] Lots of bug fixes

## [0.4.0] - 2022-12-13
### Added

- [Developer] Included `casaos-cli` command tool for debugging
- [Developer] Added message bus for events and actions - Use `casaos-cli message-bus` to manage.
- [Disk] Disk notification in Dashboard
- [System] Restart/shutdown directly from CasaOS Dashboard
### Changed

- [General] CasaOS new logo!
- [App] Redesign of Featured App
- [App] Now you can choose to delete userdata along with app uninstallation

### Security

- [System] Fixed a shell injection issue for better security

### Fixed

- [System] Re-instate default zone0 for CPU Temp ([#694](https://github.com/IceWhaleTech/CasaOS/issues/694))
- [Disk] Fixed storage name with extra `-1` after rebooting ([#698](https://github.com/IceWhaleTech/CasaOS/issues/698))
- [Disk] Fixed disk check so it does not impact disk going into idle ([#704](https://github.com/IceWhaleTech/CasaOS/issues/704))

## [0.3.8] 2022-11-21

### Added
- [System] Add system announcement
- [App] Allow to turn off the display of "Existing Docker Apps" in the settings.

### Changed
- [System] Improve the feedback function, you can submit feedback in the bottom right corner of WebUI.

### Fixed
- [System] Fix CPU Temp for other platforms ([#661](https://github.com/IceWhaleTech/CasaOS/issues/661))

## [0.3.7.1] 2022-11-04

### Fixed

- Fix memory leak issue ([#658](https://github.com/IceWhaleTech/CasaOS/issues/658)[#646](https://github.com/IceWhaleTech/CasaOS/issues/646))
- Solve the problem of local application import failure ([#490](https://github.com/IceWhaleTech/CasaOS/issues/490))

## [0.3.7] 2022-10-28

### Added
- [Storage] Disk merge (Beta), you can merge multiple disks into a single storage space (currently you need to enable this feature from the command line)

### Changed
- [Files] Changed the cache file storage location, now the file upload size is not limited by the system disk capacity.
- [Scripts] Updated installation and upgrade scripts to support more Debian-based Linux distributions.
- [Engineering] Refactored Local Storage into a standalone service as part of CasaOS modularization.

### Fixed
- [Apps] App list update mechanism improved, now you can see the latest apps in App Store immediately.
- [Storage] Fixed a lot of known issues

### Added
- [Storage] Disk merge (Beta), you can merge multiple disks into a single storage space (currently you need to enable this feature from the command line)

### Changed
- [Files] Changed the cache file storage location, now the file upload size is not limited by the system disk capacity.
- [Scripts] Updated installation and upgrade scripts to support more Debian-based Linux distributions.
- [Engineering] Refactored Local Storage into a standalone service as part of CasaOS modularization.

### Fixed
- [Apps] App list update mechanism improved, now you can see the latest apps in App Store immediately.
- [Storage] Fixed a lot of known issues


## [0.3.6] - 2022-09-06

###  Added
- [System] Added power and temperature info to performance widget (Intel)
- [Apps] Custom links can be added to Apps section

### Fixed
- [Apps] Fixed the problem of not being able to modify some App settings ([#510](https://github.com/IceWhaleTech/CasaOS/issues/510))

### Changed
- [System] Architecture optimization. Improved performance.

## [0.3.5] - 2022-08-23

### Added

- [File] Mount the shared samba
- [File] File sharing via Samba
- [System] You can share casaos on Twitter, facebook, reddit

### Changed

- [Disk] Support for mounting existing data disks

### Fixed

- [App] fixed uninstalling imported docker container apps results in wiping ALL your config data from them ([#360](https://github.com/IceWhaleTech/CasaOS/issues/360))

## [0.3.4] - 2022-07-29(UTC)

### Added

- SSH adds port-side options and prompts for connection status. ([#286](https://github.com/IceWhaleTech/CasaOS/issues/286))

### Changed

- Normalize all routes
- Application names now support spaces ([#211](https://github.com/IceWhaleTech/CasaOS/issues/211))

### Removed

- Removed  casaos connect

### Security

- Adjustment of authentication method

### Fixed

- Fixed storage format and remove password error issues ([#344](https://github.com/IceWhaleTech/CasaOS/issues/344) [#357](https://github.com/IceWhaleTech/CasaOS/issues/357))

## [0.3.3] - 2022-07-08(UTC)

### Added

- [System]Add interface call log
- Adding Developing file ([#311](https://github.com/IceWhaleTech/CasaOS/pull/311))
- [App] add new tips for app section.
- [System] UI Configurable function modules: support turning off the search bar and recommended apps module in the settings.
- [System] Custom wallpapers: two new preset wallpapers, support for custom uploads, support for setting images from Files as wallpapers, Also support right click on dashboard to change wallpaper.

### Changed

- [App] Cache app store index and category data
- [System] casaos master program adapted to FHS standards
- [App] Update casaos icons.
- [System] Update translation.

### Removed

- [System] Remove upnp function module
- [System] Remove ddns function module
- [System] Remove search function module
- [System] Remove zerotier function module
- [System] Remove task function module
- [System] Remove file share function module

### Fixed

- [Disk] Fixed hard drive won't hibernate problem ([#202](https://github.com/IceWhaleTech/CasaOS/issues/202))
- [File] Fixed the backspace key that causes the folder to rewind ([#252](https://github.com/IceWhaleTech/CasaOS/issues/252))
- [App] Fixed app logo is not loading when imported. ([#320](https://github.com/IceWhaleTech/CasaOS/issues/320))

## [0.3.2.1] - 2022-06-16(UTC)

### Changed

- [System] Adjusted the display style.

### Fixed

- [System] Fixed the issue of widgets displaying wrongly on mobile devices.
- [App] Fix the problem of application opening failure on non-80 ports ([#283](https://github.com/IceWhaleTech/CasaOS/issues/283) [#280](https://github.com/IceWhaleTech/CasaOS/issues/280))
- [System] Modify port failure problem ([#282](https://github.com/IceWhaleTech/CasaOS/issues/282))
- [App]Modify environment variables disappearing problem([#284](https://github.com/IceWhaleTech/CasaOS/issues/284))
- [System]Fix no update alert([#278](https://github.com/IceWhaleTech/CasaOS/issues/278))
- [System] Fixed some bugs of application cpu usage and memory staging([#272](https://github.com/IceWhaleTech/CasaOS/issues/272))
- [App] Fixed plex and HA network mode error issues ([#299](https://github.com/IceWhaleTech/CasaOS/issues/299))
- [App] Fix application terminal not working ([#266](https://github.com/IceWhaleTech/CasaOS/issues/266))

## [0.3.2] - 2022-06-10

### Added

- [Files] Files can now be selected multiple files and downloaded, deleted, moved, etc.
- [Apps] Support to modify the application opening address.([#204](https://github.com/IceWhaleTech/CasaOS/issues/204))

### Changed

- [Apps] Hide the display of non-essential environment variables in the application.([#196](https://github.com/IceWhaleTech/CasaOS/issues/196))
- [System] Network, disk, cpu, memory, etc. information is modified to be pushed via socket.
- [System] Optimize opening speed.([#214](https://github.com/IceWhaleTech/CasaOS/issues/214))
- [Language] Update language pack [zarevskaya](https://github.com/zarevskaya) [patrickhilker](https://github.com/patrickhilker)
- [System] Interface path adjustment

### Removed

- [Files] Remove the online preview function of PDF files

### Fixed

- [System] Fixed the problem that sync data cannot submit the device ID ([#68](https://github.com/IceWhaleTech/CasaOS/issues/68))
- [Files] Fixed the code editor center alignment display problem.([#210](https://github.com/IceWhaleTech/CasaOS/issues/210))
- [Files] Fixed the problem of wrong name when downloading files.([#240](https://github.com/IceWhaleTech/CasaOS/issues/240))
- [System] Fixed the network display as a negative number problem.([#224](https://github.com/IceWhaleTech/CasaOS/issues/224))
- [System] Fixed the problem of wireless network card traffic display.([#222](https://github.com/IceWhaleTech/CasaOS/issues/222))


## [0.3.1.1] - 2022-05-17

### Fixed

- Fix the data loss problem when importing local applications

## [0.3.1] - 2022-05-16

### Added

- CasaConnect and file add image thumbnail function
- Import of docker applications
- List support custom sorting function
- CasaConnect gives priority to LAN connections
- USB auto-mount switch (Raspberry Pi is off by default)
- Application custom installation supports Docker Compose configuration import in YAML format
- You will see the new version changelog from the next version
- Added live preview for icons in custom installed applications

### Changed

- Application data is no longer saved to the database
- Optimize app store speed issues
- Optimize the way WebUI is filled in
- Image preview has been completely upgraded and now supports switching between all images in the same folder, as well as dragging, zooming, rotating and resetting.
- Added color levels to the CPU and RAM charts
- Optimized the display of the Connect friends list right-click menu
- Change the initial display directory to /DATA

### Removed

- Historical Application Data

### Fixed

- Fixed the problem that some Docker CLI commands failed to import
- Fix the problem that the application is not easily recognized in /DATA/AppData directory and docker command line after installation, it will be shown as application name
- Fix Pi-hole installation failure
- Fixed the issue that the app could not be updated using WatchTower
- Fixed the problem that the task status was lost after closing Files when there was an upload task

## [0.3.0] - 2022-04-08

### Added

- Add CasaConnect function, now you can share private files peer-to-peer with your friends.
- Add a widget for network traffic monitoring.
- 12 new popular apps added to App Center

### Changed

- Updated the sidebar of Files.
- Updated the initial directory of Files to the Root directory.
- Armbian 22.02 armhf/arm64/amd64 platform tests passed [@igorpecovnik ](https://github.com/igorpecovnik)
- Elementary OS 6.1 Jólnir amd64 platform tests passed [@alvarosamudio ](https://github.com/alvarosamudio)

### Fixed

- Fix an issue in Files where the backspace button would trigger a return to the previous level of the directory when creating a folder.
- Fix the display problem of application list in CPU widget.
- Fix the problem that the ipv6 of the application cannot be opened

### Removed

- Interfaces related to "zerotier"

## [0.2.10] - 2022-03-10

### Added

- Added CasaOS own file manager, now you can browse, upload, download files from the system, even edit code online, preview photos and videos through it. It will appear in the first position of Apps.
- Added CPU core count display and memory capacity display.

### Changed

- Optimized the rendering performance of the home page.
- Optimized the internationalization display of the time widget.
- Show the icon of the stopped application as gray.
- Unify the animation of the drop-down menu.
- Optimize the display of the application drop-down menu.
- Replaced the default font to optimize the display.

### Fixed

- Fix the problem of failed to create storage space

## [0.2.9] - 2022-02-18

### Added

- Add a simple notification function

### Changed

- Custom installation of new parameters(Capabilities,Hostname,Privileged)
- Update front-end translation [@SemVer](https://github.com/zarevskaya) [@koboldMaki](https://github.com/koboldMaki) [@sgastol](https://github.com/sgastol) [@delki8](https://github.com/delki8)

- Modify the default location and name of the usb mount

### Fixed

- Fix the problem of being indexed by search engines
- Fix some style display issues
- Solve hard drive can't be formatted, can't finish adding storage

## [0.2.8] - 2022-01-30

### Added

- Add USB disk device display

### Changed

- Update translation [@baptiste313](https://github.com/baptiste313) [@thueske](https://github.com/thueske)
- Compatible with more types of drives

### Fixed

- Fix the language initialization bug
- Fix the problem that the login page could not be displayed
- Fix missing translated content

## [0.2.7] - 2022.01.26

### Changed

- Apply multilingual support

### Security

- Fix an injectable execution bug

## [0.2.6] - 2022.01.26

### Added

- Add a bug report panel.
- App Store apps start supporting multiple languages

### Fixed

- Fix a disk that cannot be formatted under certain circumstances

## [0.2.5] - 2022.01.24

### Added

- Storage Manager

### Changed

- Update Disk widget
- Update language files [@ImOstrovskiy](https://github.com/ImOstrovskiy) [@baptiste313](https://github.com/baptiste313)

### Fixed

- File synchronization issues
- Fix the app store classification problem

## [0.2.4] - 2021.12.30

### Changed

- Brand new App Store
- Optimize request method

### Fixed

- Fix Sync panel width display error.
- Fix App panel width display error.

## [0.2.3] - 2021.12.11

### Added

- Add detailed CPU and memory statistics.
- Add the multi-language function and add Chinese translation.
- Add the function to modify the search engine.
- Add the function of modifying the WebUI port

### Changed

- Update update script
- Preprocessing usb automounting

### Fixed

- Volume path problem when customizing the installation of applications
- Fix Cpu and Ram usage display error
- Fix translation errors
- Fixed an error when importing and exporting appfile.

## [0.2.2] - 2021.12.02

### Changed

- UI adjustment

### Fixed

- Fix the problem of data display error when manually installing apps
- Fix some spelling problems
- Fix the bug of synchronization module

## [0.2.1] - 2021.11.25

### Fixed

- Fix Sync display error
- Fix Sync Downoad url error
- Fix Smart Block display error
- Fix widgets settings dispaly error
- Fix  application installation path error

## [0.2.0] - 2021.11.25

### Added

- Add sync function


## [0.1.11] - 2021.11.10

### Changed

- Adaptation of cell phone terminals
- Optimize user experience
- Replaced the default background
- Optimized the display performance and fixed some bugs

### Fixed

- Resolve application installation path errors

## [0.1.10] - 2021.11.04

### Added

- Add application terminal
- Add application logs
- Add system logs
- Add App Store for installation

## [0.1.9] - 2021.11.01 [YANKED]

## [0.1.8] - 2021.10.27

### Added

- Add system terminal
- Add the ability to modify the user name and password

### Changed

- Experience optimization
- Improve single user management function
- Fixed Disk widget display error
- Fixed Username display error after change
- Adaptation for mobile access

## [0.1.7] - 2021.10.22

### Added

- Add user authentication module, Login page and initialization page.

### Fixed

- Fix the problem that the application could not start after the system restarted.
- Home storage space data display exception
- Script override causes application loss after installation
- Fix docker network error

## [0.1.6] - 2021.10.19

### Added

- Add app icon auto-fill via docker image name.
- Add a file selector for app install.

### Changed

- Modify import reminder.
- Optimize the application installation process

### Fixed

- Fixed an issue with the app were it would disappear when the app was modified.
- Fixed device selector default dir to /dev

## [0.1.5] - 2021.10.15

### Added

- Add CPU RAM Status with widget
- Add Disk Info with widget
- Realize automatic loading of widgets

### Changed

- Enhance the Docker cli import experience and automatically fill in the folders that need to be mounted

### Removed

- Remove Weather widget.

### Fixed

- AppFile upload does not pass verification
- The setting menu of the app is displayed abnormally when the browser window is too narrow
- The port is occupied and the program cannot start
- Fix display bugs when windows size less than 1024px

## [0.1.4] - 2021.09.30

### Added

- Import and export of application configuration files
- Automatic parsing of docker commands

### Changed

- Improve the program release process
- Application installation process UX/UI optimization

### Fixed

- Authentication failure during the operation, resulting in the need to re-login

## [0.1.3] - 2021.09.29 [YANKED]

## [0.1.2] - 2021.09.28

### Fixed

- Application modification and new creation failure issues

## [0.1.1] - 2021.09.27

## [0.1.0] - 2021.09.26

### Added

- Application Center
