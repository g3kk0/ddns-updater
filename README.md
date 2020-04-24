# DDNS Updater

Dynamic DNS updater for Google Cloud DNS.

## Requirements

- Google service account (with DNS Administrator role)

## Usage

Build the binary
```
git clone https://github.com/g3kk0/ddns-updater.git
cd ddns-updater
go build
```

Set the appropriate environment variables and run the binary
```
export GOOGLE_APPLICATION_CREDENTIALS="sa.json"
export GOOGLE_PROJECT_ID="foo-bar-12345678"
export DNS_ZONE="foo-com"
export DNS_RECORD="bar.foo.com"
./ddns-updater
```
