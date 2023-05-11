import {List} from 'antd';
import {useQuery} from 'react-query';
import ClientService from "@/services/clientService";
import ClientItem from "@/components/client/ClientItem";

export default function ClientList({appId}: { appId: number }) {
  const clientService = new ClientService();

  const {data: allClientsResp} = useQuery(['getAllClients'], () =>
    clientService.getAllClients(appId)
  );

  return (
    <div>
      <List
        renderItem={(item) => (
          <ClientItem clientDetails={item}/>
        )}
        itemLayout="horizontal"
        dataSource={allClientsResp}
      />
    </div>
  );
}
