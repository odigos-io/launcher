module github.com/keyval-dev/launcher

go 1.18

require github.com/Binject/binjection v0.0.0-20210701074423-605d46e35deb

require (
	github.com/Binject/debug v0.0.0-20210312092933-6277045c2fdf // indirect
	github.com/Binject/shellcode v0.0.0-20191101084904-a8a90e7d4563 // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
)

replace (
	github.com/Binject/debug => ./debug
	github.com/Binject/binjection => ./binjection
)
