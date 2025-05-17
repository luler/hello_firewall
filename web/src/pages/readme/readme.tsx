import {PageContainer} from '@ant-design/pro-components';
import React, {useEffect, useState} from 'react';
import MdEditor from "for-editor";

const Index: React.FC = () => {
  const [markdownContent, setMarkdownContent] = useState('');
  useEffect(() => {
    async function fetchMarkdownFile() {
      try {
        const response = await fetch('/api/README.md');
        const markdownText = await response.text();
        setMarkdownContent(markdownText);
      } catch (error) {
        console.error(error);
      }
    }

    fetchMarkdownFile();
  }, []);
  return (
    <PageContainer style={{minHeight: window.innerHeight - 150}}>

      <MdEditor
        value={markdownContent}
        preview={true}
        toolbar={{}}
        height="100%"
      />
    </PageContainer>
  );
};

export default Index;
