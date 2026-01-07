export interface SessionData {
  token: string;
  expires: number;
}

export function isSessionExisting(): boolean {
  return localStorage.getItem('nauri_session') !== null;
}

export function getSession(): SessionData | undefined {
  const session = localStorage.getItem('nauri_session');
  if (!session) return undefined;
  return JSON.parse(session);
}

export function setSession(token: string, expires: number): void {
  const session = { token, expires };
  localStorage.setItem('nauri_session', JSON.stringify(session));
}

export function destroySession(): void {
  localStorage.removeItem('nauri_session');
}