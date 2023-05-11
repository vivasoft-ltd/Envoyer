import { useRouter } from 'next/router';
import React, { ReactNode } from 'react';
import Admin from '@/layouts/Admin';
import { USER_ROLE } from '@/utils/Constants';
import { GetServerSidePropsContext } from 'next';
import { getSession } from 'next-auth/react';
import { Typography } from 'antd';
import { useQuery } from 'react-query';
import EventService from '@/services/eventService';
import TemplateList from '@/components/templates/TemplateList';

const { Title } = Typography;

const EventTemplates = () => {
  const router = useRouter();
  const eventService = new EventService();

  const appId = router.query?.['app-id'];
  const eventId = router.query?.['event-id'];

  const { data: getEven } = useQuery(['getEvent', eventId], () =>
    eventService.getEvent(Number(eventId))
  );

  return (
    <div>
      <Title level={3}>Templates - {getEven?.name}</Title>
      <TemplateList eventId={Number(eventId)} />
    </div>
  );
};

export default EventTemplates;

EventTemplates.pageOptions = {
  requiresAuth: true,
  getLayout: (children: ReactNode) => <Admin>{children}</Admin>,
};

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const session = await getSession(context);

  if (
    session?.role === USER_ROLE.SUPER_ADMIN ||
    session?.app_id === Number(context.params?.['app-id'])
  ) {
    return {
      props: {},
    };
  } else {
    return {
      redirect: {
        destination: '/template/' + session?.app_id,
        permanent: true,
      },
    };
  }
}
