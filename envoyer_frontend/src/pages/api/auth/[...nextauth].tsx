import NextAuth, { NextAuthOptions } from 'next-auth';
import CredentialsProviders from 'next-auth/providers/credentials';
import { NextApiRequest, NextApiResponse } from 'next';
import AuthService from '@/services/authServices';
import logger from '@/lib/logger';

const MAX_AGE = parseInt(process.env.NEXTAUTH_SESSION_DURATION || '0', 10);

type NextAuthOptionsCallback = (
  req: NextApiRequest,
  res: NextApiResponse
) => NextAuthOptions;

export async function refreshAccessToken(token: any) {
  const authServices = new AuthService();
  const response = await authServices.refreshToken(token.refresh_token);
  if ((response.status ?? 0) === 401) {
    token.error = 'refresh_token_error';
  } else {
    token.access_token = response.access_token;
    token.refresh_token = response.refresh_token;
    token.exp += MAX_AGE;
  }
  return token;
}

export const authOptions: NextAuthOptions = {
  callbacks: {
    async signIn({ user, account }: any) {
      return !!user.access_token;
    },

    async jwt({ token, user, account }) {
      if (Date.now() > (token.exp as number) * 1000) {
        token = await refreshAccessToken(token);
      }

      return {
        ...user,
        ...token,
        ...account,
      };
    },

    async session({ token }: any) {
      return token;
    },
    async redirect({ url, baseUrl }) {
      // Allows relative callback URLs
      if (url.startsWith('/')) return `${baseUrl}${url}`;
      // Allows callback URLs on the same origin
      else if (new URL(url).origin === baseUrl) return url;
      return baseUrl;
    },
  },

  debug: process.env.NODE_ENV === 'development',

  events: {},
  jwt: {
    maxAge: MAX_AGE,
  },
  logger: {
    error(code, ...message) {
      logger.error(`[next-auth] ${code} ${message.join(', ')}`);
    },
    warn(code, ...message) {
      logger.warn(`[next-auth] ${code} ${message.join(', ')}`);
    },
    debug(code, ...message) {
      logger.debug(`[next-auth] ${code} ${message.join(', ')}`);
    },
  },
  pages: {
    signIn: '/auth/login',
    // error: '/auth/login',
  },
  providers: [
    CredentialsProviders({
      // The name to display on the sign in form (e.g. "Sign in with...")
      name: 'Credentials',
      // `credentials` is used to generate a form on the sign in page.
      // You can specify which fields should be submitted, by adding keys to the `credentials` object.
      // e.g. domain, username, password, 2FA token, etc.
      // You can pass any HTML attribute to the <input> tag through the object.
      credentials: {
        username: { label: 'Username', type: 'text', placeholder: 'jsmith' },
        password: { label: 'Password', type: 'password' },
      },
      async authorize(credentials: any, req): Promise<any> {
        // console.log({ credentials });
        const AuthServices = new AuthService();
        const response = await AuthServices.authenticate({
          username: credentials.username,
          password: credentials.password,
        });
        // console.log({ response });

        if (response) {
          // Any object returned will be saved in `user` property of the JWT
          return response.data;
        } else {
          // If you return null then an error will be displayed advising the user to check their details.
          return null;

          // You can also Reject this callback with an Error thus the user will be sent to the error page with the error message as a query parameter
        }
      },
    }),
  ],

  // secret: process.env.NEXTAUTH_SECRET,
  session: {
    strategy: 'jwt',
    maxAge: MAX_AGE,
  },
};

export const nextAuthOptions: NextAuthOptionsCallback = () => {
  return authOptions;
};

const handler = (req: NextApiRequest, res: NextApiResponse) => {
  return NextAuth(req, res, nextAuthOptions(req, res));
};

export default handler;
