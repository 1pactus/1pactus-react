import { 
  Card, 
  CardHeader, 
  CardTitle, 
  CardDescription, 
  CardContent, 
  CardFooter,
  CardAction 
} from './card'
import { Button } from './button'

const meta = {
  title: 'UI/Card',
  component: Card,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
}

export default meta

// 基础卡片
export const Default = {
  render: () => (
    <Card className="w-96">
      <CardHeader>
        <CardTitle>卡片标题</CardTitle>
        <CardDescription>这是一个卡片的描述信息</CardDescription>
      </CardHeader>
      <CardContent>
        <p>这里是卡片的主要内容。可以放置任何你想要的内容。</p>
      </CardContent>
      <CardFooter>
        <Button>操作按钮</Button>
      </CardFooter>
    </Card>
  ),
}

// 带操作按钮的卡片
export const WithAction = {
  render: () => (
    <Card className="w-96">
      <CardHeader>
        <CardTitle>通知设置</CardTitle>
        <CardDescription>管理您的通知偏好</CardDescription>
        <CardAction>
          <Button variant="outline" size="sm">
            编辑
          </Button>
        </CardAction>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span>邮件通知</span>
            <span className="text-green-600">已开启</span>
          </div>
          <div className="flex items-center justify-between">
            <span>短信通知</span>
            <span className="text-gray-500">已关闭</span>
          </div>
        </div>
      </CardContent>
    </Card>
  ),
}

// 用户资料卡片
export const UserProfile = {
  render: () => (
    <Card className="w-96">
      <CardHeader>
        <div className="flex items-center space-x-4">
          <div className="w-12 h-12 bg-gradient-to-r from-blue-400 to-purple-600 rounded-full flex items-center justify-center text-white font-bold">
            张
          </div>
          <div>
            <CardTitle>张三</CardTitle>
            <CardDescription>高级开发工程师</CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <div className="flex justify-between">
            <span className="text-sm text-gray-500">邮箱</span>
            <span className="text-sm">zhangsan@example.com</span>
          </div>
          <div className="flex justify-between">
            <span className="text-sm text-gray-500">部门</span>
            <span className="text-sm">技术部</span>
          </div>
          <div className="flex justify-between">
            <span className="text-sm text-gray-500">入职时间</span>
            <span className="text-sm">2023年1月</span>
          </div>
        </div>
      </CardContent>
      <CardFooter className="gap-2">
        <Button variant="outline" className="flex-1">查看详情</Button>
        <Button className="flex-1">发送消息</Button>
      </CardFooter>
    </Card>
  ),
}

// 统计数据卡片
export const Statistics = {
  render: () => (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 w-full max-w-4xl">
      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-medium">总销售额</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">¥45,231</div>
          <p className="text-xs text-green-600">+20.1% 比上月</p>
        </CardContent>
      </Card>
      
      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-medium">用户数量</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">+2,350</div>
          <p className="text-xs text-green-600">+180.1% 比上月</p>
        </CardContent>
      </Card>
      
      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-medium">活跃用户</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">+12,234</div>
          <p className="text-xs text-green-600">+19% 比上月</p>
        </CardContent>
      </Card>
    </div>
  ),
}

// 产品卡片
export const Product = {
  render: () => (
    <Card className="w-80">
      <CardContent>
        <div className="aspect-square bg-gray-100 rounded-lg mb-4"></div>
        <CardTitle className="mb-2">无线蓝牙耳机</CardTitle>
        <CardDescription>高品质音乐体验，降噪技术，长续航</CardDescription>
        <div className="mt-4 flex items-center justify-between">
          <span className="text-2xl font-bold text-red-600">¥299</span>
          <span className="text-sm text-gray-500 line-through">¥399</span>
        </div>
      </CardContent>
      <CardFooter>
        <Button className="w-full">立即购买</Button>
      </CardFooter>
    </Card>
  ),
}

// 最小化卡片
export const Minimal = {
  render: () => (
    <Card className="w-64">
      <CardContent>
        <p>这是一个最简单的卡片，只包含内容。</p>
      </CardContent>
    </Card>
  ),
}
