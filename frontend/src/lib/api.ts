const BASE = "http://localhost:8088/api/v1";

async function get(path: string) {
  const res = await fetch(`${BASE}${path}`);
  return res.json();
}

async function post(path: string, body: unknown) {
  const res = await fetch(`${BASE}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return res.json();
}

async function del(path: string) {
  const res = await fetch(`${BASE}${path}`, { method: "DELETE" });
  return res.json();
}

export { get, post, del };
