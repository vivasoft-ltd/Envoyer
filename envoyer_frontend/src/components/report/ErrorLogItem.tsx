import {Button, List} from 'antd';
import { ExclamationCircleFilled} from '@ant-design/icons';
import {Modal} from 'antd';
import {useMutation, useQueryClient} from "react-query";
import {toast} from "react-toastify";
import React, {useState} from "react";
import Link from "next/link";
import ErrorLogInfo from "@/components/report/ErrorLogInfo";
import ErrorLogService from "@/services/errorlogService";
import { format } from 'date-fns'

export default function ErrorLogItem({errorDetails}: { errorDetails: ErrorLog }) {
  const {confirm} = Modal;
  const errorLogService = new ErrorLogService();
  const queryClient = useQueryClient();

  const [showDetails, setShowDetails] = useState(false);

  const hideDetailsModal = () => {
    setShowDetails(false);
  };

  const {mutate: deleteLog} = useMutation(
    (id: number) => errorLogService.deleteLog(id),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  function handleDeleteProvider(id: number) {
    deleteLog(id, {
      onError: (err: any) => {
        toast.error(err?.message ? err?.message : 'Error deleting log')
      },
      onSuccess: async () => {
        toast.success("Error log deleted successfully")
        await queryClient.invalidateQueries(['getAllErrors']);
      }
    })
  }

  const showDeleteConfirm = (id: number) => {
    confirm({
      title: 'Confirm delete error log',
      icon: <ExclamationCircleFilled/>,
      content: 'Are you sure you want to delete this log?',
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
          open={showDetails}
          title="Error Info"
          onCancel={hideDetailsModal}
          footer={null}
          width={"60%"}
        >
          <div className="pt-4">
            <ErrorLogInfo errorDetails={errorDetails}/>
          </div>
        </Modal>
      </div>

      <List.Item>

        <List.Item.Meta
          title={<Link href="#" onClick={() => setShowDetails(true)}>{errorDetails.message}</Link>}
          description={format(new Date(errorDetails.date!), 'yyyy-MM-dd hh:mm:ss a z')}
        >
        </List.Item.Meta>

        <div className="space-x-2">
          <Button type="primary" danger ghost onClick={() => showDeleteConfirm(errorDetails.ID)}>Delete</Button>
        </div>
      </List.Item>
    </div>
  );
}

