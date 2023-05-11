import {List} from 'antd';
import {useQuery} from 'react-query';
import AppServices from '@/services/appServices';
import AppItem from "@/components/app/AppItem";

export default function AppList() {
  const appServices = new AppServices();

  const {data: allAppsResp} = useQuery(['getAllApps'], () =>
    appServices.getAllApps()
  );

  return (
    <div>
      <List
        renderItem={(item) => (
          <AppItem appDetails={item}/>
        )}
        itemLayout="horizontal"
        dataSource={allAppsResp}
      />
    </div>
  );
}
