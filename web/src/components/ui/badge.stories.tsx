import type { Meta, StoryObj } from '@storybook/nextjs';
import { Home } from "lucide-react";

import { Badge } from './badge';

const meta = {
  component: Badge,
} satisfies Meta<typeof Badge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    variant: 'default',
  },
};

export const WithIcon: Story  = {
    args: {
        variant: "default",
    },
    render: (args) => (
        <Badge {...args}>
            <Home/>
            Home
        </Badge>
    )
}

export const WithIconOutline: Story = {
    args: {
        variant: "outline"
    },

    render: args => (
        <Badge {...args}>
            <Home />Home
        </Badge>
    )
};