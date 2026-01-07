export enum Role {
  Viewer = 0,
  Editor = 1,
  Manager = 2,
}

export interface User {
  id: string;
  username: string;
  email: string;
  avatar: string;
  isRoot: boolean;
  isAdmin: boolean;
  role: Role;
  otpActive: boolean;
  otpVerified: boolean;
  isBlocked?: boolean;
}