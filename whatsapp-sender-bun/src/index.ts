import { Client, NoAuth } from "whatsapp-web.js";
import { Hono } from "hono";
import { z } from "zod";

// Constant

const LISTENING_PORT = 4242;
const WHATS_APP_DEVICE_NAME = "Bar Server";
const WHATS_APP_SEND_TIMEOUT = 1000 * 60 * 30; // 30 min timeout
const WHATS_APP_SEND_MIN_DELAY = 3000; // 3 seconds
const WHATS_APP_SEND_MAX_DELAY = 6000; // 6 seconds
const SHORT_DELAY_MIN = 100;
const SHORT_DELAY_MAX = 800;

// ----- Types -----
// WhatsApp Message to send to a client
type Message = {
  number: string;
  message: string;
  error?: string;
};
// Schema for Message
const MessageSchema = z.object({
  number: z.string().min(5),
  message: z.string().min(5),
}).strict();
const MessageArraySchema = z.array(MessageSchema);

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

function rand_sleep() {
  return new Promise<void>((resolve => setTimeout(resolve, Math.random() * (SHORT_DELAY_MAX - SHORT_DELAY_MIN) + SHORT_DELAY_MIN)));
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
    job.status = "running";
  });

  whatsappClient.on("disconnected", (reason) => {
    console.log(`WhatsApp: disconnected: ${reason}`);
  });

  whatsappClient.on("auth_failure", () => {
    console.log("WhatsApp: auth_failure");
    block.fail(new Error("WhatsApp: auth_failure"));
  });

  whatsappClient.on("change_state", (state) => {
    console.log(`WhatsApp: change_state: ${state}`);
  });

  whatsappClient.once("ready", async () => {
    for (const message of job.messages) {
      if (job.cancelled) {
        job.status = "cancelled";
        job.finishedAt = Date.now();
        break;
      }
      // update total progress
      job.progress = (job.messagesSend.length + job.messagesFailed.length) / job.messages.length;

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
        const contact = await whatsappClient.getContactById(number);
        await rand_sleep();
        console.log(`Sending message to ${number}`);
        await whatsappClient.sendMessage(number, message.message);
        await rand_sleep();
        console.log(`Contact in users contacts: ${contact.isMyContact}`);
        if (!contact.isMyContact) {
          console.log('Reciever not in contacts -> archiving\nGetting Chat...');
          const chat = await contact.getChat();
          await rand_sleep();
          console.log(`Got chat: ${chat.name}, Archived: ${chat.archived}`);
          if (chat.archived) continue;
          console.log('Archiving...')
          await chat.archive();
          await rand_sleep();
        }
        job.messagesSend.push(message);
        console.log("Message send");
      } catch (err) {
        console.log(`Failed sending message: ${err}`);
        message.error = String(err);
        job.messagesFailed.push(message);
        continue;
      }
    }

    console.log('Finished sending messages. Releasing job ...');

    await rand_sleep();

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
    job.status = "success";
    job.finishedAt = Date.now();
    await whatsappClient.logout();
  } catch (err) {
    job.error = String(err);
    job.status = "failed";
    job.finishedAt = Date.now();
    await whatsappClient.logout();
  }
}

const app = new Hono();

app.post('/jobs', async (ctx) => {
  const body = await ctx.req.json().catch(() => null);
  const result = MessageArraySchema.safeParse(body);
  if (!result.success) {
    return ctx.json({ error: 'validation_vailed', details: z.treeifyError(result.error) }, 400);
  }

  const messages: Message[] = result.data;
  const jobId = enqueueJob(messages);

  return ctx.json({ jobId: jobId }, 202);
});

app.get('/jobs/:id', (ctx) => {
  const id = ctx.req.param('id');
  const job = jobs.get(id);
  if (!job) return ctx.json({ error: 'not_found' }, 404);

  return ctx.json({ id: job.id, job: job });
})

app.post('jobs/:id/cancel', (ctx) => {
  const id = ctx.req.param('id');
  const job = jobs.get(id);
  if (!job) return ctx.json({ error: 'not_found' }, 404);

  if (job.status === "queued") {
    job.cancelled = true;

    // remove from queue
    for (let i = queue.length - 1; i >= 0; i--) {
      if (queue[i] === id) queue.splice(i, 1);
    }

    job.status = "cancelled";
    job.finishedAt = Date.now();
    return ctx.json({ cancelled: true }, 200);
  } else if (job.status === "running") {
    job.cancelled = true;
    return ctx.json({ cancelled: true }, 202);
  }

  return ctx.json({ cancelled: false, reason: 'already finished' });
});

app.get('/status', (ctx) => {
  return ctx.json({
    queded: queue.length,
    workerActive,
    totalJobs: jobs.size,
  });
});

export default {
  port: LISTENING_PORT,
  fetch: app.fetch,
}
