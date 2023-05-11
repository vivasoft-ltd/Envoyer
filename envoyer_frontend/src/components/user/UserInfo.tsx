import UserService from "@/services/userService";
import {useQuery} from "react-query";
import {Descriptions} from "antd";
import React from "react";

export default function UserInfo({userId}: { userId: number }) {
  const userService = new UserService();

  const {data: usersResp} = useQuery(['getUserInfo'], () =>
    userService.getUser(userId)
  );

  return (
    <Descriptions title="User Information" layout="vertical" bordered className="pt-5">
      <Descriptions.Item label="Username">{usersResp?.user_name}</Descriptions.Item>
      <Descriptions.Item label="User Id">{usersResp?.ID}</Descriptions.Item>
      <Descriptions.Item label="Role">{usersResp?.role}</Descriptions.Item>
    </Descriptions>
  );
}
