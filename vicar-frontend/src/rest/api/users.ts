import type {User} from "@/@types/user.ts";
import {DELETE, GET, PATCH, POST} from "@/rest";

export async function getUser(id: string = "@me", inAdminView: boolean = false): Promise<User|null> {
  return await GET<User>(`/users/${id}`, {view: inAdminView ? "admin" : undefined});
}

export async function getUserSessions(): Promise<string[]> {
  return await GET<string[]>("/users/@me/sessions", undefined, []);
}

export async function beginTotp(): Promise<string|null> {
  return (await POST<{url: string}>("/users/@me/totp"))?.url;
}

export async function verifyTotp(code: string): Promise<string[]|null> {
  const res = await POST<{codes: string[]}>("/users/@me/totp/verify", {code});
  return res?.codes || null;
}

export async function requestUrlForVerification(): Promise<string|null> {
  return (await GET<{url: string}>("/users/@me/totp/verify/url"))?.url;
}

export async function disableTotp(): Promise<boolean> {
  const res = await POST("/users/@me/totp/disable");
  return !!res;
}

export async function removeTotp(userId: string): Promise<boolean> {
  const res = await POST(`/users/${userId}/totp`);
  return !!res;
}

export async function beginPasskeyLink(displayName: string): Promise<{state: string, options: any}|null> {
  return await POST<{state: string, options: any}>("/users/@me/link/fido2/begin", {}, {display_name: displayName});
}

export async function endPasskeyLink(state: string, body: any): Promise<boolean> {
  const res = await POST("/users/@me/link/fido2/end", body, {state});
  return !!res;
}

export async function getPasskeyLinks(): Promise<string[]|null> {
  return await GET<string[]>("/users/@me/link/fido2/all");
}

export async function deletePasskeyLink(name: string): Promise<boolean> {
  const res = await DELETE(`/users/@me/link/fido2/${name}`);
  return !!res;
}

export async function setUserUsername(userId: string, username: string): Promise<boolean> {
  const res = await PATCH(`/users/${userId}/username`, {username});
  return !!res;
}

export async function setUserPassword(userId: string, password: string, oldPassword: string): Promise<boolean> {
  const res = await PATCH(`/users/${userId}/password`, {password, oldPassword});
  return !!res;
}