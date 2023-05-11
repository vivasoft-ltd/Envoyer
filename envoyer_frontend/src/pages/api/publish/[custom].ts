import PublishService from '@/services/publishService';
import { NextApiRequest, NextApiResponse } from 'next';

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const custom = req.query.custom as string;

  const publishService = new PublishService();

  try {
    const publishRes: any = await publishService.publishInQueueCustom(
      custom,
      req.body
    );
    res.json(publishRes.data);
  } catch (error: any) {
    res.status(error.status).json(error.data);
  }
}
