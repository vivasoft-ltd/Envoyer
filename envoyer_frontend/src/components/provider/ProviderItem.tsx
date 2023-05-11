import {Button, List} from 'antd';
import {
  CheckCircleFilled,
  CheckCircleTwoTone,
  CloseCircleFilled,
  CloseCircleTwoTone,
  ExclamationCircleFilled
} from '@ant-design/icons';
import {Modal} from 'antd';
import {useMutation, useQueryClient} from "react-query";
import {toast} from "react-toastify";
import React, {useState} from "react";
import ProviderCreate from "@/components/provider/ProviderCreate";
import ProviderService from "@/services/providerService";
import Link from "next/link";
import ProviderInfo from "@/components/provider/ProviderInfo";

export default function ProviderItem({providerDetails}: { providerDetails: ProviderResponse }) {
  const {confirm} = Modal;
  const providerService = new ProviderService();
  const queryClient = useQueryClient();

  const [showEditProviderModal, setShowEditProviderModal] = useState(false);
  const [showDetails, setShowDetails] = useState(false);

  const hideDetailsModal = () => {
    setShowDetails(false);
  };

  const hideModal = () => {
    setShowEditProviderModal(false);
  };

  const {mutate: deleteProvider} = useMutation(
    (id: number) => providerService.deleteProvider(id),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  function handleDeleteProvider(id: number) {
    deleteProvider(id, {
      onError: (err: any) => {
        toast.error(err?.message ? err?.message : 'Error deleting provider')
      },
      onSuccess: async () => {
        toast.success("Provider deleted successfully")
        await queryClient.invalidateQueries(['getAllSmsProvidersOfApp']);
        await queryClient.invalidateQueries(['getAllEmailProvidersOfApp']);
      }
    })
  }

  const showDeleteConfirm = (id: number) => {
    confirm({
      title: 'Confirm delete provider',
      icon: <ExclamationCircleFilled/>,
      content: 'Are you sure you want to delete this provider?',
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      closable: true,
      onOk() {
        handleDeleteProvider(id)
      },
      onCancel() {
      },
    });
  };

  return (
    <div style={{border: "1px solid rgba(5, 5, 5, 0.06)"}}>
      <div>
        <Modal
          centered
          open={showEditProviderModal}
          title="Edit Provider"
          onCancel={hideModal}
          footer={null}
          maskClosable={false}
        >
          <div className="pt-4">
            <ProviderCreate IsEdit={true} EditProviderDetails={providerDetails} hideModal={hideModal}/>
          </div>
        </Modal>
        <Modal
          centered
          open={showDetails}
          title="Provider Info"
          onCancel={hideDetailsModal}
          footer={null}
          width={"60%"}
        >
          <div className="pt-4">
            <ProviderInfo providerDetails={providerDetails}/>
          </div>
        </Modal>
      </div>

      <List.Item>

        <List.Item.Meta
          avatar={providerDetails.active ?
            <CheckCircleFilled style={{color: '#609966'}}/> :
            <CloseCircleFilled style={{color: '#FF4d4f'}}/>
          }
          title={<Link href="#" onClick={() => setShowDetails(true)}>{providerDetails.name}</Link>}
          description={providerDetails.description}
        >
          Type = {providerDetails.type}
          Provider Type = {providerDetails.provider_type}
        </List.Item.Meta>

        <div className="space-x-2">
          <Button type="primary" ghost onClick={() => setShowEditProviderModal(true)}>Edit</Button>
          <Button type="primary" danger ghost onClick={() => showDeleteConfirm(providerDetails.ID)}>Delete</Button>
        </div>
      </List.Item>
    </div>
  );
}

