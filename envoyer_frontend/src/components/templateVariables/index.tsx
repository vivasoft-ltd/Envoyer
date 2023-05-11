import React, { useEffect, useRef, useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import type { InputRef } from 'antd';
import { Input, Tag, theme } from 'antd';
import { transformToVariable } from '@/utils/stringFormat';

interface Props {
  variables: string[];
  setVariables: (variables: string[]) => void;
  creatable?: boolean;
}

const TemplateVariables = ({ variables, setVariables, creatable }: Props) => {
  const { token } = theme.useToken();
  const [inputVisible, setInputVisible] = useState(false);
  const [inputValue, setInputValue] = useState('');
  const inputRef = useRef<InputRef>(null);

  useEffect(() => {
    if (inputVisible) {
      inputRef.current?.focus();
    }
  }, [inputVisible]);

  const handleClose = (removedTag: string) => {
    const newTags = variables.filter((tag) => tag !== removedTag);
    console.log(newTags);
    setVariables(newTags);
  };

  const showInput = () => setInputVisible(true);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  };

  const handleInputConfirm = () => {
    if (inputValue && variables.indexOf(inputValue) === -1) {
      setVariables([...variables, inputValue]);
    }
    setInputVisible(false);
    setInputValue('');
  };

  const forMap = (tag: string, index: number) => {
    if (tag === '') {
      return null;
    }

    return (
      <div key={tag + index} className='inline-block'>
        <Tag
          closable
          onClose={(e) => {
            e.preventDefault();
            handleClose(tag);
          }}
        >
          {transformToVariable(tag)}
        </Tag>
      </div>
    );
  };

  const tagChild = variables.map(forMap);

  const tagPlusStyle = {
    background: token.colorBgContainer,
    borderStyle: 'dashed',
  };

  if (creatable) {
    return (
      <>
        <div style={{ marginBottom: 16 }}>{tagChild}</div>
        {inputVisible ? (
          <Input
            ref={inputRef}
            type='text'
            size='small'
            style={{ width: 200 }}
            value={inputValue}
            onChange={handleInputChange}
            onBlur={handleInputConfirm}
            onPressEnter={handleInputConfirm}
          />
        ) : (
          <Tag onClick={showInput} style={tagPlusStyle}>
            <PlusOutlined /> Variables
          </Tag>
        )}
      </>
    );
  }

  return (
    <>
      <div style={{ marginBottom: 16 }}>{tagChild}</div>
    </>
  );
};

export default TemplateVariables;
