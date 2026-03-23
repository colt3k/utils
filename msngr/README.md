# msngr

Provides various ways to show a message.

- espeak
- freedesktop
- notifyicon
- nsuser
- pushbullet
- say
- slack
- speech synthesizer

## Platform Notes

- macOS: `nsuser`, `say`
- Linux: `freedesktop`, `espeak`
- Windows: `notifyicon`, `speechsynthesizer`
- cross-service integrations: `pushbullet`, `slack`

Each backend keeps its own `Notification` shape so applications can use the fields required by that transport only.

## Development

- `cd msngr && go test ./...`
