import {Button, List} from 'antd';
import {ExclamationCircleFilled} from '@ant-design/icons';
import {Modal} from 'antd';
import {useMutation, useQueryClient} from "react-query";
import {toast} from "react-toastify";
import React, {useState} from "react";
import EventCreate from "@/components/event/EventCreate";
import EventService from "@/services/eventService";
import ClientInfo from "@/components/client/ClientInfo";
import EventInfo from "@/components/event/EventInfo";
import Link from "next/link";

export default function EventItem({eventDetails}: { eventDetails: EventResponse }) {
  const {confirm} = Modal;
  const eventService = new EventService();
  const queryClient = useQueryClient();

  const [showEditEventModal, setShowEditEventModal] = useState(false);
  const [showDetails, setShowDetails] = useState(false);

  const hideDetailsModal = () => {
    setShowDetails(false);
  };

  const hideModal = () => {
    setShowEditEventModal(false);
  };

  const {mutate: deleteEvent} = useMutation(
    (id: number) => eventService.deleteEvent(id),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  function handleDeleteEvent(id: number) {
    deleteEvent(id, {
      onError: (err: any) => {
        toast.error(err?.message ? err?.message : 'Error deleting event')
      },
      onSuccess: async () => {
        toast.success("Event deleted successfully")
        await queryClient.invalidateQueries(['getAllEvents']);
      }
    })
  }

  const showDeleteConfirm = (id: number) => {
    confirm({
      title: 'Confirm delete event',
      icon: <ExclamationCircleFilled/>,
      content: 'Are you sure you want to delete this event?',
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      closable: true,
      onOk() {
        handleDeleteEvent(id)
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
          open={showEditEventModal}
          title="Edit Event"
          onCancel={hideModal}
          footer={null}
          maskClosable={false}
        >
          <div className="pt-4">
            <EventCreate IsEdit={true} EditEventDetails={eventDetails} hideModal={hideModal}/>
          </div>
        </Modal>
        <Modal
          centered
          open={showDetails}
          title="Event Info"
          onCancel={hideDetailsModal}
          footer={null}
        >
          <div className="pt-4">
            <EventInfo eventDetails={eventDetails}/>
          </div>
        </Modal>
      </div>

      <List.Item>

        <List.Item.Meta
          title={<Link href="#" onClick={() => setShowDetails(true)}>{eventDetails.name}</Link>}
          description={eventDetails.description}
        >
          {eventDetails?.variables?.join(',')}
        </List.Item.Meta>
        <div className="space-x-2">
          <Button type="primary" ghost onClick={() => setShowEditEventModal(true)}>Edit</Button>
          <Button type="primary" danger ghost onClick={() => showDeleteConfirm(eventDetails.ID)}>Delete</Button>
        </div>
      </List.Item>
    </div>
  );
}

