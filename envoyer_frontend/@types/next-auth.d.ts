import NextAuth, { DefaultSession } from 'next-auth';

declare module 'next-auth' {
  /**
   * Returned by `useSession`, `getSession` and received as a prop on the `SessionProvider` React Context
   */
  interface Session {
    role: string;
    access_token: string;
    exp: number;
    expired_at: number;
    iat: number;
    id: number;
    jti: string;
    provider: string;
    providerAccountId: number;
    refresh_token: string;
    type: string;
    app_id: number;
  }
}
