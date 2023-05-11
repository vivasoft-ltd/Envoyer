import React, { forwardRef, useEffect, useRef, useState } from 'react';
import EmailEditor from 'react-email-editor';

interface Props {
  selectedTemplate?: string;
}

const EmailTemplate = forwardRef((props: Props, ref: any) => {
  const { selectedTemplate } = props;

  const onLoad = () => {
    selectedTemplate &&
      ref?.current?.editor?.loadDesign(JSON.parse(selectedTemplate));
  };

  const onReady = () => {};

  return (
    <div>
      <EmailEditor
        minHeight={600}
        style={{ minWidth: '100%' }}
        ref={ref}
        onLoad={onLoad}
        onReady={onReady}
        projectId={138679}
        // options={{ innerWidth: 600 }}
        // appearance=''
        // options={{ blocks: [headerBlock, footerBlock] }}
      />
    </div>
  );
});

EmailTemplate.displayName = 'EmailTemplate';

export default EmailTemplate;
