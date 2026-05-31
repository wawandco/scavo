# Contributing to Scavo

Thank you for your interest in Scavo!

Scavo is primarily an internal tool built for Wawandco team events. We open-sourced it as a showcase of our engineering practices with Go, HTMX, Tailwind, and LeapKit.

## Ways to contribute

- Bug reports and security issues are welcome.
- Small improvements and refactors that align with the existing simple, explicit style are appreciated.
- Large feature additions are unlikely to be accepted unless they solve a clear problem for event-style scavenger hunts.

## Development

See the [README](README.md) for setup instructions (`go tool dev`).

Please keep changes small, focused, and well-tested where possible.

## Code style

- Follow existing patterns (raw SQL with explicit error handling, context-injected DB, Plush templates).
- Prefer clarity and simplicity over cleverness.
- Run `go fmt ./...` and `go vet ./...` before submitting.

## License

Contributions are licensed under the same O’Saasy license as the project.
