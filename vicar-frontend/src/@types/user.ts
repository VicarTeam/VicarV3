export interface User {
  id: string;
  username: string;
  avatar?: string;
  isTeam: boolean;
  otpActive: boolean;
  otpVerified: boolean;
  isBlocked?: boolean;
}