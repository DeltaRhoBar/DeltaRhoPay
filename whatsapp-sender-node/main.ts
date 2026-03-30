import { Client, NoAuth } from "whatsapp-web.js";
import qrcode from "qrcode-terminal";

const client = new Client({
  authStrategy: new NoAuth(),
  deviceName: "Bar Server",
  puppeteer: {
    headless: true,
  },
});

client.once("ready", async () => {
  console.log("Client is ready!");
  try {
    await client.sendMessage("49xxxxxxxxxxxx@c.us", "Hallo Welt!");
  } catch (e) {
    console.log(e);
  }
});

client.on("qr", (qr) => {
  qrcode.generate(qr, { small: true });
});

client.on("message", (msg) => {
  console.log(msg);
});

client.on("authenticated", () => {
  console.log("AUTHENTICATED");
});

client.on("auth_failure", (msg) => {
  // Fired if session restore was unsuccessful
  console.error("AUTHENTICATION FAILURE", msg);
});

client.initialize();
