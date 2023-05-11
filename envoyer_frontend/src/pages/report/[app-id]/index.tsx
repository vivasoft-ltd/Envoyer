import {useRouter} from 'next/router';
import React, {ReactNode} from 'react';
import Admin from '@/layouts/Admin';
import {USER_ROLE} from '@/utils/Constants';
import {GetServerSidePropsContext} from 'next';
import {getSession} from 'next-auth/react';
import ErrorLogList from "@/components/report/ErrorLogList";

const ReportId = () => {
  const router = useRouter();

  const appId = router.query?.['app-id'];

  return (
    <div>
      <div>
        <h2 className='text-2xl font-semibold pb-2'>Error Logs</h2>
        <ErrorLogList appId={Number(appId)}/>
      </div>
    </div>
  );
};

export default ReportId;

ReportId.pageOptions = {
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
        destination: '/report/' + session?.app_id,
        permanent: true,
      },
    };
  }
}
