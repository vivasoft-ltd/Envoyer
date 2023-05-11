import {PropsWithChildren} from 'react';

const SiteLayout = ({children}: PropsWithChildren<{}>) => {
  return (
    <>
      <main>{children}</main>
    </>
  );
};

export default SiteLayout;
