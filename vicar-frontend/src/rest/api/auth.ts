import {GET, POST, POST_AND_THROW, RestError, SuccessResponse} from "../index.ts";
import {setSession} from "../session.ts";

export function getStatus(): Promise<'CLAIMED'|'UNCLAIMED'|null> {
  return GET<'CLAIMED'|'UNCLAIMED'>('/auth/status');
}

export async function register(username: string, email: string, password: string, reg_type: 'root_claim'|'invite', token: string): Promise<boolean> {
  const res = await POST<SuccessResponse>(`/auth/register`, {
    username, email, password
  }, {
    reg_type, token
  });

  return !!res;
}

export async function login(username: string, password: string): Promise<boolean|{otp_missing: true, state: string}> {
  try {
    const res = await POST_AND_THROW<{token: string, expiresIn: number}>(`/auth/login`, {
      username, password
    });

    setSession(res.token, res.expiresIn);

    return true;
  } catch (e: any) {
    if (e instanceof RestError) {
      e = e as RestError<{state: string}>;

      if (e.status === 428) {
        return {otp_missing: true, state: e.data.state};
      }
    }

    console.error(e);
    return false;
  }
}

export async function loginTOTP(state: string, code: string): Promise<boolean> {
  const res = await POST<{token: string, expiresIn: number}>(`/auth/login/totp`, {
    state, code
  });

  if (!res) {
    return false;
  }

  setSession(res.token, res.expiresIn);

  return true;
}

export function beginFido2Login(): Promise<{state: string, options: any}|null> {
  return POST(`/auth/login/fido2/begin`);
}

export async function endFido2Login(state: string, body: any): Promise<boolean> {
  const res = await POST<{token: string, expiresIn: number}>(`/auth/login/fido2/end`, body, {state});

  if (!res) {
    return false;
  }

  setSession(res.token, res.expiresIn);

  return true;
}

export async function logout(): Promise<void> {
  await POST(`/auth/logout`);
}

export async function logoutAll(): Promise<boolean> {
  try {
    await POST_AND_THROW(`/auth/logout/all`);
    return true;
  } catch {
    return false;
  }
}