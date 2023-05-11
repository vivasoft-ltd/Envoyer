type PublishInput = {
  app_key: string;
  client_key: string;
  [key: string]: any;
  // event_name: string;
  // delivery_time: any;
  // receivers_with_variables: ReceiversWithVariable[];
};

type ReceiversWithVariable = {
  receiver: string;
  variables: ReceiverVariable[];
};

type ReceiverVariable = {
  name: string;
  value: string;
};
