import {useRouter} from 'next/router';
import React, {ReactNode, useState} from 'react';
import Admin from '@/layouts/Admin';
import {USER_ROLE} from '@/utils/Constants';
import {GetServerSidePropsContext} from 'next';
import {getSession} from 'next-auth/react';
import {Button, Modal} from 'antd';
import EventCreate from '@/components/event/EventCreate';
import EventList from "@/components/event/EventList";

const EventId = () => {
  const router = useRouter();
  const [showCreateEventModal, setShowCreateEventModal] = useState(false);

  const hideModal = () => {
    setShowCreateEventModal(false);
  };
  const appId = router.query?.['app-id'];

  return (
    <div>
      <div className='text-right'>
        <Button
          type='primary'
          className='bg-blue-700'
          onClick={() => setShowCreateEventModal(true)}
        >
          Create Event
        </Button>
      </div>
      <div>
        <h2 className='text-2xl font-semibold pb-2'>Events</h2>
        <EventList appId={Number(appId)}/>
      </div>
      <div>
        <Modal
          centered
          open={showCreateEventModal}
          title='Create Event'
          onCancel={hideModal}
          footer={null}
          maskClosable={false}
        >
          <div className='pt-4'>
            <EventCreate hideModal={hideModal}></EventCreate>
          </div>
        </Modal>
      </div>
    </div>
  );
};

export default EventId;

EventId.pageOptions = {
  requiresAuth: true,
  getLayout: (children: ReactNode) => <Admin>{children}</Admin>,
};

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const session = await getSession(context);

  if (
    session?.role === USER_ROLE.SUPER_ADMIN ||
    (session?.role === USER_ROLE.DEVELOPER &&
      session?.app_id === Number(context.params?.['app-id']))
  ) {
    return {
      props: {},
    };
  } else {
    return {
      redirect: {
        destination: '/event/' + session?.app_id,
        permanent: true,
      },
    };
  }
}
