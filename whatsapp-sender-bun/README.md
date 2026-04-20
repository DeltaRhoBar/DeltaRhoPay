# whatsapp-sender-bun

To install dependencies:

```bash
bun install
```

To run:

```bash
bun run 
```

This project was created using `bun init` in bun v1.3.11. [Bun](https://bun.com) is a fast all-in-one JavaScript runtime.

## Build Container

```bash
podman build -t whatsapp_sender:latest .
```

## Run Container

```bash
podman run --rm -it -p 4242:4242 --cap-add SYS_ADMIN -v ./screenshots:/usr/src/app/screenshots whatsapp_sender:latest
```
