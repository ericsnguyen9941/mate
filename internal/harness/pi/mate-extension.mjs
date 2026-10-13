// Mate-owned pi 1.0.4 extension. No dependencies outside Node and pi's API.
import { readFileSync } from "node:fs";
import { spawn } from "node:child_process";

const config = JSON.parse(
  readFileSync(new URL("./mate-launch.json", import.meta.url), "utf8"),
);

export function runHook(binary, args, payload, cwd) {
  return new Promise((resolve, reject) => {
    const child = spawn(binary, ["hook", ...args], {
      cwd,
      stdio: ["pipe", "pipe", "pipe"],
    });
    let stdout = "";
    let stderr = "";
    child.stdout.setEncoding("utf8");
    child.stderr.setEncoding("utf8");
    const timer = setTimeout(() => {
      child.kill("SIGKILL");
      reject(new Error("Mate hook timed out"));
    }, 30000);
    child.stdout.on("data", (chunk) => {
      stdout += chunk;
      if (Buffer.byteLength(stdout, "utf8") > 65536) {
        child.kill("SIGKILL");
        reject(new Error("Mate hook output too large"));
      }
    });
    child.stderr.on("data", (chunk) => {
      stderr = (stderr + chunk).slice(-4096);
    });
    child.on("error", (err) => {
      clearTimeout(timer);
      reject(err);
    });
    child.on("close", (code) => {
      clearTimeout(timer);
      if (code === 0) resolve(stdout);
      else reject(new Error(`Mate hook exited ${code}: ${stderr}`));
    });
    child.stdin.on("error", () => {});
    child.stdin.end(JSON.stringify(payload));
  });
}

export default function mateExtension(pi) {
  if (typeof pi.registerMcpServer !== "function") {
    throw new Error("Mate requires pi >=1.0.4 with native MCP support");
  }
  for (const [name, server] of Object.entries(config.mcpServers)) {
    pi.registerMcpServer(name, server);
  }
  const warn = (ctx, err) => {
    const message = `Mate hook failed: ${err.message}`;
    if (ctx.hasUI) ctx.ui.notify(message, "error");
    else console.error(message);
  };
  const call = async (ctx, args, payload) => {
    try {
      return await runHook(config.binary, args, payload, ctx.cwd);
    } catch (err) {
      warn(ctx, err);
      return "Mate hook failed. Run mate recall before acting.";
    }
  };
  const identity = (ctx) => ({
    session_id: ctx.sessionManager.getSessionId(),
    transcript_path: ctx.sessionManager.getSessionFile(),
  });
  let source = config.resume ? "resume" : "startup";
  let recallPending = true;
  pi.on("session_start", (event) => {
    if (event.reason === "new") source = "clear";
    else if (event.reason === "resume") source = "resume";
    else if (event.reason === "reload" || event.reason === "fork") source = "startup";
    recallPending = true;
  });
  pi.on("session_compact", () => {
    source = "compact";
    recallPending = true;
  });
  pi.on("input", async (event, ctx) => {
    await call(ctx, ["mate-prompt"], { prompt: event.text });
    return { action: "continue" };
  });
  pi.on("before_agent_start", async (_event, ctx) => {
    if (!recallPending) return;
    const text = await call(ctx, ["mate-session", "--harness", "pi"], {
      ...identity(ctx),
      source,
    });
    recallPending = false;
    return { message: { customType: "mate-recall", content: text, display: false } };
  });
  pi.on("agent_end", async (event, ctx) => {
    const last = [...event.messages].reverse().find(
      (message) => message.role === "assistant" && message.stopReason !== "toolUse",
    );
    if (!last) return;
    const text = (last.content ?? [])
      .filter((block) => block.type === "text")
      .map((block) => block.text)
      .join("\n");
    await call(ctx, ["mate-stop"], {
      ...identity(ctx),
      last_assistant_message: text,
    });
  });
}
