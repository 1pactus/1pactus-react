import { Button } from './button'

// Meta 配置：定义组件的基本信息
const meta = {
  title: 'UI/Button', // 在 Storybook 中的组织路径
  component: Button,
  parameters: {
    layout: 'centered', // 在画布中居中显示
  },
  tags: ['autodocs'], // 自动生成文档
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'destructive', 'outline', 'secondary', 'ghost', 'link'],
      description: '按钮的视觉变体'
    },
    size: {
      control: 'select',
      options: ['default', 'sm', 'lg', 'icon'],
      description: '按钮的尺寸'
    },
    asChild: {
      control: 'boolean',
      description: '是否作为子元素渲染'
    },
    disabled: {
      control: 'boolean',
      description: '是否禁用按钮'
    }
  },
}

export default meta

// 基础故事：展示默认状态
export const Default = {
  args: {
    children: '按钮',
  },
}

// 不同变体的故事
export const Destructive = {
  args: {
    variant: 'destructive',
    children: '删除',
  },
}

export const Outline = {
  args: {
    variant: 'outline',
    children: '轮廓按钮',
  },
}

export const Secondary = {
  args: {
    variant: 'secondary',
    children: '次要按钮',
  },
}

export const Ghost = {
  args: {
    variant: 'ghost',
    children: '幽灵按钮',
  },
}

export const Link = {
  args: {
    variant: 'link',
    children: '链接按钮',
  },
}

// 不同尺寸的故事
export const Small = {
  args: {
    size: 'sm',
    children: '小按钮',
  },
}

export const Large = {
  args: {
    size: 'lg',
    children: '大按钮',
  },
}

export const Icon = {
  args: {
    size: 'icon',
    children: '🔥',
  },
}

// 状态演示
export const Disabled = {
  args: {
    disabled: true,
    children: '禁用按钮',
  },
}

// 带图标的按钮（使用 Lucide React 图标）
export const WithIcon = {
  render: (args: Record<string, unknown>) => (
    <Button {...args}>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      >
        <path d="M5 12h14" />
        <path d="m12 5 7 7-7 7" />
      </svg>
      下一步
    </Button>
  ),
}

// 所有变体的组合展示
export const AllVariants = {
  render: () => (
    <div className="flex flex-wrap gap-4">
      <Button variant="default">默认</Button>
      <Button variant="destructive">销毁</Button>
      <Button variant="outline">轮廓</Button>
      <Button variant="secondary">次要</Button>
      <Button variant="ghost">幽灵</Button>
      <Button variant="link">链接</Button>
    </div>
  ),
}

// 所有尺寸的组合展示
export const AllSizes = {
  render: () => (
    <div className="flex items-center gap-4">
      <Button size="sm">小</Button>
      <Button size="default">默认</Button>
      <Button size="lg">大</Button>
      <Button size="icon">🎯</Button>
    </div>
  ),
}
