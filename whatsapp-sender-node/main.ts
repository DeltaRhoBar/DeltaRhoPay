import { Client, NoAuth } from "whatsapp-web.js";
import qrcode from "qrcode-terminal";
import { Application, Router } from "@oak/oak";

// Constant

const LISTENING_PORT = 4242;
const WHATS_APP_DEVICE_NAME = "Bar Server";

// ----- Types -----
// WhatsApp Message to send to a client
type Message = {
  number: string;
  message: string;
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
  messagesToSend: Message[];
  progress: number;
  messagesSend: Message[];
  loginQrCode?: string;
  error?: string;
  createdAt?: number;
  startedAt?: number;
  finishedAt?: number;
  cancelled?: boolean;
};

const jobs = new Map<string, Job>();
const queue: string[] = [];
let workerActive = false;

function enqueueJob(messagesToSend: Message[]) {
  const id = crypto.randomUUID();
  const job: Job = {
    id,
    status: "queued",
    messagesToSend: messagesToSend,
    progress: 0,
    messagesSend: [],
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
      const whatsappClient = new Client({
        authStrategy: new NoAuth(),
        deviceName: WHATS_APP_DEVICE_NAME,
      });
    } catch (err) {
      job.status = "failed";
      job.error = String(err);
      job.finishedAt = Date.now();
    }
  }

  workerActive = false;
}

const router = new Router();

const whatsappClient = new Client({
  authStrategy: new NoAuth(),
  deviceName: WHATS_APP_DEVICE_NAME,
  puppeteer: {
    headless: true,
  },
});

whatsappClient.once("ready", async () => {
  console.log("Client is ready!");
  try {
    await whatsappClient.sendMessage("49xxxxxxxxxxxx@c.us", "Hallo Welt!");
  } catch (e) {
    console.log(e);
  }
});

whatsappClient.on("qr", (qr) => {
  qrcode.generate(qr, { small: true });
});

whatsappClient.on("message", (msg) => {
  console.log(msg);
});

whatsappClient.on("authenticated", () => {
  console.log("AUTHENTICATED");
});

whatsappClient.on("auth_failure", (msg) => {
  // Fired if session restore was unsuccessful
  console.error("AUTHENTICATION FAILURE", msg);
});

whatsappClient.initialize();
