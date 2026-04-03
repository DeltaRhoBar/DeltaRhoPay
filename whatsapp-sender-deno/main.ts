import { Context, Router } from "@oak/oak";
import qrcode from "qrcode-terminal";
import { Client, NoAuth } from "whatsapp-web.js";

// Constant

const LISTENING_PORT = 4242;
const WHATS_APP_DEVICE_NAME = "Bar Server";
const WHATS_APP_SEND_TIMEOUT = 1000 * 60 * 30; // 30 min timeout
const WHATS_APP_SEND_MIN_DELAY = 1000; // 1 second
const WHATS_APP_SEND_MAX_DELAY = 4000; // 4 seconds

// ----- Types -----
// WhatsApp Message to send to a client
type Message = {
  number: string;
  message: string;
  error?: string;
};
// Status of a Message Job
type JobStatus =
  | "queued"
  | "running"
  | "success"
  | "failed"
  | "cancelled"
  | "loginRequired";

type Job = {
  id: string;
  status: JobStatus;
  messages: Message[];
  progress: number;
  messagesSend: Message[];
  messagesFailed: Message[];
  loginQrCode?: string;
  error?: string;
  createdAt?: number;
  startedAt?: number;
  finishedAt?: number;
  cancelled?: boolean;
};

// ---- Classes ----
class SimpleBlocker {
  private promise!: Promise<void>;
  private _res!: () => void;
  private _rej!: (reason?: unknown) => void;

  constructor() {
    this._reset();
  }

  _reset() {
    this.promise = new Promise<void>((res, rej) => {
      this._res = res;
      this._rej = rej;
    });
  }

  wait(): Promise<void> {
    return this.promise;
  }

  release(): void {
    this._res();
    this._reset();
  }

  fail(err?: unknown): void {
    this._rej(err);
    this._reset();
  }
}

const jobs = new Map<string, Job>();
const queue: string[] = [];
let workerActive = false;

function enqueueJob(messagesToSend: Message[]) {
  const id = crypto.randomUUID();
  const job: Job = {
    id,
    status: "queued",
    messages: messagesToSend,
    progress: 0,
    messagesSend: [],
    messagesFailed: [],
  };
  jobs.set(id, job);
  queue.push(id);
  startWorkerLoop();
  return id;
}

async function startWorkerLoop() {
  if (workerActive) return;
  workerActive = true;

  while (queue.length > 0) {
    const id = queue.shift()!;
    const job = jobs.get(id);
    if (!job) continue;
    if (job.cancelled) {
      job.status = "cancelled";
      job.finishedAt = Date.now();
      continue;
    }

    job.status = "running";
    job.startedAt = Date.now();
    job.progress = 0;

    // Send messages to users
    try {
      await sendMessages(job);
    } catch (err) {
      job.status = "failed";
      job.error = String(err);
      job.finishedAt = Date.now();
    }
  }

  workerActive = false;
}

function sleep(ms: number) {
  return new Promise<void>((resolve) => setTimeout(resolve, ms));
}

async function sendMessages(job: Job) {
  const whatsappClient = new Client({
    authStrategy: new NoAuth(),
    deviceName: WHATS_APP_DEVICE_NAME,
    puppeteer: {
      headless: true,
    },
  });

  const block = new SimpleBlocker();

  whatsappClient.on("qr", (qr) => {
    console.log("WhatsApp: QR code for login required");
    job.loginQrCode = qr;
    job.status = "loginRequired";
  });

  whatsappClient.on("authenticated", () => {
    console.log("WhatsApp: authenticated");
  });

  whatsappClient.on("disconnected", () => {
    console.log("WhatsApp: disconnected");
  });

  whatsappClient.on("auth_failure", () => {
    console.log("WhatsApp: auth_failure");
    block.fail(new Error("WhatsApp: auth_failure"));
  });

  whatsappClient.once("ready", async () => {
    for (const message of job.messages) {
      // update total progress
      job.progress = job.messages.length /
        (job.messagesSend.length + job.messagesFailed.length);

      // sleep for random delay before sending each message
      await sleep(
        Math.random() * (WHATS_APP_SEND_MAX_DELAY - WHATS_APP_SEND_MIN_DELAY) +
          WHATS_APP_SEND_MIN_DELAY,
      );

      // Convert 0049xxxxx -> 49xxxxx@c.us
      // convert number to whatss id (if not already)
      let number = message.number.includes("@c.us")
        ? message.number
        : `${message.number}@c.us`;
      // remove leading 0s
      number = number.replace(/^(?:0)+/, "");

      // send message to client
      try {
        console.log(`Sending message to ${number}`);
        await whatsappClient.sendMessage(number, message.message);
        job.messagesSend.push(message);
        console.log("Message send");
      } catch (err) {
        console.log(`Failed sending message: ${err}`);
        message.error = String(err);
        job.messagesFailed.push(message);
        continue;
      }
    }

    block.release();
  });

  whatsappClient.initialize();

  // wait (max time: WHATS_APP_SEND_TIMEOUT) for promise completion
  setTimeout(
    () => block.fail(new Error("Timeout on sending messages")),
    WHATS_APP_SEND_TIMEOUT,
  );
  try {
    await block.wait();
  } catch (err) {
    job.error = String(err);
    job.status = "failed";
    job.finishedAt = Date.now();
  }
}

const router = new Router();

router.post("/jobs", async (ctx: Context) => {
  const bodyReader = ctx.request.body;
  if (bodyReader.type() !== "json") {
    ctx.response.status = 400;
    ctx.response.body = { error: "Expected application/json" };
    return;
  }

  try {
    const body = await bodyReader.json();
  } catch (err) {
    ctx.response.status = 400;
    ctx.response.body = { error: `Invalid JSON: ${err}` };
    return;
  }
});
