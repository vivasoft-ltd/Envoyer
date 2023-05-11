import React, {useEffect, useState} from "react";
import {Button, List, Modal} from "antd";
import ReactDragListView from "react-drag-listview";
import {CheckCircleFilled, CloseCircleFilled} from "@ant-design/icons";
import ProviderItem from "@/components/provider/ProviderItem";
import {useMutation, useQueryClient} from "react-query";
import ProviderService from "@/services/providerService";
import {toast} from "react-toastify";

export default function ProviderListByType({
                                             allProvidersResp,
                                             appId,
                                             notificationType
                                           }: { allProvidersResp: ProviderResponse[] | undefined, appId: number, notificationType: string }) {
  const [listData, setListData] = useState(allProvidersResp || []);
  const [showPriorityModal, setShowPriorityModal] = useState(false)
  const providerService = new ProviderService();
  const queryClient = useQueryClient();
  const hideModal = () => {
    setShowPriorityModal(false);
  };

  useEffect(() => {
    setListData(allProvidersResp || [])
  }, [allProvidersResp])

  const clearForm = () => {
    setListData(allProvidersResp || []);
  };

  function handleCreateProvider(e: React.MouseEvent<HTMLAnchorElement, MouseEvent> | React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    e.preventDefault();
    let ediData: Priority[] = []
    for (let i = 0; i < listData.length; i++) {
      let a: Priority = {
        priority: i + 1,
        id: listData[i].ID,
      }
      ediData.push(a)
    }

    editPriorityProvider(
      ediData
      , {
        onError: (err: any) => {
          toast.error(err?.message ? err?.message : 'Error setting priority to the provider');
        },
        onSuccess: async () => {
          hideModal();
          await queryClient.invalidateQueries(['getAllSmsProvidersOfApp']);
          await queryClient.invalidateQueries(['getAllEmailProvidersOfApp']);
          await queryClient.invalidateQueries(['getAllPushProvidersOfApp']);
          await queryClient.invalidateQueries(['getAllWebhookProvidersOfApp']);
          toast.success('Provider priority updated successfully');
          clearForm();
        },
      });
  }

  const {mutate: editPriorityProvider} = useMutation(
    (data: Priority[]) =>
      providerService.editProviderPriority(data, appId, notificationType),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  return (
    <div>
      <Modal
        centered
        open={showPriorityModal}
        title="Edit Provider Priority"
        onCancel={hideModal}
        footer={null}
        maskClosable={false}
      >
        <div className="pt-4">
          <div className="pb-4 text-sm">
            Drag and drop to update the priority of the providers and click save.
          </div>
          <ReactDragListView nodeSelector=".drag-handle"
                             onDragEnd={(fromIndex, toIndex) => {
                               const items = Array.from(listData);
                               const [reorderedItem] = items.splice(fromIndex, 1);
                               items.splice(toIndex, 0, reorderedItem);
                               setListData(items);
                             }}>
            <List
              renderItem={(item) => (
                <div className='drag-handle cursor-move' style={{border: "1px solid rgba(5, 5, 5, 0.06)"}}>
                  <List.Item>
                    <List.Item.Meta
                      avatar={item.active ?
                        <CheckCircleFilled style={{color: '#609966'}}/> :
                        <CloseCircleFilled style={{color: '#FF4d4f'}}/>
                      }
                      title={item.name}
                      description={item.description}
                    >
                    </List.Item.Meta>
                    <div className="space-x-2">
                      Priority = {item.priority}
                      <br/>
                      Provider = {item.provider_type}
                    </div>
                  </List.Item>
                </div>
              )}
              itemLayout="horizontal"
              dataSource={listData}
            />
          </ReactDragListView>
          <div className="lg:text-right">
            <Button type='primary' className='bg-blue-700 mt-5' onClick={handleCreateProvider}>Save</Button>
          </div>
        </div>
      </Modal>
      <div className="lg:text-right">
        {listData.length === 0 ||
            <Button type='primary' className='bg-blue-700 mb-5' onClick={() => setShowPriorityModal(true)}>Edit
                Priority</Button>}
      </div>
      <div style={{borderRadius: "8px"}}>
        <List
          renderItem={(item) => (
            <ProviderItem providerDetails={item}/>
          )}
          itemLayout="horizontal"
          dataSource={allProvidersResp}
        />
      </div>
    </div>
  );
}


