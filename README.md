# DDNS Updater

Dynamic DNS updater for Google Cloud DNS.

## Requirements

- Google service account (with DNS Administrator role)

## Usage

Download the latest version
```
curl xxx -o /usr/local/bin/
chmod +x /usr/local/bin/ddns-updater
```

Set the appropriate environment variables and run the tool
```
export GOOGLE_APPLICATION_CREDENTIALS="sa.json"
export GOOGLE_PROJECT_ID="foo-bar-12345678"
export DNS_ZONE="foo-com"
export DNS_RECORD="bar.foo.com"
./ddns-updater
```
