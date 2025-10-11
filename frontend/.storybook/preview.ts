import type { Preview } from '@storybook/nextjs';
import React from 'react';
import '../src/app/globals.css';

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
    nextjs: {
      appDirectory: true,
      navigation: {
        pathname: '/dashboard/home',
        query: {},
      },
    },
  },
  decorators: [
    (Story) => React.createElement('div', { style: { minHeight: '100vh' } }, React.createElement(Story)),
  ],
};

export default preview;
