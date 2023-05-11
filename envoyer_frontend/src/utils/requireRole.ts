import { NextApiRequest, NextApiResponse } from 'next';

export function requireRole(role: string) {
  return (req: any, res: NextApiResponse, next: () => void) => {
    if (req.session?.role === role) {
      next(); // Allow the request to proceed
    } else {
      return {
        redirect: {
          destination: '/dashboard/',
          permanent: true,
        },
      };
    }
  };
}
