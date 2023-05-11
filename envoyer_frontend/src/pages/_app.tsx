import '@/styles/globals.css';
import type { AppProps } from 'next/app';
import { QueryClient, QueryClientProvider } from 'react-query';
import { ReactQueryDevtools } from 'react-query/devtools';
import { SessionProvider, useSession } from 'next-auth/react';
import { PropsWithChildren, ReactNode, useEffect } from 'react';
import SiteLayout from '@/layouts/siteLayout';
import { USER_ROLE } from '@/utils/Constants';
import Spinner from '@/components/spinner';
import { ToastContainer } from 'react-toastify';
import RouteChangeIndicator from '@/components/RouteChangeIndicator';
import { ConfigProvider } from 'antd';

function getDefaultLayout(children: ReactNode): ReactNode {
  return <SiteLayout>{children}</SiteLayout>;
}

export default function App(props: AppProps | any) {
  const {
    Component,
    pageProps: { session, ...pageProps },
  } = props;

  const queryClient = new QueryClient();

  const renderAppLayout = () => {
    const children = (
      <ConfigProvider
        theme={{
          token: {},
        }}
      >
        <Component {...pageProps} />
      </ConfigProvider>
    );
    const { getLayout = getDefaultLayout } = Component.pageOptions || {};

    return getLayout(children);
  };
  if (session === 'loading') return <Spinner.FullPage />;

  return (
    <QueryClientProvider client={queryClient}>
      <SessionProvider session={session}>
        <AuthManager {...props}>
          {renderAppLayout()}
          <ReactQueryDevtools initialIsOpen={false} />
        </AuthManager>
        <RouteChangeIndicator />
        <ToastContainer
          className={'font-Metropolis'}
          pauseOnFocusLoss={false}
          autoClose={500}
          hideProgressBar={true}
        />
      </SessionProvider>
    </QueryClientProvider>
  );
}

function AuthManager({ Component, children, router }: PropsWithChildren<any>) {
  const { data: session, status }: any = useSession();
  const { redirectIfAuthenticated = false, requiresAuth = false } =
    Component.pageOptions || {};

  useEffect(() => {
    if (status === 'loading') return;
    if (requiresAuth && !session) {
      router.replace(
        `/auth/login?callbackUrl=${encodeURIComponent(window.location.href)}`
      );
      return;
    }

    if (!!session && redirectIfAuthenticated) {
      let redirectUrl = '/dashboard/' + session?.app_id;
      if (session?.role === USER_ROLE.SUPER_ADMIN) {
        redirectUrl = '/super-admin/dashboard';
      }

      router.replace(redirectUrl);
    }
  }, [status, redirectIfAuthenticated, requiresAuth, router, session]);

  if (
    status === 'loading' ||
    (requiresAuth && !session) ||
    (!!session && redirectIfAuthenticated)
  ) {
    return <Spinner.FullPage />;
  }
  return <>{children}</>;
}
