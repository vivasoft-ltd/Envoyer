import {List} from 'antd';
import {useQuery} from 'react-query';
import ErrorLogService from "@/services/errorlogService";
import ErrorLogItem from "@/components/report/ErrorLogItem";

export default function ErrorLogList({appId}: { appId: number }) {
  const errorLogService = new ErrorLogService();

  const {data: allErrorLogResp} = useQuery(['getAllErrors'], () =>
    errorLogService.getAllErrorLogs(appId)
  );

  return (
    <div>
      <List
        renderItem={(item) => (
          <ErrorLogItem errorDetails={item}/>
        )}
        itemLayout="horizontal"
        dataSource={allErrorLogResp}
      />
    </div>
  );
}
