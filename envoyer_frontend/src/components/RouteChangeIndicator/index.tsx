import { useEffect } from 'react';
import { useRouter } from 'next/router';
import NProgress from 'nprogress';

const RouteChangeIndicator = () => {
  NProgress.configure({ showSpinner: false, speed: 300, trickleSpeed: 100 });
  const router = useRouter();

  useEffect(() => {
    const handleRouteChangeStart = () => NProgress.start();
    const handleRouteChangeFinish = () => NProgress.done();

    router.events.on('routeChangeStart', handleRouteChangeStart);
    router.events.on('routeChangeComplete', handleRouteChangeFinish);
    router.events.on('routeChangeError', handleRouteChangeFinish);

    return () => {
      router.events.off('routeChangeStart', handleRouteChangeStart);
      router.events.off('routeChangeComplete', handleRouteChangeFinish);
      router.events.off('routeChangeError', handleRouteChangeFinish);
    };
  }, [router.events]);

  return null;
};

export default RouteChangeIndicator;
