import {Button, List} from 'antd';
import Link from "next/link";
import {CheckCircleFilled, CheckCircleTwoTone, CloseCircleFilled, CloseCircleTwoTone} from '@ant-design/icons';
import {ExclamationCircleFilled} from '@ant-design/icons';
import {Modal} from 'antd';
import {useMutation, useQueryClient} from "react-query";
import AppServices from "@/services/appServices";
import {toast} from "react-toastify";
import React, {useState} from "react";
import AppCreate from "@/components/app/AppCreate";

export default function AppItem({appDetails}: { appDetails: ListAppResponse }) {
  const {confirm} = Modal;
  const appServices = new AppServices();
  const queryClient = useQueryClient();

  const [showCreateAppModal, setShowCreateAppModal] = useState(false);

  const hideModal = () => {
    setShowCreateAppModal(false);
  };

  const {mutate: deleteApp} = useMutation(
    (id: number) => appServices.deleteApp(id),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  function handleDeleteApp(id: number) {
    deleteApp(id, {
      onError: (err: any) => {
        toast.error(err?.message ? err?.message : 'Error deleting app')
      },
      onSuccess: async () => {
        toast.success("App deleted successfully")
        await queryClient.invalidateQueries(['getAllApps']);
      }
    })
  }

  const showDeleteConfirm = (id: number) => {
    confirm({
      title: 'Confirm delete app',
      icon: <ExclamationCircleFilled/>,
      content: 'Are you sure you want to delete this app?',
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      closable: true,
      onOk() {
        handleDeleteApp(id)
      },
      onCancel() {
      },
    });
  };

  return (
    <div style={{border: "1px solid rgba(5, 5, 5, 0.06)"}}>
      <div>
        {showCreateAppModal && (
          /*<GeneralModal title="Edit App" hideModal={hideModal} id="edit-app">
            <AppCreate IsEdit={true} EditAppDetails={appDetails} hideModal={hideModal}/>
          </GeneralModal>*/
          <Modal
            centered
            open={showCreateAppModal}
            title="Edit App"
            onCancel={hideModal}
            footer={null}
            maskClosable={false}
          >
            <div className="pt-4">
              <AppCreate IsEdit={true} EditAppDetails={appDetails} hideModal={hideModal}/>
            </div>
          </Modal>
        )
        }
      </div>

      <List.Item>
        <List.Item.Meta
          avatar={appDetails.active ?
            <CheckCircleFilled style={{color: '#609966'}}/> :
            <CloseCircleFilled style={{color: '#FF4d4f'}}/>
          }
          title={<Link href={"/dashboard/" + appDetails.ID}>{appDetails.name}</Link>}
          description={appDetails.description}
        />
        <div className="space-x-2">
          <Button type="primary" ghost onClick={() => setShowCreateAppModal(true)}>Edit</Button>
          <Button type="primary" danger ghost onClick={() => showDeleteConfirm(appDetails.ID)}>Delete</Button>
        </div>
      </List.Item>
    </div>
  );
}

