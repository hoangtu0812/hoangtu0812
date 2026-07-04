# Go Install

You can find the relevant installation files at https://golang.org/dl/.

Follow the instructions related to your operating system. To check if Go was installed successfully, you can run the following command in a terminal window:
```powershell
go version
```
Which should show the version of your Go installation.

## Go Quickstart

Let's create our first Go program.

Launch the VS Code editor
Open the extension manager or alternatively, press Ctrl + Shift + x
In the search box, type "go" and hit enter
Find the Go extension by the GO team at Google and install the extension
After the installation is complete, open the command palette by pressing Ctrl + Shift + p
Run the Go: Install/Update Tools command
Select all the provided tools and click OK
VS Code is now configured to use Go.

Open up a terminal window and type:

```powershell
go mod init example.com/hello
```

Run
```powershell
go run .\helloworld.go
```

If you want to save the program as an executable, type and run:
```powershell
go build .\helloworld.go
```