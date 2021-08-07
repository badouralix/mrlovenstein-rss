# Contributing

- [Bootstrap](#bootstrap)
- [Run Locally](#run-locally)

## Bootstrap

To create the go module

```bash
go mod init github.com/badouralix/mrlovenstein-rss
go mod tidy
```

## Run Locally

```bash
go run . &
curl http://localhost:8080/mrlovenstein.xml
```

Note that a popup might show up and ask for a confirmation to allow the app to listen on port 8080.

![Firewall Warning](https://user-images.githubusercontent.com/19719047/128614620-89cf8c62-8f39-4ddd-a2eb-8837cf063bbb.png)
