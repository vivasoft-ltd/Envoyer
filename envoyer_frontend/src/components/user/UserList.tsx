import {List} from 'antd';
import {useQuery} from 'react-query';
import UserService from "@/services/userService";
import UserItem from "@/components/user/UserItem";

export default function UserList({appId, canEdit = true}: { appId: number, canEdit?: boolean }) {
  const userService = new UserService();

  const {data: allUsersResp} = useQuery(['getAllUsers'], () =>
    userService.getAllUsers(appId)
  );

  return (
    <div>
      <List
        renderItem={(item) => (
          <UserItem userDetails={item} canEdit={canEdit}/>
        )}
        itemLayout="horizontal"
        dataSource={allUsersResp}
      />
    </div>
  );
}
