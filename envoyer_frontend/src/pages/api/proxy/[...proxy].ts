import {NextApiRequest, NextApiResponse} from 'next';
import proxyHandler from '@/lib/proxy-request-handle';
import {getToken} from 'next-auth/jwt';

const publicAPIs = [
  '/api/proxy/v2/auth/login',
  '/api/proxy/v2/auth/user/register',
  '/api/proxy/v2/auth/password/forgot',
  '/api/proxy/v2/auth/password/reset',
];

const publicTokenAPIs = [
  '/api/proxy/v1/user_verification',
  '/api/proxy/v1/auth/password',
];

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const jwtPayload = await getToken({
    req,
    secret: process.env.NEXTAUTH_SECRET!,
  });

  const requiresAuth = !publicAPIs.includes(req.url || '');

  if (
    !jwtPayload &&
    requiresAuth &&
    !!publicTokenAPIs.includes(req.url || '')
  ) {
    return res.status(401).send({
      status: 401,
      message: 'Unauthorized',
      timestamp: new Date().toISOString(),
    });
  }

  if (jwtPayload) {
    req.headers['authorization'] = `Bearer ${jwtPayload.access_token}`;
  }
  // Delete data that's not needed in backend
  delete req.headers['cookie'];
  delete req.query['proxy'];

  return proxyHandler(req, res);
}

export const config = {
  api: {
    // - https://nextjs.org/docs/api-routes/api-middlewares#custom-config
    externalResolver: true,
    bodyParser: false, // not to use url encoded form like streaming POST request
  },
};
